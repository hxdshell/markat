import IconBtn from '../ui/IconBtn'
import { FolderDown, FolderUp, Trash2 } from 'lucide-react'
import { moveToFolder } from '../../api/message'

export default function EnvelopesToolbar({
  mb,
  selectedIds,
  allIdsSelected,
  someIdsSelected,
  setSelectedIds,
  refresh,
  toggleAll,
}: {
  mb: string
  selectedIds: Set<number>
  allIdsSelected: boolean
  someIdsSelected: boolean
  toggleAll: (val: boolean) => void
  refresh: () => Promise<void>
  setSelectedIds: React.Dispatch<React.SetStateAction<Set<number>>>
}) {
  const mbGeneric = mb.toLowerCase()

  async function move(mb: string, selectedIds: Set<number>, dest: string) {
    if (selectedIds.size > 0) {
      const uids = Array.from(selectedIds)
      const res = await moveToFolder(mb, uids, dest)
      if (res.status === 200) {
        setSelectedIds(new Set<number>())
        refresh()
      }
    }
  }
  return (
    <div className="toolbar">
      <input
        type="checkbox"
        name="select-all"
        id="select-all"
        checked={allIdsSelected}
        onChange={(e) => {
          e.stopPropagation()
          toggleAll(e.target.checked)
        }}
      />
      {someIdsSelected || allIdsSelected ? (
        <>
          {mbGeneric === 'inbox' ? (
            <IconBtn
              abbr="Archive"
              onClick={() => move(mb, selectedIds, 'Archive')}
            >
              <FolderDown />
            </IconBtn>
          ) : (
            <IconBtn
              abbr="Move to Inbox"
              onClick={() => move(mb, selectedIds, 'INBOX')}
            >
              <FolderUp />
            </IconBtn>
          )}

          {/* <MarkSeenUnseen mb={mb} uid={uid} /> */}
          {mbGeneric !== 'trash' ? (
            <IconBtn
              abbr="Move to trash"
              onClick={() => move(mb, selectedIds, 'Trash')}
            >
              <Trash2 />
            </IconBtn>
          ) : null}
        </>
      ) : null}
    </div>
  )
}
