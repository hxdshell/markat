import { ChevronLeft, ChevronRight, RotateCw } from 'lucide-react'
import { mbNameRoute } from '../router'
import { useState } from 'react'
import { fetchEnvelopes } from '../api/mailbox'
import Loading from '../components/ui/Loading'

export default function MailBoxPage() {
  const loaderData: ApiResponseType = mbNameRoute.useLoaderData()
  const envelopeData: EnvelopeResponse = loaderData.data

  const [data, setData] = useState(envelopeData)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function refersh() {
    try {
      setLoading(true)
      const resp = await fetchEnvelopes(1)
      if (resp.status !== 200) {
        setError(resp.message)
      } else {
        setData(resp.data)
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('unkown error, please debug')
      }
    } finally {
      setLoading(false)
    }
  }

  async function prev() {
    try {
      setLoading(true)
      const resp = await fetchEnvelopes(data.page - 1)
      if (resp.status !== 200) {
        setError(resp.message)
      } else {
        setData(resp.data)
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('unkown error, please debug')
      }
    } finally {
      setLoading(false)
    }
  }
  async function next() {
    try {
      setLoading(true)
      const resp = await fetchEnvelopes(data.page + 1)
      if (resp.status !== 200) {
        setError(resp.message)
      } else {
        setData(resp.data)
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('unkown error, please debug')
      }
    } finally {
      setLoading(false)
    }
  }

  if (loaderData.status != 200 || error != null) {
    return (
      <>
        <p>Error {error}</p>
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
          {loading ? <Loading /> : null}
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
