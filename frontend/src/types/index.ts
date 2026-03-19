export type TransactionType = 'income' | 'expense'
export type EntityType = 'account' | 'transaction' | 'category' | 'subcategory' | 'quick_transaction'
export type AuditAction = 'create' | 'update' | 'delete' | 'restore'
export type Period = 'day' | 'week' | 'month' | 'year'

export interface User {
    id: string
    firebase_uid: string
    email: string
    display_name: string
    currency_code: string
    created_at: string
    updated_at: string
}

export interface Account {
    id: string
    user_id: string
    name: string
    account_type: string
    opening_balance: number
    current_balance: number
    is_archived: boolean
    created_at: string
    updated_at: string
    last_transaction_date?: string
}

export interface Subcategory {
    id: string
    category_id: string
    user_id: string
    name: string
    created_at: string
    updated_at: string
}

export interface Category {
    id: string
    user_id: string
    name: string
    created_at: string
    updated_at: string
    subcategories?: Subcategory[]
}

export interface Transaction {
    id: string
    user_id: string
    account_id: string
    category_id: string
    subcategory_id: string
    name: string
    amount: number
    transaction_date: string
    notes?: string
    is_scheduled: boolean
    created_at: string
    updated_at: string
}

export interface QuickTransaction {
    id: string
    user_id: string
    label: string
    account_id?: string
    category_id?: string
    subcategory_id?: string
    amount?: number
    notes?: string
    sort_order: number
    created_at: string
    updated_at: string
}

export interface AuditLog {
    id: string
    user_id: string
    entity_type: EntityType
    entity_id: string
    action: AuditAction
    diff: {
        before: unknown
        after: unknown
    }
    created_at: string
}

export interface TimeseriesPoint {
    label: string
    income: number
    expense: number
}

export interface CategoryBreakdownItem {
    category: string
    subcategory: string
    amount: number
    percentage: number
}

export interface StatsSummary {
    total_income: number
    total_expense: number
    net: number
    category_breakdown_expense: CategoryBreakdownItem[]
    category_breakdown_income: CategoryBreakdownItem[]
    timeseries: TimeseriesPoint[]
}

export interface PaginationMeta {
    page: number
    per_page: number
    total: number
}

export interface PaginatedResponse<T> {
    data: T
    meta: PaginationMeta
}

export interface ApiError {
    error: {
        code: string
        message: string
    }
}

export interface TransactionFilters {
    account_id?: string
    category_id?: string
    subcategory_id?: string
    date_from?: string
    date_to?: string
    search?: string
    type?: 'income' | 'expense' | 'all'
    sort_by?: 'transaction_date' | 'amount' | 'name' | 'created_at'
    sort_dir?: 'asc' | 'desc'
    page?: number
    per_page?: number
    passbook_mode?: boolean
}

export interface StatsFilters {
    period: Period
    value: string
    account_id?: string
}

export interface ComparePeriodData {
    label: string
    total_income: number
    total_expense: number
    net: number
    top_expense: CategoryBreakdownItem[]
    top_income: CategoryBreakdownItem[]
}

export interface CompareResponse {
    period1: ComparePeriodData
    period2: ComparePeriodData
    income_change_pct: number
    expense_change_pct: number
    net_change: number
}

export interface CompareFilters {
    period: Period
    value1: string
    value2: string
    account_id?: string
}

export interface RecurringTransaction {
    id: string
    user_id: string
    account_id: string
    category_id: string
    subcategory_id: string
    name: string
    amount: number
    notes?: string
    frequency: 'daily' | 'weekly' | 'monthly' | 'yearly'
    start_date: string
    next_due_date: string
    end_date?: string
    is_active: boolean
    last_executed_at?: string
    created_at: string
    updated_at: string
}

export interface Budget {
    id: string
    user_id: string
    category_id: string
    amount: number
    period: 'monthly' | 'yearly'
    is_active: boolean
    created_at: string
    updated_at: string
}

export interface BudgetProgress extends Budget {
    spent: number
    remaining: number
    percent: number
}
