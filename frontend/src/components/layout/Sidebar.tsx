import { Link, useLocation } from 'react-router-dom'
import {
  BarChart3,
  CalendarClock,
  CreditCard,
  GitCompare,
  HelpCircle,
  LayoutDashboard,
  MessageCircle,
  PiggyBank,
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
  { to: '/recurring', label: 'Recurring', icon: CalendarClock },
  { to: '/budgets', label: 'Budgets', icon: PiggyBank },
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
      className={`border-r border-border bg-surface transition-all ${sidebarCollapsed ? 'w-20' : 'w-64'
        }`}
    >
      <div className="flex items-center justify-between border-b border-border px-4 py-3 border-border">
        <span className="text-sm font-semibold text-foreground">ledgerA</span>
        <button className="text-xs text-muted" onClick={toggleSidebar}>
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
                ? 'bg-accent text-white'
                : 'text-secondary hover:bg-surface-hover hover:text-foreground text-secondary hover:bg-surface-hover '
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
      <div className="mt-auto border-t border-border p-3 text-xs text-muted border-border text-muted">
        {user ? `${user.displayName ?? user.email}` : 'Guest'}
      </div>
    </aside>
  )
}
