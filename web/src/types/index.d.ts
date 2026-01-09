interface ApiResponseType {
  status: number
  message: string
  data: any
  cause?: string
}

interface Envelope {
  uid: number
  internalDate: string
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
