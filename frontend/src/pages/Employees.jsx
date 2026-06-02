import { useCallback, useEffect, useState } from 'react'
import client, { apiError } from '../api/client'
import { ErrorBox, Modal, Pagination, Spinner } from '../components/ui'
import { fmtDate } from '../constants'

const empty = { name: '', email: '', phone: '', role: '', status: 'active' }

export default function Employees() {
  const [page, setPage] = useState(1)
  const [data, setData] = useState(null)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [error, setError] = useState('')
  const [editing, setEditing] = useState(null)
  const [posFor, setPosFor] = useState(null)

  const load = useCallback(() => {
    const params = { page, page_size: 10 }
    if (search) params.search = search
    if (statusFilter) params.status = statusFilter
    client
      .get('/employees', { params })
      .then((res) => setData(res.data))
      .catch((err) => setError(apiError(err)))
  }, [page, search, statusFilter])

  useEffect(() => { load() }, [load])

  async function remove(id) {
    if (!confirm('Confirmar exclusao logica deste empregado?')) return
    try {
      await client.delete(`/employees/${id}`)
      load()
    } catch (err) {
      setError(apiError(err))
    }
  }

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
        <h1 className="text-xl font-semibold">Empregados</h1>
        <button className="btn-primary" onClick={() => setEditing(empty)}>+ Novo Empregado</button>
      </div>

      <div className="card mb-4 flex flex-wrap gap-3 items-end">
        <div>
          <label className="label">Buscar</label>
          <input
            className="input"
            value={search}
            onChange={(e) => {
              setPage(1)
              setSearch(e.target.value)
            }}
            placeholder="Nome ou codigo"
          />
        </div>
        <div>
          <label className="label">Status</label>
          <select
            className="input"
            value={statusFilter}
            onChange={(e) => {
              setPage(1)
              setStatusFilter(e.target.value)
            }}
          >
            <option value="">Todos</option>
            <option value="active">Ativo</option>
            <option value="inactive">Inativo</option>
          </select>
        </div>
      </div>

      <ErrorBox message={error} />

      {!data ? <Spinner /> : (
        <div className="card overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-slate-500 border-b">
                <th className="py-2 pr-4">Codigo</th>
                <th className="py-2 pr-4">Nome</th>
                <th className="py-2 pr-4">Cargo</th>
                <th className="py-2 pr-4">Contato</th>
                <th className="py-2 pr-4">Status</th>
                <th className="py-2 pr-4">Ultima posicao</th>
                <th className="py-2 pr-4 text-right">Acoes</th>
              </tr>
            </thead>
            <tbody>
              {data.items.map((employee) => (
                <tr key={employee.id} className="border-b last:border-0">
                  <td className="py-2 pr-4 font-mono">{employee.code}</td>
                  <td className="py-2 pr-4">{employee.name}</td>
                  <td className="py-2 pr-4">{employee.role || '-'}</td>
                  <td className="py-2 pr-4">{employee.email || employee.phone || '-'}</td>
                  <td className="py-2 pr-4">
                    <span className={`badge ${employee.status === 'active' ? 'bg-green-600' : 'bg-gray-500'}`}>
                      {employee.status === 'active' ? 'Ativo' : 'Inativo'}
                    </span>
                  </td>
                  <td className="py-2 pr-4 text-xs text-slate-500">
                    {employee.last_position ? fmtDate(employee.last_position.recorded_at) : '-'}
                  </td>
                  <td className="py-2 pr-4 text-right whitespace-nowrap">
                    <button className="text-brand-600 hover:underline mr-3" onClick={() => setPosFor(employee)}>Posicao</button>
                    <button className="text-brand-600 hover:underline mr-3" onClick={() => setEditing(employee)}>Editar</button>
                    <button className="text-red-600 hover:underline" onClick={() => remove(employee.id)}>Excluir</button>
                  </td>
                </tr>
              ))}
              {data.items.length === 0 && (
                <tr><td colSpan={7} className="text-center text-slate-400 py-6">Nenhum empregado.</td></tr>
              )}
            </tbody>
          </table>
          <Pagination page={data.page} totalPages={data.total_pages} onChange={setPage} />
        </div>
      )}

      {editing && (
        <EmployeeForm
          initial={editing}
          onClose={() => setEditing(null)}
          onSaved={() => {
            setEditing(null)
            load()
          }}
        />
      )}
      {posFor && (
        <PositionForm
          employee={posFor}
          onClose={() => setPosFor(null)}
          onSaved={() => {
            setPosFor(null)
            load()
          }}
        />
      )}
    </div>
  )
}

