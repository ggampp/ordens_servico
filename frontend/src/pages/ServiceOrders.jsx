import { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import client, { apiError } from '../api/client'
import { Modal, Pagination, Spinner, ErrorBox, StatusBadge, PriorityBadge } from '../components/ui'
import { STATUS_OPTIONS, PRIORITY_OPTIONS, fmtDate } from '../constants'
import { useAuth } from '../auth/AuthContext'

const empty = {
  title: '', description: '', priority: 'medium', employee_id: '',
  address: '', latitude: '', longitude: '', due_at: '', notes: '',
}

export default function ServiceOrders() {
  const { canManage } = useAuth()
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [data, setData] = useState(null)
  const [employees, setEmployees] = useState([])
  const [filters, setFilters] = useState({ status: '', priority: '', employee_id: '', date_from: '', date_to: '' })
  const [error, setError] = useState('')
  const [creating, setCreating] = useState(false)

  const load = useCallback(() => {
    const params = { page, page_size: 10 }
    Object.entries(filters).forEach(([k, v]) => { if (v) params[k] = v })
    client.get('/service-orders', { params })
      .then((res) => setData(res.data))
      .catch((err) => setError(apiError(err)))
  }, [page, filters])

  useEffect(() => { load() }, [load])
  useEffect(() => {
    if (canManage) {
      client.get('/employees', { params: { page_size: 100, status: 'active' } })
        .then((res) => setEmployees(res.data.items)).catch(() => {})
    }
  }, [canManage])

  function setFilter(k, v) { setPage(1); setFilters((f) => ({ ...f, [k]: v })) }

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
        <h1 className="text-xl font-semibold">Ordens de Serviço</h1>
        {canManage && <button className="btn-primary" onClick={() => setCreating(true)}>+ Nova OS</button>}
      </div>

      <div className="card mb-4 grid grid-cols-2 md:grid-cols-5 gap-3">
        <div>
          <label className="label">Status</label>
          <select className="input" value={filters.status} onChange={(e) => setFilter('status', e.target.value)}>
            <option value="">Todos</option>
            {STATUS_OPTIONS.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
          </select>
        </div>
        <div>
          <label className="label">Prioridade</label>
          <select className="input" value={filters.priority} onChange={(e) => setFilter('priority', e.target.value)}>
            <option value="">Todas</option>
            {PRIORITY_OPTIONS.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
          </select>
        </div>
        {canManage && (
          <div>
            <label className="label">Responsável</label>
            <select className="input" value={filters.employee_id} onChange={(e) => setFilter('employee_id', e.target.value)}>
              <option value="">Todos</option>
              {employees.map((e) => <option key={e.id} value={e.id}>{e.name}</option>)}
            </select>
          </div>
        )}
        <div>
          <label className="label">De</label>
          <input type="date" className="input" value={filters.date_from} onChange={(e) => setFilter('date_from', e.target.value)} />
        </div>
        <div>
          <label className="label">Até</label>
          <input type="date" className="input" value={filters.date_to} onChange={(e) => setFilter('date_to', e.target.value)} />
        </div>
      </div>

      <ErrorBox message={error} />

      {!data ? <Spinner /> : (
        <div className="card overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-slate-500 border-b">
                <th className="py-2 pr-4">Número</th>
                <th className="py-2 pr-4">Título</th>
                <th className="py-2 pr-4">Responsável</th>
                <th className="py-2 pr-4">Prioridade</th>
                <th className="py-2 pr-4">Status</th>
                <th className="py-2 pr-4">Abertura</th>
              </tr>
            </thead>
            <tbody>
              {data.items.map((o) => (
                <tr key={o.id} className="border-b last:border-0 hover:bg-slate-50 cursor-pointer"
                    onClick={() => navigate(`/service-orders/${o.id}`)}>
                  <td className="py-2 pr-4 font-mono">{o.number}</td>
                  <td className="py-2 pr-4">{o.title}</td>
                  <td className="py-2 pr-4">{o.employee_name || '—'}</td>
                  <td className="py-2 pr-4"><PriorityBadge priority={o.priority} /></td>
                  <td className="py-2 pr-4"><StatusBadge status={o.status} /></td>
                  <td className="py-2 pr-4 text-xs text-slate-500">{fmtDate(o.opened_at)}</td>
                </tr>
              ))}
              {data.items.length === 0 && (
                <tr><td colSpan={6} className="text-center text-slate-400 py-6">Nenhuma ordem.</td></tr>
              )}
            </tbody>
          </table>
          <Pagination page={data.page} totalPages={data.total_pages} onChange={setPage} />
        </div>
      )}

      {creating && (
        <OrderForm employees={employees} onClose={() => setCreating(false)} onSaved={() => { setCreating(false); load() }} />
      )}
    </div>
  )
}

function OrderForm({ employees, onClose, onSaved }) {
  const [form, setForm] = useState(empty)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  function up(k, v) { setForm((f) => ({ ...f, [k]: v })) }

  async function submit(e) {
    e.preventDefault()
    setSaving(true); setError('')
    const payload = {
      title: form.title,
      description: form.description || null,
      priority: form.priority,
      employee_id: form.employee_id ? Number(form.employee_id) : null,
      address: form.address || null,
      latitude: form.latitude ? parseFloat(form.latitude) : null,
      longitude: form.longitude ? parseFloat(form.longitude) : null,
      due_at: form.due_at ? new Date(form.due_at).toISOString() : null,
      notes: form.notes || null,
    }
    try {
      await client.post('/service-orders', payload)
      onSaved()
    } catch (err) { setError(apiError(err)); setSaving(false) }
  }

  return (
    <Modal open title="Nova Ordem de Serviço" onClose={onClose}>
      <form onSubmit={submit}>
        <ErrorBox message={error} />
        <div className="mb-3">
          <label className="label">Título *</label>
          <input className="input" value={form.title} onChange={(e) => up('title', e.target.value)} required />
        </div>
        <div className="mb-3">
          <label className="label">Descrição</label>
          <textarea className="input" rows={2} value={form.description} onChange={(e) => up('description', e.target.value)} />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">Prioridade</label>
            <select className="input" value={form.priority} onChange={(e) => up('priority', e.target.value)}>
              {PRIORITY_OPTIONS.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
            </select>
          </div>
          <div className="mb-3">
            <label className="label">Responsável</label>
            <select className="input" value={form.employee_id} onChange={(e) => up('employee_id', e.target.value)}>
              <option value="">— Não atribuir —</option>
              {employees.map((e) => <option key={e.id} value={e.id}>{e.name}</option>)}
            </select>
          </div>
        </div>
        <div className="mb-3">
          <label className="label">Endereço</label>
          <input className="input" value={form.address} onChange={(e) => up('address', e.target.value)} />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">Latitude</label>
            <input className="input" value={form.latitude} onChange={(e) => up('latitude', e.target.value)} />
          </div>
          <div className="mb-3">
            <label className="label">Longitude</label>
            <input className="input" value={form.longitude} onChange={(e) => up('longitude', e.target.value)} />
          </div>
        </div>
        <div className="mb-4">
          <label className="label">Data prevista</label>
          <input type="datetime-local" className="input" value={form.due_at} onChange={(e) => up('due_at', e.target.value)} />
        </div>
        <div className="flex justify-end gap-2">
          <button type="button" className="btn-secondary" onClick={onClose}>Cancelar</button>
          <button className="btn-primary" disabled={saving}>{saving ? 'Salvando…' : 'Criar'}</button>
        </div>
      </form>
    </Modal>
  )
}
