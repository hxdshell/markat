import type { ReactNode } from 'react'

export default function IconBtn({
  children,
  onClick,
  abbr,
}: {
  children: ReactNode
  onClick: () => void
  abbr: string
}) {
  return (
    <button className="icon-btn" onClick={onClick}>
      <abbr title={abbr}>{children}</abbr>
    </button>
  )
}
