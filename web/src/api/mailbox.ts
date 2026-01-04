import { client } from './api'

export async function fetchAllMailboxes() {
  return client.get('/mb/list')
}

export async function prepareMailBox(mbName: string) {
  return client.put('/mb/select', { mailbox: mbName })
}
