import { RefreshCw } from 'lucide-react'

export default function Loading() {
  return (
    <div className="loading-container">
      <div>
        <RefreshCw />
        <p>Loading...</p>
      </div>
    </div>
  )
}
