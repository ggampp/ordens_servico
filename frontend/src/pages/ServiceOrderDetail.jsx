import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import client, { apiError } from '../api/client'
import ServiceOrderForm from '../components/ServiceOrderForm'
import { ErrorBox, PriorityBadge, Spinner, StatusBadge } from '../components/ui'
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
  const [editing, setEditing] = useState(false)

  const load = useCallback(() => {
    Promise.all([
      client.get(`/service-orders/${id}`),
      client.get(`/service-orders/${id}/history`),
    ])
      .then(([orderRes, historyRes]) => {
        setOrder(orderRes.data)
        setHistory(historyRes.data)
        setNewStatus(orderRes.data.status)
      })
      .catch((err) => setError(apiError(err)))
  }, [id])

  useEffect(() => { load() }, [load])

  useEffect(() => {
    if (!canManage) return
    client.get('/employees', { params: { page_size: 100, status: 'active' } })
      .then((res) => setEmployees(res.data.items))
      .catch(() => {})
  }, [canManage])

  async function changeStatus() {
    try {
      await client.patch(`/service-orders/${id}/status`, { status: newStatus, note })
      setNote('')
      load()
    } catch (err) {
      setError(apiError(err))
    }
  }

  async function assign() {
    if (!assignTo) return
    try {
      await client.patch(`/service-orders/${id}/assign`, { employee_id: Number(assignTo) })
      setAssignTo('')
      load()
    } catch (err) {
      setError(apiError(err))
    }
  }

  async function remove() {
    if (!confirm('Excluir esta ordem de servico?')) return
    try {
      await client.delete(`/service-orders/${id}`)
      navigate('/service-orders')
    } catch (err) {
      setError(apiError(err))
    }
  }

  if (error && !order) return <ErrorBox message={error} />
  if (!order) return <Spinner />

  const canEdit = order.status === 'open'

  return (
    <div className="max-w-3xl">
      <Link to="/service-orders" className="text-sm text-brand-600 hover:underline">Voltar</Link>
      <div className="flex items-start justify-between gap-3 mt-2 mb-4">
        <div>
          <h1 className="text-xl font-semibold">{order.title}</h1>
          <p className="font-mono text-sm text-slate-500">{order.number}</p>
        </div>
        <div className="flex flex-wrap justify-end gap-2 items-center">
          {canEdit && <button className="btn-secondary" onClick={() => setEditing(true)}>Editar OS</button>}
          <StatusBadge status={order.status} />
          <PriorityBadge priority={order.priority} />
        </div>
      </div>

      <ErrorBox message={error} />

      <div className="card mb-4 grid grid-cols-1 sm:grid-cols-2 gap-y-2 text-sm">
        <Field label="Responsavel" value={order.employee_name || '-'} />
        <Field label="Endereco" value={order.address || '-'} />
        <Field label="Abertura" value={fmtDate(order.opened_at)} />
        <Field label="Prevista" value={fmtDate(order.due_at)} />
        <Field label="Conclusao" value={fmtDate(order.completed_at)} />
        <Field label="Coordenadas" value={order.latitude ? `${order.latitude}, ${order.longitude}` : '-'} />
        <div className="sm:col-span-2"><Field label="Descricao" value={order.description || '-'} /></div>
        <div className="sm:col-span-2"><Field label="Observacoes" value={order.notes || '-'} /></div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div className="card">
          <h2 className="font-medium mb-3">Alterar Status</h2>
          <select className="input mb-2" value={newStatus} onChange={(e) => setNewStatus(e.target.value)}>
            {STATUS_OPTIONS.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
          </select>
          <input className="input mb-2" placeholder="Observacao opcional" value={note} onChange={(e) => setNote(e.target.value)} />
          <button className="btn-primary w-full" onClick={changeStatus}>Atualizar Status</button>
        </div>

        {canManage && (
          <div className="card">
            <h2 className="font-medium mb-3">Atribuir Responsavel</h2>
            <select className="input mb-2" value={assignTo} onChange={(e) => setAssignTo(e.target.value)}>
              <option value="">Selecione</option>
              {employees.map((employee) => <option key={employee.id} value={employee.id}>{employee.name}</option>)}
            </select>
            <button className="btn-primary w-full mb-2" onClick={assign}>Atribuir</button>
            <button className="btn-danger w-full" onClick={remove}>Excluir OS</button>
          </div>
        )}
      </div>

      <div className="card">
        <h2 className="font-medium mb-3">Historico de Status</h2>
        <ol className="relative border-l border-slate-200 ml-2">
          {history.map((item) => (
            <li key={item.id} className="mb-4 ml-4">
              <div
                className="absolute w-2.5 h-2.5 rounded-full -left-[5px] mt-1.5"
                style={{ background: STATUS[item.new_status]?.color || '#94a3b8' }}
              />
              <p className="text-sm">
                {item.old_status ? `${STATUS[item.old_status]?.label || item.old_status} -> ` : ''}
                <strong>{STATUS[item.new_status]?.label || item.new_status}</strong>
              </p>
              <p className="text-xs text-slate-400">
                {fmtDate(item.changed_at)}{item.changed_by_name ? ` - ${item.changed_by_name}` : ''}{item.note ? ` - ${item.note}` : ''}
              </p>
            </li>
          ))}
          {history.length === 0 && <p className="text-sm text-slate-400 ml-4">Sem historico.</p>}
        </ol>
      </div>

      {editing && (
        <ServiceOrderForm
          initial={order}
          employees={employees}
          onClose={() => setEditing(false)}
          onSaved={() => {
            setEditing(false)
            load()
          }}
        />
      )}
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
