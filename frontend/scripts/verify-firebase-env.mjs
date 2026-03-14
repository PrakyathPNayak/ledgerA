import fs from 'node:fs'
import path from 'node:path'

const envPath = path.resolve('.env')

if (!fs.existsSync(envPath)) {
    console.error('FAIL: frontend/.env is missing')
    process.exit(1)
}

const content = fs.readFileSync(envPath, 'utf8')
const required = [
    'VITE_FIREBASE_API_KEY',
    'VITE_FIREBASE_AUTH_DOMAIN',
    'VITE_FIREBASE_PROJECT_ID',
    'VITE_FIREBASE_APP_ID',
]

const envMap = new Map()
for (const rawLine of content.split(/\r?\n/)) {
    const line = rawLine.trim()
    if (!line || line.startsWith('#')) continue
    const idx = line.indexOf('=')
    if (idx <= 0) continue
    const key = line.slice(0, idx).trim()
    const value = line.slice(idx + 1).trim()
    envMap.set(key, value)
}

const missing = required.filter((key) => !envMap.get(key))
if (missing.length > 0) {
    console.error(`FAIL: Missing Firebase env vars: ${missing.join(', ')}`)
    process.exit(1)
}

const apiKey = envMap.get('VITE_FIREBASE_API_KEY')
if (!apiKey.startsWith('AIza')) {
    console.error('FAIL: VITE_FIREBASE_API_KEY format looks invalid (expected prefix AIza)')
    process.exit(1)
}

console.log('PASS: Firebase frontend env vars are present and non-empty')
