import type { LeaderboardItem, PlayerResult, RankItem } from './types'

async function request<T>(input: string, init?: RequestInit): Promise<T> {
  const resp = await fetch(input, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  })
  const data = await resp.json()
  if (!resp.ok) {
    throw new Error(data?.error || `请求失败 (${resp.status})`)
  }
  return data as T
}

export const fetchPlayer = (session: string): Promise<PlayerResult> =>
  request<PlayerResult>('/api/v1/player', {
    method: 'POST',
    body: JSON.stringify({ session }),
  })

export const fetchRankTable = (): Promise<RankItem[]> =>
  request<RankItem[]>('/api/v1/rank_table')

export const fetchLeaderboard = (): Promise<LeaderboardItem[]> =>
  request<LeaderboardItem[]>('/api/v1/leaderboard')