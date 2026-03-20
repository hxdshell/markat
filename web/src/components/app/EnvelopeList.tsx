import { ChevronLeft, ChevronRight, RotateCw } from 'lucide-react'
import Loading from '../ui/Loading'
import { useNavigate } from '@tanstack/react-router'
import EnvelopesToolbar from './EnvelopesToolbar'
import { useMemo, useState } from 'react'

export default function EnvelopeList({
  mbName,
  data,
  loading,
  refresh,
  prev,
  next,
}: {
  mbName: string
  data: EnvelopeResponse
  loading: boolean
  refresh: () => Promise<void>
  prev: (data: EnvelopeResponse) => Promise<void>
  next: (data: EnvelopeResponse) => Promise<void>
}) {
  const navigate = useNavigate()

  const [selectedIds, setSelectedIds] = useState(new Set<number>())

  const currentPageIds = useMemo(
    () => data.envelopes.map((envlp) => envlp.uid),
    [data],
  )

  const allIdsSelected = selectedIds.size === currentPageIds.length
  const someIdsSelected = selectedIds.size > 0 && !allIdsSelected

  function toggleRowSelect(value: boolean, id: number) {
    if (value) {
      const newSet = new Set<number>([...selectedIds, id])
      setSelectedIds(newSet)
    } else {
      const newSet = new Set<number>(selectedIds)
      if (newSet.delete(id)) {
        setSelectedIds(newSet)
      }
    }
  }

  function toggleAll(value: boolean) {
    if (value) {
      const newSet = new Set<number>([...currentPageIds])
      setSelectedIds(newSet)
    } else {
      const newSet = new Set<number>()
      setSelectedIds(newSet)
    }
  }

  function openMessage(uid: number) {
    navigate({ to: `/message/${mbName}/${uid}` })
  }
  return (
    <div className="envelope-container">
      <div className="navbar">
        <div>
          <button className="icon-btn" onClick={refresh}>
            <abbr title="Refresh">
              <RotateCw />
            </abbr>
          </button>
        </div>
        {loading ? <Loading /> : null}
        <div className="pagination-control">
          <button
            disabled={data.page <= 1}
            className="icon-btn"
            onClick={() => prev(data)}
          >
            <abbr title="Prev">
              <ChevronLeft className={data.page <= 1 ? 'disabled-icon' : ''} />
            </abbr>
          </button>
          {data.start}-{data.end} of {data.total}
          <button
            disabled={data.end === data.total}
            className="icon-btn"
            onClick={() => next(data)}
          >
            <abbr title="Next">
              <ChevronRight
                className={data.end === data.total ? 'disabled-icon' : ''}
              />
            </abbr>
          </button>
        </div>
      </div>
      <EnvelopesToolbar
        allIdsSelected={allIdsSelected}
        someIdsSelected={someIdsSelected}
        mb={mbName}
        refresh={refresh}
        setSelectedIds={setSelectedIds}
        selectedIds={selectedIds}
        toggleAll={toggleAll}
      />
      {data.envelopes.length > 0 ? (
        <table className="envelope-list">
          <tbody>
            {data.envelopes.map((envlp: Envelope) => (
              <tr
                key={envlp.uid}
                className={
                  envlp.flags.includes('\\Seen') ? 'envelope' : 'envelope bold'
                }
                onClick={() => openMessage(envlp.uid)}
              >
                <td className="checkbox">
                  <input
                    type="checkbox"
                    name={`chkbx-${envlp.uid}`}
                    id={`chkbx-${envlp.uid}`}
                    checked={selectedIds.has(envlp.uid)}
                    onClick={(e) => {
                      e.stopPropagation()
                    }}
                    onChange={(e) =>
                      toggleRowSelect(e.target.checked, envlp.uid)
                    }
                  />
                </td>
                <td className="from">{envlp.from}</td>
                <td className="subject">{envlp.subject}</td>
                <td className="date">{envlp.date}</td>
              </tr>
            ))}
          </tbody>
        </table>
      ) : (
        <div className="flex mt-1 ml-1">There are no mail in this folder</div>
      )}
    </div>
  )
}
