import { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import client, { apiError } from '../api/client'
import { Spinner, ErrorBox, StatusBadge, PriorityBadge } from '../components/ui'
import { STATUS, STATUS_OPTIONS, fmtDate } from '../constants'
import { useAuth } from '../auth/AuthContext'

export default function ServiceOrderDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { canManage } = useAuth()
  const [order, setOrder] = useState(null)
  const [history, setHistory] = useState([])
  const [employees, setEmployees] = useState([])
  const [error, setError] = useState('')
  const [newStatus, setNewStatus] = useState('')
  const [note, setNote] = useState('')
  const [assignTo, setAssignTo] = useState('')

  const load = useCallback(() => {
    Promise.all([
      client.get(`/service-orders/${id}`),
      client.get(`/service-orders/${id}/history`),
    ])
      .then(([o, h]) => { setOrder(o.data); setHistory(h.data); setNewStatus(o.data.status) })
      .catch((err) => setError(apiError(err)))
  }, [id])

  useEffect(() => { load() }, [load])
  useEffect(() => {
    if (canManage) {
      client.get('/employees', { params: { page_size: 100, status: 'active' } })
        .then((res) => setEmployees(res.data.items)).catch(() => {})
    }
  }, [canManage])

  async function changeStatus() {
    try {
      await client.patch(`/service-orders/${id}/status`, { status: newStatus, note })
      setNote(''); load()
    } catch (err) { setError(apiError(err)) }
  }

  async function assign() {
    if (!assignTo) return
    try {
      await client.patch(`/service-orders/${id}/assign`, { employee_id: Number(assignTo) })
      setAssignTo(''); load()
    } catch (err) { setError(apiError(err)) }
  }

  async function remove() {
    if (!confirm('Excluir esta ordem de serviço?')) return
    try { await client.delete(`/service-orders/${id}`); navigate('/service-orders') }
    catch (err) { setError(apiError(err)) }
  }

  if (error && !order) return <ErrorBox message={error} />
  if (!order) return <Spinner />

  return (
    <div className="max-w-3xl">
      <Link to="/service-orders" className="text-sm text-brand-600 hover:underline">← Voltar</Link>
      <div className="flex items-start justify-between mt-2 mb-4">
        <div>
          <h1 className="text-xl font-semibold">{order.title}</h1>
          <p className="font-mono text-sm text-slate-500">{order.number}</p>
        </div>
        <div className="flex gap-2 items-center">
          <StatusBadge status={order.status} />
          <PriorityBadge priority={order.priority} />
        </div>
      </div>

      <ErrorBox message={error} />

      <div className="card mb-4 grid grid-cols-1 sm:grid-cols-2 gap-y-2 text-sm">
        <Field label="Responsável" value={order.employee_name || '—'} />
        <Field label="Endereço" value={order.address || '—'} />
        <Field label="Abertura" value={fmtDate(order.opened_at)} />
        <Field label="Prevista" value={fmtDate(order.due_at)} />
        <Field label="Conclusão" value={fmtDate(order.completed_at)} />
        <Field label="Coordenadas" value={order.latitude ? `${order.latitude}, ${order.longitude}` : '—'} />
        <div className="sm:col-span-2"><Field label="Descrição" value={order.description || '—'} /></div>
        <div className="sm:col-span-2"><Field label="Observações" value={order.notes || '—'} /></div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div className="card">
          <h2 className="font-medium mb-3">Alterar Status</h2>
          <select className="input mb-2" value={newStatus} onChange={(e) => setNewStatus(e.target.value)}>
            {STATUS_OPTIONS.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
          </select>
          <input className="input mb-2" placeholder="Observação (opcional)" value={note} onChange={(e) => setNote(e.target.value)} />
          <button className="btn-primary w-full" onClick={changeStatus}>Atualizar Status</button>
        </div>

        {canManage && (
          <div className="card">
            <h2 className="font-medium mb-3">Atribuir Responsável</h2>
            <select className="input mb-2" value={assignTo} onChange={(e) => setAssignTo(e.target.value)}>
              <option value="">Selecione…</option>
              {employees.map((e) => <option key={e.id} value={e.id}>{e.name}</option>)}
            </select>
            <button className="btn-primary w-full mb-2" onClick={assign}>Atribuir</button>
            <button className="btn-danger w-full" onClick={remove}>Excluir OS</button>
          </div>
        )}
      </div>

      <div className="card">
        <h2 className="font-medium mb-3">Histórico de Status</h2>
        <ol className="relative border-l border-slate-200 ml-2">
          {history.map((h) => (
            <li key={h.id} className="mb-4 ml-4">
              <div className="absolute w-2.5 h-2.5 rounded-full -left-[5px] mt-1.5"
                   style={{ background: STATUS[h.new_status]?.color || '#94a3b8' }} />
              <p className="text-sm">
                {h.old_status ? `${STATUS[h.old_status]?.label || h.old_status} → ` : ''}
                <strong>{STATUS[h.new_status]?.label || h.new_status}</strong>
              </p>
              <p className="text-xs text-slate-400">
                {fmtDate(h.changed_at)}{h.changed_by_name ? ` · ${h.changed_by_name}` : ''}{h.note ? ` · ${h.note}` : ''}
              </p>
            </li>
          ))}
          {history.length === 0 && <p className="text-sm text-slate-400 ml-4">Sem histórico.</p>}
        </ol>
      </div>
    </div>
  )
}

function Field({ label, value }) {
  return (
    <div>
      <span className="text-slate-400 text-xs uppercase tracking-wide">{label}</span>
      <p>{value}</p>
    </div>
  )
}
