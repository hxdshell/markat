import {
  createRootRoute,
  createRoute,
  createRouter,
  notFound,
  Outlet,
  redirect,
} from '@tanstack/react-router'
import Layout from '../pages/Layout'
import { fetchAllMailboxes } from '../api/mailbox'
import MailBoxPage from '../pages/MailBoxPage'
import ErrorPage from '../pages/ErrorPage'
import { fetchMessage, fetchMeta } from '../api/message'
import MessagePage from '../pages/MessagePage'
import Loading from '../components/ui/Loading'

const rootRoute = createRootRoute({
  component: Outlet,
  errorComponent: ({ error }) => <ErrorPage error={error.message} />,
})

export const appRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'app',
  loader: () => fetchAllMailboxes(),
  shouldReload: false,
  component: Layout,
})
const indexRoute = createRoute({
  getParentRoute: () => appRoute,
  path: '/',
  beforeLoad: () => {
    // this redirect is based on assumption which is bad but is here for speedy development and should be removed later
    // prepare a mailbox only after some kind of CAPABILITY check
    throw redirect({ to: '/mb/INBOX' })
  },
})

const mbRoute = createRoute({
  getParentRoute: () => appRoute,
  path: 'mb',
  component: Outlet,
})
export const mbNameRoute = createRoute({
  getParentRoute: () => mbRoute,
  path: '$mbName',
  validateSearch: (search: Record<string, unknown>) => {
    const page = parseInt(String(search.page), 10)
    return { page: page > 0 ? page : 1 }
  },
  component: MailBoxPage,
})

export const msgRoute = createRoute({
  getParentRoute: () => appRoute,
  path: 'message/$mbName/$uid',
  loader: async ({ params }) => {
    const uid = parseInt(String(params.uid), 10)
    if (uid <= 0) {
      throw notFound()
    }
    const meta = await fetchMeta(params.mbName, uid)
    const data = await fetchMessage(params.mbName, uid)
    return { meta: meta.data, data: data }
  },
  component: MessagePage,
  pendingComponent: Loading,
  pendingMs: 0,
})

const routeTree = rootRoute.addChildren([
  appRoute.addChildren([
    indexRoute,
    mbRoute.addChildren([mbNameRoute]),
    msgRoute,
  ]),
])
export const router = createRouter({ routeTree })
