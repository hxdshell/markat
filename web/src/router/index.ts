import {
  createRootRoute,
  createRoute,
  createRouter,
  Outlet,
  redirect,
} from '@tanstack/react-router'
import Layout from '../pages/Layout'
import { fetchAllMailboxes, prepareMailBox } from '../api/mailbox'
import MailBoxPage from '../pages/MailBoxPage'
import Loading from '../components/ui/Loading'

const rootRoute = createRootRoute({ component: Outlet })

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
  loaderDeps: ({ search: { page } }) => ({ page }),
  loader: ({ params, deps }) => prepareMailBox(params.mbName, deps.page),
  pendingComponent: Loading,
  pendingMinMs: 0,
  component: MailBoxPage,
})
const routeTree = rootRoute.addChildren([
  appRoute.addChildren([indexRoute, mbRoute.addChildren([mbNameRoute])]),
])
export const router = createRouter({ routeTree })
