import { ChevronLeft, ChevronRight, RotateCw } from 'lucide-react'
import Loading from '../ui/Loading'

export default function EnvelopeList({
  data,
  loading,
  refresh,
  prev,
  next,
}: {
  data: EnvelopeResponse
  loading: boolean
  refresh: () => Promise<void>
  prev: (data: EnvelopeResponse) => Promise<void>
  next: (data: EnvelopeResponse) => Promise<void>
}) {
  return (
    <div className="envelope-container">
      <div className="toolbar">
        <div>
          <button className="icon-btn" onClick={refresh}>
            <abbr title="Refresh">
              <RotateCw />
            </abbr>
          </button>
        </div>
        {loading ? <Loading /> : null}
        <div>
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
      {data.envelopes.length > 0 ? (
        <table className="envelope-list">
          <tbody>
            {data.envelopes.map((envlp: Envelope) => (
              <tr
                key={envlp.uid}
                className={
                  envlp.flags.includes('\\Seen') ? 'envelope' : 'envelope bold'
                }
              >
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
