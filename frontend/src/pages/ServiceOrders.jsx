import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import client, { apiError } from '../api/client'
import ServiceOrderForm from '../components/ServiceOrderForm'
import { ErrorBox, Pagination, PriorityBadge, Spinner, StatusBadge } from '../components/ui'
import { PRIORITY_OPTIONS, STATUS_OPTIONS, fmtDate } from '../constants'
import { useAuth } from '../auth/AuthContext'

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
    Object.entries(filters).forEach(([key, value]) => {
      if (value) params[key] = value
    })
    client.get('/service-orders', { params })
      .then((res) => setData(res.data))
      .catch((err) => setError(apiError(err)))
  }, [page, filters])

  useEffect(() => { load() }, [load])

  useEffect(() => {
    if (!canManage) return
    client.get('/employees', { params: { page_size: 100, status: 'active' } })
      .then((res) => setEmployees(res.data.items))
      .catch(() => {})
  }, [canManage])

  function setFilter(key, value) {
    setPage(1)
    setFilters((current) => ({ ...current, [key]: value }))
  }

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
        <h1 className="text-xl font-semibold">Ordens de Servico</h1>
        {canManage && <button className="btn-primary" onClick={() => setCreating(true)}>+ Nova OS</button>}
      </div>

      <div className="card mb-4 grid grid-cols-2 md:grid-cols-5 gap-3">
        <div>
          <label className="label">Status</label>
          <select className="input" value={filters.status} onChange={(e) => setFilter('status', e.target.value)}>
            <option value="">Todos</option>
            {STATUS_OPTIONS.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
          </select>
        </div>
        <div>
          <label className="label">Prioridade</label>
          <select className="input" value={filters.priority} onChange={(e) => setFilter('priority', e.target.value)}>
            <option value="">Todas</option>
            {PRIORITY_OPTIONS.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
          </select>
        </div>
        {canManage && (
          <div>
            <label className="label">Responsavel</label>
            <select className="input" value={filters.employee_id} onChange={(e) => setFilter('employee_id', e.target.value)}>
              <option value="">Todos</option>
              {employees.map((employee) => <option key={employee.id} value={employee.id}>{employee.name}</option>)}
            </select>
          </div>
        )}
        <div>
          <label className="label">De</label>
          <input type="date" className="input" value={filters.date_from} onChange={(e) => setFilter('date_from', e.target.value)} />
        </div>
        <div>
          <label className="label">Ate</label>
          <input type="date" className="input" value={filters.date_to} onChange={(e) => setFilter('date_to', e.target.value)} />
        </div>
      </div>

      <ErrorBox message={error} />

      {!data ? <Spinner /> : (
        <div className="card overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-slate-500 border-b">
                <th className="py-2 pr-4">Numero</th>
                <th className="py-2 pr-4">Titulo</th>
                <th className="py-2 pr-4">Responsavel</th>
                <th className="py-2 pr-4">Prioridade</th>
                <th className="py-2 pr-4">Status</th>
                <th className="py-2 pr-4">Abertura</th>
              </tr>
            </thead>
            <tbody>
              {data.items.map((order) => (
                <tr
                  key={order.id}
                  className="border-b last:border-0 hover:bg-slate-50 cursor-pointer"
                  onClick={() => navigate(`/service-orders/${order.id}`)}
                >
                  <td className="py-2 pr-4 font-mono">{order.number}</td>
                  <td className="py-2 pr-4">{order.title}</td>
                  <td className="py-2 pr-4">{order.employee_name || '-'}</td>
                  <td className="py-2 pr-4"><PriorityBadge priority={order.priority} /></td>
                  <td className="py-2 pr-4"><StatusBadge status={order.status} /></td>
                  <td className="py-2 pr-4 text-xs text-slate-500">{fmtDate(order.opened_at)}</td>
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
        <ServiceOrderForm
          employees={employees}
          onClose={() => setCreating(false)}
          onSaved={() => {
            setCreating(false)
            load()
          }}
        />
      )}
    </div>
  )
}
