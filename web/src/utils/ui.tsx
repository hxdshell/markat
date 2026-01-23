import {
  Archive,
  File,
  Folder,
  Inbox,
  Send,
  ShieldAlert,
  Trash,
} from 'lucide-react'

export function mapIcons(name: string) {
  const k = name.toLowerCase()
  switch (k) {
    case 'archive':
      return <Archive />
    case 'drafts':
      return <File />
    case 'inbox':
      return <Inbox />
    case 'sent':
      return <Send />
    case 'spam':
      return <ShieldAlert />
    case 'trash':
      return <Trash />
    default:
      return <Folder />
  }
}
