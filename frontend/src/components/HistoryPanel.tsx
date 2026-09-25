import { useEffect } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { byPrefixAndName } from '../lib/fontawesome'
import type { HistoryEntry } from '../types'
import './HistoryPanel.css'

interface HistoryPanelProps {
  open: boolean
  entries: HistoryEntry[]
  loading: boolean
  error: string | null
  onClose: () => void
  onSelect: (entry: HistoryEntry) => void
  onRetry: () => void
}

function formatTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function HistoryPanel({ open, entries, loading, error, onClose, onSelect, onRetry }: HistoryPanelProps) {
  useEffect(() => {
    if (!open) return
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [open, onClose])

  return (
    <div className={`history-root ${open ? 'history-root--open' : ''}`} aria-hidden={!open}>
      <div className="history-overlay" onClick={onClose} />
      <div className="history-modal" role="dialog" aria-label="Calculation history" aria-modal="true">
        <div className="history-modal__header">
          <span>
            <FontAwesomeIcon icon={byPrefixAndName.fas['clock-rotate-left']} /> History
          </span>
          <button
            type="button"
            className="history-modal__close"
            onClick={onClose}
            aria-label="Close history"
          >
            <FontAwesomeIcon icon={byPrefixAndName.fas.xmark} />
          </button>
        </div>

        <div className="history-modal__body glow-scroll">
          {loading && <p className="history-modal__note">Loading history…</p>}

          {!loading && error && (
            <div className="history-modal__note history-modal__note--error">
              <p>{error}</p>
              <button type="button" className="history-modal__retry" onClick={onRetry}>
                Retry
              </button>
            </div>
          )}

          {!loading && !error && entries.length === 0 && (
            <p className="history-modal__note">No calculations yet.</p>
          )}

          {!loading && !error && entries.length > 0 && (
            <ul className="history-list">
              {entries.map((entry) => (
                <li key={entry.id}>
                  <button
                    type="button"
                    className="history-list__item"
                    onClick={() => onSelect(entry)}
                  >
                    <span className="history-list__expr">{entry.expression}</span>
                    <span className="history-list__result">= {entry.result}</span>
                    <span className="history-list__time">{formatTime(entry.created_at)}</span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}

export default HistoryPanel
