interface ApiResponseType {
  status: number
  message: string
  data: any
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
