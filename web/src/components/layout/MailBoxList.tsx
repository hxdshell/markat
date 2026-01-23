import { Link } from '@tanstack/react-router'
import { mapIcons } from '../../utils/ui'

export default function MailBoxList({ list }: { list: string[] }) {
  return (
    <ul className="mailbox-ul">
      {list.map((mb, index) => (
        <li key={index}>
          <Link to={`/mb/$mbName`} params={{ mbName: mb }}>
            <span>{mapIcons(mb)}</span>
            <span>{mb}</span>
          </Link>
        </li>
      ))}
    </ul>
  )
}
