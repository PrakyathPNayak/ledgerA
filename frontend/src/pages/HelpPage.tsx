import { PageShell } from '@/components/layout/PageShell'
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from '@/components/ui/accordion'

const sections = [
    {
        id: 'getting-started',
        title: 'Getting Started',
        body: 'Create at least one account (e.g. Savings, Wallet) before adding transactions. Navigate to Accounts → Add Account. Then set up Categories to organise your spending. Once you have an account and at least one category, go to Dashboard and press the + button to record your first transaction.',
    },
    {
        id: 'adding-transactions',
        title: 'Adding Transactions',
        body: 'Tap the + FAB on Dashboard or go to the Transactions page. Choose Income for money coming in (salary, freelance) and Expense for money going out (groceries, bills). Enter the amount as a positive number — the app stores expenses as negative values internally. Pick an account, category, optional subcategory, and a date. Add a description for searchability.',
    },
    {
        id: 'managing-accounts',
        title: 'Managing Accounts',
        body: 'Each account tracks its own running balance. Open an account from the Accounts page to see a passbook-style chronological view of every transaction affecting that account. You can create, rename, or delete accounts. Deleting an account removes all associated transactions — this action is irreversible.',
    },
    {
        id: 'stats-charts',
        title: 'Stats & Charts',
        body: 'The Stats page shows income vs expense breakdowns with a bar chart and category doughnut for the selected period. Use the period selector to switch between Day, Week, Month, and Year views. Previous/Next navigation steps through time. Press Export PDF to download a shareable snapshot of the current view, including totals and category breakdown.',
    },
    {
        id: 'searching-transactions',
        title: 'Searching Transactions',
        body: 'The Search page combines full-text search with filters for account, category, transaction type (income/expense), and date range. Results are paginated — use Load More to fetch additional rows without losing what is already visible. All active filters are shown as dismissible chips above the results list.',
    },
    {
        id: 'quick-transactions',
        title: 'Quick Transactions',
        body: 'Quick Transactions are pre-configured transaction templates for recurring spends like rent, subscriptions, or fuel. Create one from the Quick Transactions page by filling in a label, amount, account, and category. Use the Execute button (▶) to instantly record that transaction for today. Drag the handle icon to reorder your quick transaction tiles.',
    },
    {
        id: 'data-privacy',
        title: 'Data & Privacy',
        body: 'All data is stored in a PostgreSQL database scoped to your Firebase account. ledgerA does not share your financial data with third parties. The app works offline for viewing cached data — mutations (creating or editing records) require an active internet connection. You can export a PDF report at any time from the Stats page.',
    },
]

export function HelpPage() {
    return (
        <PageShell title="Help">
            <div className="space-y-4">
                <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
                    <h2 className="text-2xl font-bold tracking-tight text-slate-900">ledgerA Help Center</h2>
                    <p className="mt-2 text-sm leading-6 text-slate-600">
                        Find quick answers about accounts, transactions, reports, and productivity workflows.
                    </p>
                </section>

                <section className="rounded-2xl border border-slate-200 bg-white px-6 shadow-sm">
                    <Accordion type="multiple">
                        {sections.map((section) => (
                            <AccordionItem key={section.id} value={section.id}>
                                <AccordionTrigger>{section.title}</AccordionTrigger>
                                <AccordionContent>{section.body}</AccordionContent>
                            </AccordionItem>
                        ))}
                    </Accordion>
                </section>
            </div>
        </PageShell>
    )
}

