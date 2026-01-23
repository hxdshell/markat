import { mbNameRoute } from '../router'
import { useEffect, useState } from 'react'
import { fetchEnvelopes, prepareMailBox } from '../api/mailbox'
import EnvelopeList from '../components/app/EnvelopeList'
import Loading from '../components/ui/Loading'

export default function MailBoxPage() {
  const { mbName } = mbNameRoute.useParams()
  const { page } = mbNameRoute.useSearch()
  const navigate = mbNameRoute.useNavigate()

  const [data, setData] = useState<EnvelopeResponse>()
  const [loading, setLoading] = useState(false)
  const [selectedMb, setSelectedMb] = useState('')
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setSelectedMb('')
    async function fetchData() {
      try {
        const resp = await prepareMailBox(mbName)
        if (resp.status !== 200) {
          setError(resp.message)
        } else {
          setSelectedMb(mbName)
        }
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message)
        } else {
          setError('unkown err')
        }
      }
    }
    fetchData()
  }, [mbName])

  // pagination
  useEffect(() => {
    if (selectedMb !== mbName) return
    async function fetchData() {
      try {
        setLoading(true)
        const resp = await fetchEnvelopes(page)
        if (resp.status !== 200) {
          if (resp.status === 404) {
            navigate({ search: (s: any) => ({ ...s, page: 1 }) })
          } else {
            setError(resp.message)
          }
        } else {
          setData(resp.data)
        }
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message)
        } else {
          setError('unkown err')
        }
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [page, selectedMb])

  async function refresh() {
    if (page === 1) {
      try {
        setLoading(true)
        const resp = await fetchEnvelopes(page)
        setData(resp.data)
        if (resp.status !== 200) {
          setError(resp.message)
        }
      } catch (err) {
        if (err instanceof Error) {
          setError(err.message)
        } else {
          setError('unkown error')
        }
      } finally {
        setLoading(false)
      }
    } else {
      navigate({ search: (s: any) => ({ ...s, page: 1 }) })
    }
  }
  async function prev(data: EnvelopeResponse) {
    if (data.page > 1) {
      navigate({ search: (s: any) => ({ ...s, page: data.page - 1 }) })
    }
  }
  async function next(data: EnvelopeResponse) {
    if (data.end !== data.total) {
      navigate({ search: (s: any) => ({ ...s, page: data.page + 1 }) })
    }
  }

  if (error) {
    return (
      <>
        <p>{error}</p>
      </>
    )
  } else {
    if (data) {
      return (
        <EnvelopeList
          mbName={mbName}
          data={data}
          loading={loading}
          refresh={refresh}
          prev={prev}
          next={next}
        />
      )
    } else {
      return <Loading />
    }
  }
}
