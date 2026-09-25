import { useEffect } from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { byPrefixAndName } from '../lib/fontawesome'
import './Drawer.css'

interface DrawerProps {
  open: boolean
  onClose: () => void
}

const LINKS = [
  { label: 'Calculator', href: '#calculator' },
  { label: 'About', href: '#about' },
  { label: 'Contact', href: '#contact' },
  { label: 'Products', href: '#products' },
  { label: 'Services', href: '#services' },
]

function Drawer({ open, onClose }: DrawerProps) {
  useEffect(() => {
    if (!open) return
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [open, onClose])

  return (
    <div className={`drawer-root ${open ? 'drawer-root--open' : ''}`} aria-hidden={!open}>
      <div className="drawer-overlay" onClick={onClose} />
      <nav className="drawer" aria-label="Main menu">
        <div className="drawer__header">
          <span className="drawer__brand">Grey Matter</span>
          <button
            type="button"
            className="drawer__close"
            onClick={onClose}
            aria-label="Close menu"
          >
            <FontAwesomeIcon icon={byPrefixAndName.fas.xmark} />
          </button>
        </div>
        <ul className="drawer__links">
          {LINKS.map((link) => (
            <li key={link.label}>
              <a href={link.href} onClick={onClose}>
                {link.label}
              </a>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  )
}

export default Drawer
