import type { HistoryEntry } from '../types'

export class ApiError extends Error {}

async function parseError(res: Response): Promise<string> {
  try {
    const body = (await res.json()) as { error?: string }
    return body.error ?? `request failed with status ${res.status}`
  } catch {
    return `request failed with status ${res.status}`
  }
}

export async function calculate(expression: string): Promise<HistoryEntry> {
  const res = await fetch('/api/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ expression }),
  })
  if (!res.ok) {
    throw new ApiError(await parseError(res))
  }
  return (await res.json()) as HistoryEntry
}

export async function fetchHistory(): Promise<HistoryEntry[]> {
  const res = await fetch('/api/history')
  if (!res.ok) {
    throw new ApiError(await parseError(res))
  }
  return (await res.json()) as HistoryEntry[]
}
