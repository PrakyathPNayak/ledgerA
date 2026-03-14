import { useQuery } from '@tanstack/react-query'

import api from '@/lib/api'
import type { StatsFilters, StatsSummary } from '@/types'

export const statsQueryKeys = {
    all: ['stats'] as const,
    detail: (filters: StatsFilters) => ['stats', filters] as const,
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
