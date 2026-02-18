import { useNavigate } from '@tanstack/react-router'
import IconBtn from '../ui/IconBtn'
import {
  ArrowLeft,
  FileCodeCorner,
  FolderDown,
  FolderUp,
  Trash2,
} from 'lucide-react'
import MarkSeenUnseen from '../ui/MarkSeenUnseen'
import { moveToFolder } from '../../api/message'

export default function MessageToolbar({
  mb,
  uid,
}: {
  mb: string
  uid: number
}) {
  const navigate = useNavigate()
  const mbGeneric = mb.toLowerCase()

  async function move(mb: string, uid: number, dest: string) {
    const res = await moveToFolder(mb, uid, dest)
    if (res.status === 200) {
      navigate({ to: `/mb/${mb}` })
    }
  }
  return (
    <div className="toolbar">
      <IconBtn
        abbr="Go back"
        onClick={() => {
          navigate({ to: `/mb/${mb}` })
        }}
      >
        <ArrowLeft />
      </IconBtn>
      {mbGeneric === 'inbox' ? (
        <IconBtn abbr="Archive" onClick={() => move(mb, uid, 'Archive')}>
          <FolderDown />
        </IconBtn>
      ) : (
        <IconBtn abbr="Move to Inbox" onClick={() => move(mb, uid, 'INBOX')}>
          <FolderUp />
        </IconBtn>
      )}

      <MarkSeenUnseen mb={mb} uid={uid} />
      {mbGeneric !== 'trash' ? (
        <IconBtn abbr="Move to trash" onClick={() => move(mb, uid, 'Trash')}>
          <Trash2 />
        </IconBtn>
      ) : null}

      <IconBtn abbr="Read original message" onClick={() => {}}>
        <FileCodeCorner />
      </IconBtn>
    </div>
  )
}
