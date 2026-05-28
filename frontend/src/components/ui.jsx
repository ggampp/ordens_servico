import { STATUS, PRIORITY } from '../constants'

export function StatusBadge({ status }) {
  const s = STATUS[status] || { label: status, badge: 'bg-slate-400' }
  return <span className={`badge ${s.badge}`}>{s.label}</span>
}

export function PriorityBadge({ priority }) {
  const p = PRIORITY[priority] || { label: priority, badge: 'bg-slate-400' }
  return <span className={`badge ${p.badge}`}>{p.label}</span>
}

export function Modal({ open, title, onClose, children }) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-[1000] flex items-center justify-center bg-black/40 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between border-b px-5 py-3">
          <h3 className="font-semibold">{title}</h3>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-600 text-xl leading-none">×</button>
        </div>
        <div className="p-5">{children}</div>
      </div>
    </div>
  )
}

export function Pagination({ page, totalPages, onChange }) {
  if (totalPages <= 1) return null
  return (
    <div className="flex items-center justify-center gap-2 mt-4 text-sm">
      <button className="btn-secondary" disabled={page <= 1} onClick={() => onChange(page - 1)}>
        Anterior
      </button>
      <span className="px-2">
        Página {page} de {totalPages}
      </span>
      <button className="btn-secondary" disabled={page >= totalPages} onClick={() => onChange(page + 1)}>
        Próxima
      </button>
    </div>
  )
}

export function Spinner() {
  return <div className="text-center text-slate-400 py-8">Carregando…</div>
}

export function ErrorBox({ message }) {
  if (!message) return null
  return <div className="rounded bg-red-50 border border-red-200 text-red-700 text-sm px-3 py-2 mb-3">{message}</div>
}
