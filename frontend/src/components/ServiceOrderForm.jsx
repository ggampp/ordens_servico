import { useMemo, useState } from 'react'
import client, { apiError } from '../api/client'
import { PRIORITY_OPTIONS } from '../constants'
import { ErrorBox, Modal } from './ui'
import LocationPicker, { geocodeAddress } from './LocationPicker'

const emptyOrder = {
  title: '',
  description: '',
  priority: 'medium',
  employee_id: '',
  address: '',
  latitude: '',
  longitude: '',
  due_at: '',
  notes: '',
}

function formatCoordinate(value) {
  const n = Number(value)
  return Number.isFinite(n) ? n.toFixed(6) : ''
}

function toDateTimeLocal(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 16)
}

function toForm(initial) {
  if (!initial) return emptyOrder
  return {
    title: initial.title || '',
    description: initial.description || '',
    priority: initial.priority || 'medium',
    employee_id: initial.employee_id ? String(initial.employee_id) : '',
    address: initial.address || '',
    latitude: initial.latitude == null ? '' : String(initial.latitude),
    longitude: initial.longitude == null ? '' : String(initial.longitude),
    due_at: toDateTimeLocal(initial.due_at),
    notes: initial.notes || '',
  }
}

export default function ServiceOrderForm({ employees = [], initial, onClose, onSaved }) {
  const isEdit = Boolean(initial?.id)
  const initialForm = useMemo(() => toForm(initial), [initial])
  const [form, setForm] = useState(initialForm)
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)
  const [locating, setLocating] = useState(false)
  const [mapOpen, setMapOpen] = useState(false)

  function up(key, value) {
    setForm((current) => ({ ...current, [key]: value }))
  }

  function applyLocation(location) {
    setForm((current) => ({
      ...current,
      address: location.address || current.address,
      latitude: formatCoordinate(location.latitude),
      longitude: formatCoordinate(location.longitude),
    }))
  }

  async function findCoordinates() {
    setLocating(true)
    setError('')
    try {
      const location = await geocodeAddress(form.address)
      applyLocation(location)
    } catch (err) {
      setError(err.message)
    } finally {
      setLocating(false)
    }
  }

  function payload() {
    const body = {
      title: form.title,
      description: form.description || null,
      priority: form.priority,
      address: form.address || null,
      latitude: form.latitude ? parseFloat(form.latitude) : null,
      longitude: form.longitude ? parseFloat(form.longitude) : null,
      due_at: form.due_at ? new Date(form.due_at).toISOString() : null,
      notes: form.notes || null,
    }
    if (!isEdit) {
      body.employee_id = form.employee_id ? Number(form.employee_id) : null
    }
    return body
  }

  async function submit(event) {
    event.preventDefault()
    setSaving(true)
    setError('')
    try {
      if (isEdit) {
        await client.put(`/service-orders/${initial.id}`, payload())
      } else {
        await client.post('/service-orders', payload())
      }
      onSaved()
    } catch (err) {
      setError(apiError(err))
      setSaving(false)
    }
  }

  return (
    <Modal open title={isEdit ? 'Editar Ordem de Servico' : 'Nova Ordem de Servico'} onClose={onClose}>
      <form onSubmit={submit}>
        <ErrorBox message={error} />

        <div className="mb-3">
          <label className="label">Titulo *</label>
          <input className="input" value={form.title} onChange={(e) => up('title', e.target.value)} required />
        </div>

        <div className="mb-3">
          <label className="label">Descricao</label>
          <textarea className="input" rows={2} value={form.description} onChange={(e) => up('description', e.target.value)} />
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">Prioridade</label>
            <select className="input" value={form.priority} onChange={(e) => up('priority', e.target.value)}>
              {PRIORITY_OPTIONS.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
            </select>
          </div>
          {!isEdit && (
            <div className="mb-3">
              <label className="label">Responsavel</label>
              <select className="input" value={form.employee_id} onChange={(e) => up('employee_id', e.target.value)}>
                <option value="">Nao atribuir</option>
                {employees.map((employee) => <option key={employee.id} value={employee.id}>{employee.name}</option>)}
              </select>
            </div>
          )}
        </div>

        <div className="mb-3">
          <label className="label">Endereco</label>
          <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto_auto] gap-2">
            <input className="input" value={form.address} onChange={(e) => up('address', e.target.value)} />
            <button type="button" className="btn-secondary" onClick={findCoordinates} disabled={locating}>
              {locating ? 'Buscando...' : 'Buscar coordenadas'}
            </button>
            <button type="button" className="btn-secondary" onClick={() => setMapOpen(true)}>
              Abrir mapa
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div className="mb-3">
            <label className="label">Latitude</label>
            <input className="input" value={form.latitude} onChange={(e) => up('latitude', e.target.value)} />
          </div>
          <div className="mb-3">
            <label className="label">Longitude</label>
            <input className="input" value={form.longitude} onChange={(e) => up('longitude', e.target.value)} />
          </div>
        </div>

        <div className="mb-3">
          <label className="label">Data prevista</label>
          <input type="datetime-local" className="input" value={form.due_at} onChange={(e) => up('due_at', e.target.value)} />
        </div>

        <div className="mb-4">
          <label className="label">Observacoes</label>
          <textarea className="input" rows={2} value={form.notes} onChange={(e) => up('notes', e.target.value)} />
        </div>

        <div className="flex justify-end gap-2">
          <button type="button" className="btn-secondary" onClick={onClose}>Cancelar</button>
          <button className="btn-primary" disabled={saving}>{saving ? 'Salvando...' : (isEdit ? 'Salvar' : 'Criar')}</button>
        </div>
      </form>

      <LocationPicker
        open={mapOpen}
        value={form}
        onClose={() => setMapOpen(false)}
        onConfirm={(location) => {
          applyLocation(location)
          setMapOpen(false)
        }}
      />
    </Modal>
  )
}
