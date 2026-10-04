export const getRating = (fullCombo: boolean, score: number): string => {
  if (score >= 1000000) return 'phi'
  if (fullCombo) return 'FC'
  if (score >= 960000) return 'V'
  if (score >= 920000) return 'S'
  if (score >= 880000) return 'A'
  if (score >= 820000) return 'B'
  if (score >= 700000) return 'C'
  return 'F'
}

// Phigros stores the challenge mode rank as <mode><rank>, e.g. 32 => mode 3, rank 2.
export const challengeMode = (challenge: number): number => {
  let n = challenge
  while (n > 10) n = Math.floor(n / 10)
  return n
}

export const challengeRank = (challenge: number): number => {
  if (challenge >= 100) return challenge % 100
  if (challenge >= 10) return challenge % 10
  return challenge
}

export const formatDate = (raw: string): string => {
  if (!raw) return ''
  return new Date(raw)
    .toLocaleString('zh-CN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
    })
    .replace(/\//g, '-')
    .replace(',', '')
}

export const songRks = (acc: number, rank: number): number => {
  if (acc === 100) return rank
  if (acc < 70) return 0
  const x = (acc - 55) / 45
  return rank * x * x
}

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))