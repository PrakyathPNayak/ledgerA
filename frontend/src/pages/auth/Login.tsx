import { useState } from 'react'
import type { FormEvent } from 'react'
import { Navigate, useNavigate } from 'react-router-dom'

import { Alert } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useSyncUser } from '@/hooks/useAuth'
import { signInWithEmail, signInWithGoogle } from '@/lib/firebase'
import { useAuthStore } from '@/store/authStore'

export function LoginPage() {
    const navigate = useNavigate()
    const syncUser = useSyncUser()
    const { user, isLoading } = useAuthStore()

    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [errorMessage, setErrorMessage] = useState<string | null>(null)
    const [isEmailLoading, setIsEmailLoading] = useState(false)
    const [isGoogleLoading, setIsGoogleLoading] = useState(false)

    if (!isLoading && user) {
        return <Navigate to="/" replace />
    }

    async function syncAuthenticatedUser(token: string, displayName: string, userEmail: string) {
        await syncUser.mutateAsync({
            firebase_token: token,
            display_name: displayName,
            email: userEmail,
        })
        navigate('/', { replace: true })
    }

    async function handleEmailSignIn(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        setErrorMessage(null)
        setIsEmailLoading(true)

        try {
            const credential = await signInWithEmail(email.trim(), password)
            const token = await credential.user.getIdToken()
            await syncAuthenticatedUser(
                token,
                credential.user.displayName ?? credential.user.email ?? 'User',
                credential.user.email ?? email.trim(),
            )
        } catch (error) {
            const message = error instanceof Error ? error.message : 'Unable to sign in. Please try again.'
            setErrorMessage(message)
        } finally {
            setIsEmailLoading(false)
        }
    }

    async function handleGoogleSignIn() {
        setErrorMessage(null)
        setIsGoogleLoading(true)

        try {
            const credential = await signInWithGoogle()
            const token = await credential.user.getIdToken()
            await syncAuthenticatedUser(
                token,
                credential.user.displayName ?? credential.user.email ?? 'User',
                credential.user.email ?? 'unknown@example.com',
            )
        } catch (error) {
            const message = error instanceof Error ? error.message : 'Unable to continue with Google.'
            setErrorMessage(message)
        } finally {
            setIsGoogleLoading(false)
        }
    }

    return (
        <div className="min-h-screen bg-app px-4">
            <div className="mx-auto flex min-h-screen w-full max-w-md items-center justify-center">
                <div className="w-full rounded-2xl border border-border bg-surface p-7 shadow-lg">
                    <h1 className="text-2xl font-bold tracking-tight text-foreground">Sign In</h1>
                    <p className="mt-1 text-sm text-muted">Continue to Expenditure Tracker</p>

                    <form className="mt-6 space-y-4" onSubmit={handleEmailSignIn}>
                        <div className="space-y-1">
                            <label htmlFor="email" className="text-sm font-medium text-secondary">
                                Email
                            </label>
                            <Input
                                id="email"
                                type="email"
                                autoComplete="email"
                                value={email}
                                onChange={(event) => setEmail(event.target.value)}
                                placeholder="you@example.com"
                                required
                            />
                        </div>

                        <div className="space-y-1">
                            <label htmlFor="password" className="text-sm font-medium text-secondary">
                                Password
                            </label>
                            <Input
                                id="password"
                                type="password"
                                autoComplete="current-password"
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                                placeholder="••••••••"
                                required
                            />
                        </div>

                        {errorMessage ? (
                            <Alert variant="destructive" title="Sign-in failed">
                                {errorMessage}
                            </Alert>
                        ) : null}

                        <Button type="submit" className="w-full" disabled={isEmailLoading || isGoogleLoading}>
                            {isEmailLoading ? 'Signing In...' : 'Sign In'}
                        </Button>
                    </form>

                    <div className="my-4 h-px bg-surface-hover" />

                    <Button
                        type="button"
                        variant="outline"
                        className="w-full"
                        onClick={handleGoogleSignIn}
                        disabled={isGoogleLoading || isEmailLoading}
                    >
                        {isGoogleLoading ? 'Connecting...' : 'Continue with Google'}
                    </Button>
                </div>
            </div>
        </div>
    )
}
