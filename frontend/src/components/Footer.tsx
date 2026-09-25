import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { byPrefixAndName } from '../lib/fontawesome'
import './Footer.css'

const COLUMNS: { title: string; id: string; links: string[] }[] = [
  { title: 'About', id: 'about', links: ['Our Story', 'Team', 'Careers'] },
  { title: 'Contact', id: 'contact', links: ['Support', 'Sales', 'Feedback'] },
  { title: 'Products', id: 'products', links: ['Calculator', 'History Sync', 'API'] },
  { title: 'Services', id: 'services', links: ['Enterprise', 'Integrations', 'Status'] },
]

const SOCIALS = [
  { label: 'GitHub', href: 'https://github.com/boopathylakshmanan', icon: byPrefixAndName.fab.github },
  {
    label: 'LinkedIn',
    href: 'https://linkedin.com/in/boopathylakshmanan',
    icon: byPrefixAndName.fab['square-linkedin'],
  },
  { label: 'YouTube', href: 'https://youtube.com', icon: byPrefixAndName.fab.youtube },
  { label: 'Instagram', href: 'https://instagram.com/boopathy_lakshmanan', icon: byPrefixAndName.fab.instagram },
]

function Footer() {
  return (
    <footer className="footer">
      <div className="footer__columns">
        {COLUMNS.map((col) => (
          <div className="footer__col" id={col.id} key={col.title}>
            <h3>{col.title}</h3>
            <ul>
              {col.links.map((link) => (
                <li key={link}>
                  <a href="#">{link}</a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>

      <div className="footer__socials">
        {SOCIALS.map((s) => (
          <a
            key={s.label}
            href={s.href}
            target="_blank"
            rel="noreferrer"
            aria-label={s.label}
            className="footer__social-link"
          >
            <FontAwesomeIcon icon={s.icon} />
          </a>
        ))}
      </div>

      <div className="footer__bottom">
        <span>&trade; Grey Matter &mdash; 2026. All Rights Reserved.</span>
      </div>
    </footer>
  )
}

export default Footer