function EmployeeForm({ initial, onClose, onSaved }) {
  const isNew = !initial.id
  const [form, setForm] = useState(initial)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)

  function up(key, value) {
    setForm((current) => ({ ...current, [key]: value }))
  }

  async function submit(event) {
    event.preventDefault()
    setSaving(true)
    setError('')
    const payload = {
      name: form.name,
      email: form.email || null,
      phone: form.phone || null,
      role: form.role || null,
      status: form.status,
    }
    try {
      if (isNew) {
        await client.post('/employees', payload)
      } else {
        await client.put(`/employees/${initial.id}`, payload)
      }
      onSaved()
    } catch (err) {
      setError(apiError(err))
      setSaving(false)
    }
  }

  return (
    <Modal open title={isNew ? 'Novo Empregado' : 'Editar Empregado'} onClose={onClose}>
      <form onSubmit={submit}>
        <ErrorBox message={error} />
        {!isNew && (
          <div className="mb-3">
            <label className="label">Codigo</label>
            <input className="input font-mono bg-slate-50" value={form.code || ''} disabled />
          </div>
        )}
        <div className="mb-3">
          <label className="label">Nome *</label>
          <input className="input" value={form.name} onChange={(e) => up('name', e.target.value)} required />
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">E-mail</label>
            <input className="input" type="email" value={form.email || ''} onChange={(e) => up('email', e.target.value)} />
          </div>
          <div className="mb-3">
            <label className="label">Telefone</label>
            <input className="input" value={form.phone || ''} onChange={(e) => up('phone', e.target.value)} />
          </div>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">Cargo/Funcao</label>
            <input className="input" value={form.role || ''} onChange={(e) => up('role', e.target.value)} />
          </div>
          <div className="mb-4">
            <label className="label">Status</label>
            <select className="input" value={form.status} onChange={(e) => up('status', e.target.value)}>
              <option value="active">Ativo</option>
              <option value="inactive">Inativo</option>
            </select>
          </div>
        </div>
        <div className="flex justify-end gap-2">
          <button type="button" className="btn-secondary" onClick={onClose}>Cancelar</button>
          <button className="btn-primary" disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</button>
        </div>
      </form>
    </Modal>
  )
}

function PositionForm({ employee, onClose, onSaved }) {
  const [lat, setLat] = useState(employee.last_position?.latitude || '')
  const [lng, setLng] = useState(employee.last_position?.longitude || '')
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)

  function useBrowser() {
    if (!navigator.geolocation) {
      setError('Geolocalizacao nao suportada')
      return
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setLat(pos.coords.latitude)
        setLng(pos.coords.longitude)
      },
      () => setError('Nao foi possivel obter a localizacao')
    )
  }

  async function submit(event) {
    event.preventDefault()
    setSaving(true)
    setError('')
    try {
      await client.post(`/employees/${employee.id}/position`, {
        latitude: parseFloat(lat),
        longitude: parseFloat(lng),
      })
      onSaved()
    } catch (err) {
      setError(apiError(err))
      setSaving(false)
    }
  }

  return (
    <Modal open title={`Posicao - ${employee.name}`} onClose={onClose}>
      <form onSubmit={submit}>
        <ErrorBox message={error} />
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">Latitude</label>
            <input className="input" value={lat} onChange={(e) => setLat(e.target.value)} required />
          </div>
          <div className="mb-3">
            <label className="label">Longitude</label>
            <input className="input" value={lng} onChange={(e) => setLng(e.target.value)} required />
          </div>
        </div>
        <button type="button" className="btn-secondary mb-4" onClick={useBrowser}>Usar minha localizacao</button>
        <div className="flex justify-end gap-2">
          <button type="button" className="btn-secondary" onClick={onClose}>Cancelar</button>
          <button className="btn-primary" disabled={saving}>{saving ? 'Salvando...' : 'Registrar'}</button>
        </div>
      </form>
    </Modal>
  )
}
