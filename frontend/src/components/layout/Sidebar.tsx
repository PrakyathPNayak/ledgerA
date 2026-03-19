import { Link, useLocation } from 'react-router-dom'
import {
  BarChart3,
  CreditCard,
  GitCompare,
  HelpCircle,
  LayoutDashboard,
  MessageCircle,
  Search,
  Zap,
} from 'lucide-react'

import { useAuthStore } from '@/store/authStore'
import { useUIStore } from '@/store/uiStore'

const navItems = [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/stats', label: 'Stats', icon: BarChart3 },
  { to: '/accounts', label: 'Accounts', icon: CreditCard },
  { to: '/search', label: 'Search', icon: Search },
  { to: '/quick-transactions', label: 'Quick Actions', icon: Zap },
  { to: '/compare', label: 'Compare', icon: GitCompare },
  { to: '/chat', label: 'Chat', icon: MessageCircle },
  { to: '/help', label: 'Help', icon: HelpCircle },
]

export function Sidebar() {
  const location = useLocation()
  const { sidebarCollapsed, toggleSidebar } = useUIStore()
  const { user } = useAuthStore()

  return (
    <aside
      className={`border-r border-slate-200 bg-white transition-all dark:border-slate-800 dark:bg-slate-900 ${sidebarCollapsed ? 'w-20' : 'w-64'
        }`}
    >
      <div className="flex items-center justify-between border-b border-slate-200 px-4 py-3 dark:border-slate-800">
        <span className="text-sm font-semibold text-slate-700 dark:text-slate-100">ledgerA</span>
        <button className="text-xs text-slate-500 dark:text-slate-400" onClick={toggleSidebar}>
          {sidebarCollapsed ? '>>' : '<<'}
        </button>
      </div>
      <nav className="flex flex-col gap-1 p-2">
        {navItems.map((item) => {
          const Icon = item.icon
          const active = location.pathname === item.to
          return (
            <Link
              key={item.to}
              to={item.to}
              className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm ${active
                ? 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900'
                : 'text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-300 dark:hover:bg-slate-800 dark:hover:text-slate-100'
                }`}
            >
              <Icon size={16} />
              {!sidebarCollapsed && (
                <>
                  <span>{item.label}</span>
                  {item.soon && (
                    <span className="ml-auto rounded-full bg-amber-100 px-2 py-0.5 text-xs text-amber-700">
                      Soon
                    </span>
                  )}
                </>
              )}
            </Link>
          )
        })}
      </nav>
      <div className="mt-auto border-t border-slate-200 p-3 text-xs text-slate-500 dark:border-slate-800 dark:text-slate-400">
        {user ? `${user.displayName ?? user.email}` : 'Guest'}
      </div>
    </aside>
  )
}
