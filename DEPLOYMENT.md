# Deploying ledgerA

This guide covers a production deployment on a single Ubuntu 22.04 LTS server using PostgreSQL 15, systemd for process management, and nginx as a reverse proxy.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Server Preparation](#server-preparation)
3. [PostgreSQL Setup](#postgresql-setup)
4. [Application User and Directory](#application-user-and-directory)
5. [Build the Binary and Frontend](#build-the-binary-and-frontend)
6. [Environment Configuration](#environment-configuration)
7. [Run Database Migrations](#run-database-migrations)
8. [systemd Service](#systemd-service)
9. [nginx Reverse Proxy](#nginx-reverse-proxy)
10. [TLS with Certbot](#tls-with-certbot)
11. [Verify](#verify)
12. [Updates and Rollback](#updates-and-rollback)

---

## Prerequisites

| Requirement | Notes |
|-------------|-------|
| Ubuntu 22.04 LTS | Any 64-bit Linux works with minor adjustments |
| Go 1.22+ | Install from [go.dev/dl](https://go.dev/dl/) |
| Node.js 20+ | Use [nvm](https://github.com/nvm-sh/nvm) or NodeSource PPA |
| PostgreSQL 15 | Via apt or Docker |
| nginx | `apt install nginx` |
| Certbot | `snap install --classic certbot` |
| Firebase service-account JSON | Downloaded from Firebase Console |

---

## Server Preparation

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y git build-essential nginx certbot python3-certbot-nginx
```

---

## PostgreSQL Setup

```bash
sudo apt install -y postgresql-15

# Switch to postgres user
sudo -u postgres psql <<'SQL'
CREATE USER exptracker WITH PASSWORD 'STRONG_PASSWORD_HERE';
CREATE DATABASE exptracker OWNER exptracker;
GRANT ALL PRIVILEGES ON DATABASE exptracker TO exptracker;
SQL
```

Note the connection string — you'll use it as `DATABASE_URL`:
```
postgres://exptracker:STRONG_PASSWORD_HERE@localhost:5432/exptracker?sslmode=disable
```

---

## Application User and Directory

```bash
sudo useradd --system --shell /usr/sbin/nologin --create-home ledgera

sudo mkdir -p /opt/ledgera/bin
sudo mkdir -p /opt/ledgera/frontend/dist
sudo mkdir -p /opt/ledgera/secrets
sudo chown -R ledgera:ledgera /opt/ledgera
```

---

## Build the Binary and Frontend

On your **build machine** (or the server itself):

```bash
# Clone the repository
git clone https://github.com/<your-org>/ledgerA.git
cd ledgerA

# Build Go binary (static, linux/amd64)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ledgera ./cmd/server

# Build frontend
cd frontend
npm ci
npm run build
cd ..
```

Copy artefacts to the server:

```bash
# Binary
scp ledgera user@your-server:/tmp/ledgera
ssh user@your-server "sudo mv /tmp/ledgera /opt/ledgera/bin/ledgera && sudo chown ledgera:ledgera /opt/ledgera/bin/ledgera && sudo chmod 755 /opt/ledgera/bin/ledgera"

# Frontend dist
rsync -avz frontend/dist/ user@your-server:/tmp/dist/
ssh user@your-server "sudo rsync -a /tmp/dist/ /opt/ledgera/frontend/dist/ && sudo chown -R ledgera:ledgera /opt/ledgera/frontend/dist"

# Migrations
rsync -avz migrations/ user@your-server:/tmp/migrations/
ssh user@your-server "sudo rsync -a /tmp/migrations/ /opt/ledgera/migrations/ && sudo chown -R ledgera:ledgera /opt/ledgera/migrations"
```

---

## Environment Configuration

```bash
sudo -u ledgera tee /opt/ledgera/.env <<'ENV'
DATABASE_URL=postgres://exptracker:STRONG_PASSWORD_HERE@localhost:5432/exptracker?sslmode=disable
FIREBASE_CREDENTIALS_PATH=/opt/ledgera/secrets/firebase-service-account.json
PORT=8080
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=https://your-domain.example.com
ENV

# Upload Firebase credentials
scp firebase-service-account.json user@your-server:/tmp/
ssh user@your-server "sudo mv /tmp/firebase-service-account.json /opt/ledgera/secrets/ && sudo chown ledgera:ledgera /opt/ledgera/secrets/firebase-service-account.json && sudo chmod 600 /opt/ledgera/secrets/firebase-service-account.json"
```

---

## Run Database Migrations

ledgerA uses raw SQL migration files. Apply them with `psql` before starting the service:

```bash
sudo -u ledgera bash -c '
  export DATABASE_URL="postgres://exptracker:STRONG_PASSWORD_HERE@localhost:5432/exptracker?sslmode=disable"
  for f in /opt/ledgera/migrations/*.up.sql; do
    echo "Applying $f"
    psql "$DATABASE_URL" -f "$f"
  done
'
```

Or use [golang-migrate](https://github.com/golang-migrate/migrate) CLI:

```bash
migrate -database "$DATABASE_URL" -path /opt/ledgera/migrations up
```

---

## systemd Service

The unit file is at [`deployments/exptracker.service`](deployments/exptracker.service).

Install it:

```bash
sudo cp deployments/exptracker.service /etc/systemd/system/exptracker.service
sudo systemctl daemon-reload
sudo systemctl enable exptracker
sudo systemctl start exptracker
sudo systemctl status exptracker
```

View logs:

```bash
sudo journalctl -u exptracker -f
```

---

## nginx Reverse Proxy

The config is at [`deployments/nginx.conf`](deployments/nginx.conf).

```bash
sudo cp deployments/nginx.conf /etc/nginx/sites-available/ledgera
sudo ln -s /etc/nginx/sites-available/ledgera /etc/nginx/sites-enabled/ledgera
sudo nginx -t
sudo systemctl reload nginx
```

---

## TLS with Certbot

```bash
sudo certbot --nginx -d your-domain.example.com
# Follow prompts — certbot will update nginx config and schedule auto-renewal
```

---

## Verify

```bash
# API health
curl https://your-domain.example.com/api/v1/health
# Expected: {"status":"ok"}

# Frontend (expect HTML)
curl -I https://your-domain.example.com/
```

---

## Updates and Rollback

### Update

```bash
# 1. Build new binary + frontend on build machine
# 2. Stop service
sudo systemctl stop exptracker

# 3. Backup current binary
sudo cp /opt/ledgera/bin/ledgera /opt/ledgera/bin/ledgera.bak

# 4. Deploy new artefacts (same rsync commands as above)

# 5. Apply any new migrations

# 6. Start service
sudo systemctl start exptracker
sudo journalctl -u exptracker -n 50 --no-pager
```

### Rollback

```bash
sudo systemctl stop exptracker
sudo cp /opt/ledgera/bin/ledgera.bak /opt/ledgera/bin/ledgera
# If migrations were applied, run the corresponding .down.sql files
sudo systemctl start exptracker
```

---

## Security Checklist

- [ ] `DATABASE_URL` password is strong and unique
- [ ] Firebase credentials file has mode `600` and is owned by `ledgera`
- [ ] `.env` file has mode `640` and is owned by `ledgera:ledgera`
- [ ] nginx is configured to forward `X-Forwarded-For` and `X-Real-IP`
- [ ] TLS certificate auto-renewal is verified: `sudo certbot renew --dry-run`
- [ ] Firewall allows only 80, 443 inbound; port 8080 is not publicly exposed
- [ ] `CORS_ALLOWED_ORIGINS` is set to the exact production origin — not `*`
