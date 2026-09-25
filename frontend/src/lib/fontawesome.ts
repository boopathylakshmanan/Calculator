import type { IconDefinition } from '@fortawesome/fontawesome-svg-core'
import {
  faGithub,
  faInstagram,
  faSquareLinkedin,
  faYoutube,
} from '@fortawesome/free-brands-svg-icons'
import {
  faBars,
  faXmark,
  faClockRotateLeft,
  faDeleteLeft,
  faTrashCan,
} from '@fortawesome/free-solid-svg-icons'

const fab = {
  github: faGithub,
  instagram: faInstagram,
  'square-linkedin': faSquareLinkedin,
  youtube: faYoutube,
}

const fas = {
  bars: faBars,
  xmark: faXmark,
  'clock-rotate-left': faClockRotateLeft,
  'delete-left': faDeleteLeft,
  'trash-can': faTrashCan,
}

// Small local stand-in for a Font Awesome Kit's generated lookup, built from
// the free npm icon packages so `byPrefixAndName.fab['github']`-style access
// works without a hosted Kit script.
export const byPrefixAndName: {
  fab: Record<keyof typeof fab, IconDefinition>
  fas: Record<keyof typeof fas, IconDefinition>
} = { fab, fas }
