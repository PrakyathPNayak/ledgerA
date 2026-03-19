import { useQuery } from '@tanstack/react-query'

import api from '@/lib/api'
import type { CompareFilters, CompareResponse, MonthlyFilters, MonthlyReport, StatsFilters, StatsSummary } from '@/types'

export const statsQueryKeys = {
    all: ['stats'] as const,
    detail: (filters: StatsFilters) => ['stats', filters] as const,
    compare: (filters: CompareFilters) => ['stats', 'compare', filters] as const,
    monthly: (filters: MonthlyFilters) => ['stats', 'monthly', filters] as const,
}

/**
 * @description Fetches summary stats for a period.
 */
export function useStats(filters: StatsFilters) {
    const hasValidFilters = Boolean(filters.period && filters.value)

    return useQuery<StatsSummary>({
        queryKey: statsQueryKeys.detail(filters),
        queryFn: () => api.get('/stats/summary', { params: filters }),
        enabled: hasValidFilters,
        staleTime: 60 * 1000,
    })
}

/**
 * @description Fetches comparison between two periods.
 */
export function useCompare(filters: CompareFilters) {
    const hasValidFilters = Boolean(filters.period && filters.value1 && filters.value2)

    return useQuery<CompareResponse>({
        queryKey: statsQueryKeys.compare(filters),
        queryFn: () => api.get('/stats/compare', { params: filters }),
        enabled: hasValidFilters,
        staleTime: 60 * 1000,
    })
}

/**
 * @description Fetches monthly income/expense report.
 */
export function useMonthlyReport(filters: MonthlyFilters) {
    return useQuery<MonthlyReport>({
        queryKey: statsQueryKeys.monthly(filters),
        queryFn: () => api.get('/stats/monthly', { params: filters }),
        staleTime: 5 * 60 * 1000,
    })
}
