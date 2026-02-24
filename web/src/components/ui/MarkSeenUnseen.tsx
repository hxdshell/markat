import { useState } from 'react'
import IconBtn from './IconBtn'
import { Eye, EyeOff } from 'lucide-react'
import { markSeenUnseen } from '../../api/message'

export default function MarkSeenUnseen({
  mb,
  uid,
}: {
  mb: string
  uid: number
}) {
  const [seen, setSeen] = useState(true)
  const seenValue = seen ? 'unseen' : 'seen'

  async function mark(mb: string, uid: number) {
    const uids = new Array()
    uids.push(uid)
    const res = await markSeenUnseen(mb, uids, !seen)
    if (res.status == 200) {
      setSeen((seen) => !seen)
    }
  }

  return (
    <IconBtn abbr={`Mark as ${seenValue}`} onClick={() => mark(mb, uid)}>
      {seen ? <EyeOff /> : <Eye />}
    </IconBtn>
  )
}
