import { API_URL, client } from './api'

export async function fetchMeta(mb: string, uid: number) {
  return client.get(`/meta/${mb}/${uid}`)
}
export async function fetchMessage(mb: string, uid: number) {
  return client.get(`/message/${mb}/${uid}`)
}

type AttachmentResponse = {
  blob: Blob
  filename: string
}
export async function fetchAttachment(
  mb: string,
  uid: number,
  specifier: string,
): Promise<AttachmentResponse | null> {
  const res = await fetch(`${API_URL}/attachment/${mb}/${uid}/${specifier}`)

  if (!res.ok) {
    console.log(res.status)
    return null // ignore error for now since server response will be empty
  }
  const disposition = res.headers.get('Content-Disposition')
  const blob = await res.blob()

  const match = disposition?.match(/filename="([^"]+)"/)

  let filename = `${mb}-${uid}-${specifier}`
  if (match) {
    if (match.length > 1) {
      filename = match[1]
    }
  }

  return { blob: blob, filename: filename }
}

export async function markSeenUnseen(
  mb: string,
  uid: number | string,
  seen: boolean,
) {
  return client.put(`/mark-seen-unseen/${mb}/${uid}`, { seen: seen })
}

export async function moveToFolder(
  mb: string,
  uid: number | string,
  dest: string,
) {
  return client.put(`/move/${mb}/${uid}`, { dest: dest })
}
