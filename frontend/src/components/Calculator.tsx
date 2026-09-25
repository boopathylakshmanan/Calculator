import { useEffect, useState } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { byPrefixAndName } from '../lib/fontawesome'
import { ApiError, calculate, fetchHistory } from '../api/client'
import type { HistoryEntry } from '../types'
import HistoryPanel from './HistoryPanel'
import './Calculator.css'

type ButtonKind = 'digit' | 'operator' | 'action' | 'equals'

interface CalcButton {
  label: string
  token?: string
  kind: ButtonKind
  ariaLabel?: string
}

const ROWS: CalcButton[][] = [
  [
    { label: '(', kind: 'digit' },
    { label: ')', kind: 'digit' },
    { label: 'AC', kind: 'action' },
    { label: 'DEL', kind: 'action' },
  ],
  [
    { label: '7', kind: 'digit' },
    { label: '8', kind: 'digit' },
    { label: '9', kind: 'digit' },
    { label: '÷', token: '/', kind: 'operator' },
  ],
  [
    { label: '4', kind: 'digit' },
    { label: '5', kind: 'digit' },
    { label: '6', kind: 'digit' },
    { label: '×', token: '*', kind: 'operator' },
  ],
  [
    { label: '1', kind: 'digit' },
    { label: '2', kind: 'digit' },
    { label: '3', kind: 'digit' },
    { label: '−', token: '-', kind: 'operator' },
  ],
  [
    { label: '0', kind: 'digit' },
    { label: '.', kind: 'digit' },
    { label: '+', kind: 'operator' },
    { label: '=', kind: 'equals' },
  ],
]

const RESETS_ON_FRESH_INPUT = new Set(['(', '.', ...'0123456789'.split('')])

function Calculator() {
  const [expression, setExpression] = useState('')
  const [prevExpression, setPrevExpression] = useState('')
  const [justEvaluated, setJustEvaluated] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const [history, setHistory] = useState<HistoryEntry[]>([])
  const [historyLoading, setHistoryLoading] = useState(true)
  const [historyError, setHistoryError] = useState<string | null>(null)
  const [historyOpen, setHistoryOpen] = useState(false)

  const loadHistory = () => {
    setHistoryLoading(true)
    setHistoryError(null)
    fetchHistory()
      .then(setHistory)
      .catch((e: unknown) => {
        setHistoryError(e instanceof ApiError ? e.message : 'Could not load history.')
      })
      .finally(() => setHistoryLoading(false))
  }

  useEffect(() => {
    fetchHistory()
      .then(setHistory)
      .catch((e: unknown) => {
        setHistoryError(e instanceof ApiError ? e.message : 'Could not load history.')
      })
      .finally(() => setHistoryLoading(false))
  }, [])

  const pressToken = (label: string, token: string | undefined) => {
    setError(null)
    const value = token ?? label
    if (justEvaluated) {
      setJustEvaluated(false)
      setExpression(RESETS_ON_FRESH_INPUT.has(label) ? value : expression + value)
      return
    }
    setExpression((prev) => prev + value)
  }

  const onAllClear = () => {
    setExpression('')
    setPrevExpression('')
    setError(null)
    setJustEvaluated(false)
  }

  const onDelete = () => {
    setError(null)
    setJustEvaluated(false)
    setExpression((prev) => prev.slice(0, -1))
  }

  const onEquals = async () => {
    if (!expression.trim() || busy) return
    setBusy(true)
    setError(null)
    try {
      const entry = await calculate(expression)
      setPrevExpression(expression)
      setExpression(entry.result)
      setJustEvaluated(true)
      setHistory((prev) => [entry, ...prev].slice(0, 100))
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Something went wrong.')
    } finally {
      setBusy(false)
    }
  }

  const onPress = (btn: CalcButton) => {
    switch (btn.kind) {
      case 'action':
        if (btn.label === 'AC') onAllClear()
        else onDelete()
        return
      case 'equals':
        void onEquals()
        return
      default:
        pressToken(btn.label, btn.token)
    }
  }

  const onSelectHistory = (entry: HistoryEntry) => {
    setExpression(entry.expression)
    setPrevExpression('')
    setJustEvaluated(false)
    setError(null)
    setHistoryOpen(false)
  }

  return (
    <div className="calculator" id="calculator">
      <div className="calculator__display">
        <button
          type="button"
          className="calculator__history-btn"
          onClick={() => setHistoryOpen(true)}
          aria-label="Open calculation history"
          title="History"
        >
          <FontAwesomeIcon icon={byPrefixAndName.fas['clock-rotate-left']} />
        </button>

        <div className="calculator__display-lines">
          <div className={`calculator__line-top ${error ? 'calculator__line-top--error' : ''}`}>
            {error ? error : justEvaluated ? `${prevExpression} =` : ' '}
          </div>
          <div className="calculator__line-main" role="status" aria-live="polite">
            {expression || '0'}
          </div>
        </div>
      </div>

      <div className="calculator__pad">
        {ROWS.flat().map((btn, i) => (
          <button
            key={`${btn.label}-${i}`}
            type="button"
            className={`calc-btn calc-btn--${btn.kind}`}
            onClick={() => onPress(btn)}
            aria-label={btn.ariaLabel ?? btn.label}
            disabled={btn.kind === 'equals' && busy}
          >
            {btn.label}
          </button>
        ))}
      </div>

      <HistoryPanel
        open={historyOpen}
        entries={history}
        loading={historyLoading}
        error={historyError}
        onClose={() => setHistoryOpen(false)}
        onSelect={onSelectHistory}
        onRetry={loadHistory}
      />
    </div>
  )
}

export default Calculator
