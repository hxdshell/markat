interface ApiResponseType {
  status: number
  message: string
  data: any
  contentType?: string
  cause?: string
}

interface EnvelopeResponse {
  page: number
  start: number
  end: number
  total: number
  envelopes: Envelope[]
}
interface Envelope {
  uid: number
  from: string[]
  fromName: string[]
  to: string[]
  toName: string[]
  sender: string[]
  date: string[]
  subject: string
  size: string
  flags: string[]
}
interface MessageAttachment {
  specifier: string
  contentType: string
  encoding: string
  fileName: string
  size: string
}

interface MessageMeta {
  uid: number
  mb: string
  from: string
  to: string
  subject: string
  date: string
  attachments: MessageAttachment[] | null
}
