import { Link, Outlet } from '@tanstack/react-router'
import MailBoxList from '../components/layout/MailBoxList'
import { appRoute } from '../router'

export default function Layout() {
  const loaderData: ApiResponseType = appRoute.useLoaderData()
  return (
    <div className="container">
      <header>
        <div className="logo-box">
          <Link to={'/'}>
            <h1>Markat</h1>
          </Link>
        </div>
        <div className="search-bar-box"></div>
      </header>
      <div className="layout-container">
        <div className="left-nav">
          <MailBoxList list={loaderData.data} />
        </div>
        <div className="main-section">
          <Outlet />
        </div>
      </div>
      <footer></footer>
    </div>
  )
}
