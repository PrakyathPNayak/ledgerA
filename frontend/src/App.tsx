import { Navigate, Outlet, createBrowserRouter } from 'react-router-dom'

import { AccountsPage } from '@/pages/AccountsPage'
import { ChatPage } from '@/pages/ChatPage'
import { ComparePage } from '@/pages/ComparePage'
import { DashboardPage } from '@/pages/DashboardPage'
import { HelpPage } from '@/pages/HelpPage'
import { PassbookPage } from '@/pages/PassbookPage'
import { QuickTransactionsPage } from '@/pages/QuickTransactionsPage'
import { SearchPage } from '@/pages/SearchPage'
import { StatsPage } from '@/pages/Stats'
import { LoginPage } from '@/pages/auth/Login'
import { useAuthStore } from '@/store/authStore'

function RequireAuth() {
    const { user, isLoading } = useAuthStore()
    if (isLoading) {
        return <div className="grid min-h-screen place-items-center text-sm text-muted">Checking session...</div>
    }
    if (!user) {
        return <Navigate to="/login" replace />
    }
    return <Outlet />
}

export const appRouter = createBrowserRouter([
    { path: '/login', element: <LoginPage /> },
    {
        element: <RequireAuth />,
        children: [
            { path: '/', element: <DashboardPage /> },
            { path: '/stats', element: <StatsPage /> },
            { path: '/accounts', element: <AccountsPage /> },
            { path: '/accounts/:id', element: <PassbookPage /> },
            { path: '/search', element: <SearchPage /> },
            { path: '/quick-transactions', element: <QuickTransactionsPage /> },
            { path: '/compare', element: <ComparePage /> },
            { path: '/chat', element: <ChatPage /> },
            { path: '/help', element: <HelpPage /> },
        ],
    },
])
