import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { byPrefixAndName } from '../lib/fontawesome'
import './Navbar.css'

interface NavbarProps {
  onMenuClick: () => void
}

function Navbar({ onMenuClick }: NavbarProps) {
  return (
    <header className="navbar">
      <button
        type="button"
        className="navbar__menu-btn"
        onClick={onMenuClick}
        aria-label="Open menu"
      >
        <FontAwesomeIcon icon={byPrefixAndName.fas.bars} />
      </button>

      <div className="navbar__title">
        <span className="navbar__title-text">Grey Matter</span>
      </div>

      <div className="navbar__auth">
        <button type="button" className="navbar__auth-btn" title="Coming soon">
          Log In
        </button>
        <button
          type="button"
          className="navbar__auth-btn navbar__auth-btn--primary"
          title="Coming soon"
        >
          Sign Up
        </button>
      </div>
    </header>
  )
}

export default Navbar
