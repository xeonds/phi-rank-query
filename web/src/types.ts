export interface SongRecord {
  id: string
  song: string
  level: string
  difficulty: string
  rks: number
  score: number
  acc: number
  fullCombo: boolean
  illustration: string
}

export interface PlayerResult {
  player: string
  rks: number
  challenge: number
  date: string
  records: SongRecord[]
  phi: SongRecord[]
}

export interface RankItem {
  id: string
  title: string
  EZ: string
  HD: string
  IN: string
  AT: string
  Legacy: string
}

export interface LeaderboardItem {
  id: number
  username: string
  rks: number
}

export interface Session {
  token: string
  alias: string
}