import { mbNameRoute } from '../router'

export default function MailBoxPage() {
  const loaderData: ApiResponseType = mbNameRoute.useLoaderData()

  if (loaderData.status != 200) {
    return <></>
  } else {
    const envelopes: Envelope[] = loaderData.data
    if (envelopes.length === 0) {
      return <div>There are no mail in this mailbox</div>
    }
    return (
      <div className="envelope-container">
        <table className="envelope-list">
          <tbody>
            {loaderData.data.map((envlp: Envelope) => (
              <tr
                className={
                  envlp.flags.includes('\\Seen') ? 'envelope' : 'envelope bold'
                }
              >
                <td className="from">{envlp.from}</td>
                <td className="subject">{envlp.subject}</td>
                <td className="meta">
                  <span className="date">{envlp.date}</span>
                  <span className="size">{envlp.size}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )
  }
}
