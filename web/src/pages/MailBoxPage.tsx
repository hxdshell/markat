import { ChevronLeft, ChevronRight, RotateCw } from 'lucide-react'
import { mbNameRoute } from '../router'

export default function MailBoxPage() {
  const loaderData: ApiResponseType = mbNameRoute.useLoaderData()
  const navigate = mbNameRoute.useNavigate()
  const data: EnvelopeResponse = loaderData.data

  async function refersh() {
    navigate({ search: (s: any) => ({ ...s, page: 1 }) })
  }

  async function prev() {
    if (data.page > 1) {
      navigate({ search: (s: any) => ({ ...s, page: data.page - 1 }) })
    }
  }
  async function next() {
    if (data.end !== data.total) {
      navigate({ search: (s: any) => ({ ...s, page: data.page + 1 }) })
    }
  }

  if (loaderData.status != 200) {
    return (
      <>
        <p>{loaderData.message}</p>
      </>
    )
  } else {
    const envelopes: Envelope[] = data.envelopes
    if (envelopes.length === 0) {
      return <div>There are no mail in this mailbox</div>
    }
    return (
      <div className="envelope-container">
        <div className="toolbar">
          <div>
            <button className="icon-btn" onClick={refersh}>
              <RotateCw />
            </button>
          </div>
          <div>
            <button
              disabled={data.page === 1}
              className="icon-btn"
              onClick={prev}
            >
              <ChevronLeft className={data.page === 1 ? 'disabled-icon' : ''} />
            </button>
            {data.start}-{data.end} of {data.total}
            <button
              disabled={data.end === data.total}
              className="icon-btn"
              onClick={next}
            >
              <ChevronRight
                className={data.end === data.total ? 'disabled-icon' : ''}
              />
            </button>
          </div>
        </div>

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
      </div>
    )
  }
}
