import { client } from './api'

export async function fetchAllMailboxes() {
  return client.get('/mb/list')
}

export async function prepareMailBox(mbName: string, page: number) {
  const resp = client.put('/mb/select', { mailbox: mbName })
  if ((await resp).status !== 200) {
    return resp
  }
  return fetchEnvelopes(page)
}

export async function fetchEnvelopes(page: number) {
  return client.get(`/envelopes/${page}`)
}
