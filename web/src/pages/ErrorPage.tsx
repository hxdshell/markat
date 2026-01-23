export default function ErrorPage({ error }: { error: string }) {
  return (
    <div className="error-container">
      <h1>Error</h1>
      <p>{error}</p>
    </div>
  )
}
