import { Link } from '@tanstack/react-router'

export default function MailBoxList({ list }: { list: string[] }) {
  return (
    <ul className="mailbox-ul">
      {list.map((mb, index) => (
        <li key={index}>
          <Link to={`/mb/$mbName`} params={{ mbName: mb }}>
            {mb}
          </Link>
        </li>
      ))}
    </ul>
  )
}
