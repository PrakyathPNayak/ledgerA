import os
import sys

import requests
from google.auth.transport.requests import Request
from google.oauth2 import service_account

cred_path = os.environ.get("FIREBASE_CREDENTIALS_PATH")
project_id = os.environ.get("FIREBASE_PROJECT_ID")
if not cred_path or not project_id:
    print("MISSING_ENV")
    sys.exit(1)

scopes = ["https://www.googleapis.com/auth/cloud-platform"]
creds = service_account.Credentials.from_service_account_file(cred_path, scopes=scopes)
creds.refresh(Request())
headers = {"Authorization": f"Bearer {creds.token}"}

list_url = f"https://firebase.googleapis.com/v1beta1/projects/{project_id}/webApps"
resp = requests.get(list_url, headers=headers, timeout=30)
print("LIST_STATUS", resp.status_code)
if resp.status_code != 200:
    print(resp.text[:500])
    sys.exit(1)

apps = resp.json().get("apps", [])
if not apps:
    print("NO_WEB_APPS")
    sys.exit(2)

app_name = apps[0]["name"]
config_url = f"https://firebase.googleapis.com/v1beta1/{app_name}/config"
cfg_resp = requests.get(config_url, headers=headers, timeout=30)
print("CONFIG_STATUS", cfg_resp.status_code)
if cfg_resp.status_code != 200:
    print(cfg_resp.text[:500])
    sys.exit(1)

cfg = cfg_resp.json()
api_key = cfg.get("apiKey", "")
auth_domain = cfg.get("authDomain", "")
project = cfg.get("projectId", "")
app_id = cfg.get("appId", "")

print("API_KEY_PREFIX", api_key[:4])
print("AUTH_DOMAIN", auth_domain)
print("PROJECT_ID", project)
print("APP_ID_PREFIX", app_id[:6])

frontend_env = os.path.join(os.getcwd(), "frontend", ".env")
with open(frontend_env, "w", encoding="utf-8") as f:
    f.write(f"VITE_FIREBASE_API_KEY={api_key}\n")
    f.write(f"VITE_FIREBASE_AUTH_DOMAIN={auth_domain}\n")
    f.write(f"VITE_FIREBASE_PROJECT_ID={project}\n")
    f.write(f"VITE_FIREBASE_APP_ID={app_id}\n")
print("WROTE_FRONTEND_ENV", frontend_env)
