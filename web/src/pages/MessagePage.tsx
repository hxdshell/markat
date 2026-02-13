import { ArrowLeft, Download, FileIcon, RefreshCw } from 'lucide-react'
import { msgRoute } from '../router'
import IconBtn from '../components/ui/IconBtn'
import { useNavigate } from '@tanstack/react-router'
import { fetchAttachment } from '../api/message'
import { useState } from 'react'

export default function MessagePage() {
  const loaderData: { meta: MessageMeta; data: ApiResponseType } =
    msgRoute.useLoaderData()
  const params = msgRoute.useParams()
  const meta = loaderData.meta
  const message = loaderData.data.data

  const navigate = useNavigate()

  const [attchmntLoading, setAttchmntLoading] = useState(false)

  async function downloadAttachment(
    mb: string,
    uid: number,
    specifier: string,
  ) {
    const attchmnt = await fetchAttachment(mb, uid, specifier)
    if (!attchmnt) return

    setAttchmntLoading(true)

    const url = window.URL.createObjectURL(attchmnt.blob)
    const tempLink = document.createElement('a')
    tempLink.href = url
    tempLink.download = attchmnt.filename

    tempLink.click()
    tempLink.remove()

    window.URL.revokeObjectURL(url)

    setAttchmntLoading(false)
  }
  return (
    <div className="envelope-container">
      <div className="toolbar">
        <IconBtn
          abbr="Go back"
          onClick={() => {
            navigate({ to: `/mb/${meta.mb}` })
          }}
        >
          <ArrowLeft />
        </IconBtn>
      </div>
      <div className="message-container">
        <div className="message-header">
          <p className="subject">{meta.subject}</p>
          <div className="top">
            <p className="from">{meta.from}</p>
            <p className="date">{meta.date}</p>
          </div>
          <p className="to">to {meta.to}</p>
        </div>
        <div className="message-body">
          <pre>{message}</pre>
        </div>
        {meta.attachments ? (
          <div className="attachments">
            <div className="title">
              <p>Attachments</p>
              <i>{meta.attachments.length} Files</i>
              {attchmntLoading ? (
                <span>
                  <RefreshCw />
                </span>
              ) : null}
            </div>
            <ul>
              {meta.attachments.map((item) => (
                <li key={item.specifier}>
                  <span>
                    <FileIcon />
                  </span>
                  <p className="filename">
                    <abbr title={`${item.contentType} ${item.specifier}`}>
                      {item.fileName}
                    </abbr>
                  </p>
                  <p className="size">{item.size}</p>
                  <span>
                    <IconBtn
                      abbr="Download"
                      onClick={() =>
                        downloadAttachment(
                          params.mbName,
                          meta.uid,
                          item.specifier,
                        )
                      }
                    >
                      <Download />
                    </IconBtn>
                  </span>
                </li>
              ))}
            </ul>
          </div>
        ) : null}
      </div>
    </div>
  )
}
