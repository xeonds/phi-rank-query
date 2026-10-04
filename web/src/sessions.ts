import { reactive, watch } from 'vue'
import type { PlayerResult, Session } from './types'

const SESSIONS_KEY = 'sessions'
const SELECTED_KEY = 'selectedSession'

function loadSessions(): Session[] {
  const raw = localStorage.getItem(SESSIONS_KEY)
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) return []
    const aliases: string[] = JSON.parse(localStorage.getItem('aliases') || '[]')
    return parsed.map((s: unknown, i: number) =>
      typeof s === 'string' ? { token: s, alias: aliases[i] || s } : (s as Session),
    )
  } catch {
    return []
  }
}

export const sessionStore = reactive({
  list: loadSessions(),
  selected: localStorage.getItem(SELECTED_KEY) || '',
})

watch(() => sessionStore.list, (v) => localStorage.setItem(SESSIONS_KEY, JSON.stringify(v)), { deep: true })
watch(() => sessionStore.selected, (v) => { if (v) localStorage.setItem(SELECTED_KEY, v) })

export function addSession(token: string, alias: string) {
  sessionStore.list.push({ token, alias: alias || token })
}

export function removeSession(index: number) {
  const removed = sessionStore.list[index]
  sessionStore.list.splice(index, 1)
  if (removed && sessionStore.selected === removed.token) {
    sessionStore.selected = sessionStore.list[0]?.token ?? ''
  }
}

export function selectSession(token: string) {
  sessionStore.selected = token
}

// --- query history (per session, kept in localStorage) ---

export function loadHistory(token: string): PlayerResult[] {
  const raw = localStorage.getItem('history-' + token)
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed : (parsed.data ?? [])
  } catch {
    return []
  }
}

export function saveHistory(token: string, data: PlayerResult[]) {
  localStorage.setItem('history-' + token, JSON.stringify({ timestamp: Date.now(), data }))
}

export function appendHistory(token: string, result: PlayerResult) {
  const data = loadHistory(token).filter((r) => r.date !== result.date)
  data.push(result)
  saveHistory(token, data.slice(-32))
}

export function removeHistory(token: string, index: number) {
  const data = loadHistory(token)
  data.splice(index, 1)
  saveHistory(token, data)
}