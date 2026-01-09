import { client } from './api'

export async function fetchAllMailboxes() {
  return client.get('/mb/list')
}

export async function prepareMailBox(mbName: string) {
  const resp = client.put('/mb/select', { mailbox: mbName })
  if ((await resp).status !== 200) {
    return resp
  }
  return client.get('/mb/envelopes')
}
