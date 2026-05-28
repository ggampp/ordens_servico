// Status metadata: labels and marker/badge colors per the specification.
export const STATUS = {
  open: { label: 'Aberta', color: '#dc2626', badge: 'bg-red-600' },
  assigned: { label: 'Atribuída', color: '#ea580c', badge: 'bg-orange-500' },
  in_progress: { label: 'Em Atendimento', color: '#2563eb', badge: 'bg-blue-600' },
  completed: { label: 'Concluída', color: '#16a34a', badge: 'bg-green-600' },
  cancelled: { label: 'Cancelada', color: '#6b7280', badge: 'bg-gray-500' },
}

export const PRIORITY = {
  low: { label: 'Baixa', badge: 'bg-slate-400' },
  medium: { label: 'Média', badge: 'bg-sky-500' },
  high: { label: 'Alta', badge: 'bg-amber-500' },
  urgent: { label: 'Urgente', badge: 'bg-red-600' },
}

export const ROLE_LABELS = {
  admin: 'Administrador',
  supervisor: 'Supervisor',
  operator: 'Operador',
}

export const STATUS_OPTIONS = Object.entries(STATUS).map(([value, v]) => ({ value, label: v.label }))
export const PRIORITY_OPTIONS = Object.entries(PRIORITY).map(([value, v]) => ({ value, label: v.label }))

export function fmtDate(s) {
  if (!s) return '—'
  return new Date(s).toLocaleString('pt-BR')
}
