import { initializeApp } from 'firebase/app'
import {
    GoogleAuthProvider,
    getAuth,
    onAuthStateChanged,
    signInWithEmailAndPassword,
    signInWithPopup,
    signOut as firebaseSignOut,
} from 'firebase/auth'
import type { User } from 'firebase/auth'

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID,
}

const app = initializeApp(firebaseConfig)
export const auth = getAuth(app)

const googleProvider = new GoogleAuthProvider()

export async function signInWithEmail(email: string, password: string) {
    return signInWithEmailAndPassword(auth, email, password)
}

export async function signInWithGoogle() {
    return signInWithPopup(auth, googleProvider)
}

export async function signOut() {
    await firebaseSignOut(auth)
}

export function onAuthReady(callback: (user: User | null) => void) {
    return onAuthStateChanged(auth, callback)
}

export async function getValidToken(): Promise<string | null> {
    const user = auth.currentUser
    if (!user) {
        return null
    }

    const tokenResult = await user.getIdTokenResult()
    const expirationTime = tokenResult.expirationTime
    const expiresAt = new Date(expirationTime).getTime()
    const now = Date.now()
    const fiveMinutes = 5 * 60 * 1000
    const forceRefresh = expiresAt - now < fiveMinutes

    return user.getIdToken(forceRefresh)
}
