import { useEffect, useState } from 'react'
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet'
import { Link } from 'react-router-dom'
import L from 'leaflet'
import client, { apiError } from '../api/client'
import { STATUS, STATUS_OPTIONS, fmtDate } from '../constants'
import { ErrorBox } from '../components/ui'

// Colored circular marker for a service order, keyed by status color.
function orderIcon(color) {
  return L.divIcon({
    className: 'os-marker',
    html: `<div style="background:${color};width:18px;height:18px;border-radius:50%;border:3px solid white;box-shadow:0 0 4px rgba(0,0,0,.5)"></div>`,
    iconSize: [18, 18],
    iconAnchor: [9, 9],
  })
}

// Distinct icon for an employee position.
const employeeIcon = L.divIcon({
  className: 'emp-marker',
  html: `<div style="font-size:22px;line-height:22px;filter:drop-shadow(0 1px 2px rgba(0,0,0,.5))">📍</div>`,
  iconSize: [22, 22],
  iconAnchor: [11, 22],
})

// FitBounds adjusts the viewport to contain all plotted points once.
function FitBounds({ points }) {
  const map = useMap()
  useEffect(() => {
    if (points.length > 0) {
      map.fitBounds(points, { padding: [50, 50], maxZoom: 14 })
    }
  }, [points, map])
  return null
}

export default function MapPage() {
  const [data, setData] = useState({ employees: [], service_orders: [] })
  const [statusFilter, setStatusFilter] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    const params = {}
    if (statusFilter) params.status = statusFilter
    client.get('/map/overview', { params })
      .then((res) => setData(res.data))
      .catch((err) => setError(apiError(err)))
  }, [statusFilter])

  const orders = (data.service_orders || []).filter((o) => o.latitude && o.longitude)
  const employees = data.employees || []
  const points = [
    ...orders.map((o) => [o.latitude, o.longitude]),
    ...employees.map((e) => [e.last_position.latitude, e.last_position.longitude]),
  ]
  const center = points[0] || [-15.78, -47.93] // Brasília fallback

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
        <h1 className="text-xl font-semibold">Mapa Operacional</h1>
        <div className="flex items-center gap-3">
          <select className="input w-48" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
            <option value="">Todos os status</option>
            {STATUS_OPTIONS.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
          </select>
        </div>
      </div>

      <ErrorBox message={error} />

      <div className="flex flex-wrap gap-3 mb-3 text-xs">
        {Object.entries(STATUS).map(([k, v]) => (
          <span key={k} className="flex items-center gap-1">
            <span className="inline-block w-3 h-3 rounded-full" style={{ background: v.color }} />
            {v.label}
          </span>
        ))}
        <span className="flex items-center gap-1">📍 Empregado</span>
      </div>

      <div className="card p-0 overflow-hidden" style={{ height: '70vh' }}>
        <MapContainer center={center} zoom={12} style={{ height: '100%', width: '100%' }}>
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <FitBounds points={points} />

          {orders.map((o) => (
            <Marker key={`o-${o.id}`} position={[o.latitude, o.longitude]} icon={orderIcon(STATUS[o.status]?.color || '#666')}>
              <Popup>
                <div className="text-sm">
                  <p className="font-mono text-xs text-slate-500">{o.number}</p>
                  <p className="font-semibold">{o.title}</p>
                  <p>Responsável: {o.employee_name || '—'}</p>
                  <p>Status: {STATUS[o.status]?.label}</p>
                  <p>Endereço: {o.address || '—'}</p>
                  <Link to={`/service-orders/${o.id}`} className="text-brand-600 underline">Ver detalhes</Link>
                </div>
              </Popup>
            </Marker>
          ))}

          {employees.map((e) => (
            <Marker key={`e-${e.id}`} position={[e.last_position.latitude, e.last_position.longitude]} icon={employeeIcon}>
              <Popup>
                <div className="text-sm">
                  <p className="font-semibold">{e.name}</p>
                  <p>Cargo: {e.role || '—'}</p>
                  <p className="text-xs text-slate-500">Atualizado: {fmtDate(e.last_position.recorded_at)}</p>
                </div>
              </Popup>
            </Marker>
          ))}
        </MapContainer>
      </div>
    </div>
  )
}
