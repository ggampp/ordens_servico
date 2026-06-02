import { useEffect, useState } from 'react'
import { MapContainer, Marker, TileLayer, useMap, useMapEvents } from 'react-leaflet'
import L from 'leaflet'
import { ErrorBox, Modal } from './ui'

const DEFAULT_CENTER = [-15.7801, -47.9292]

const selectedIcon = L.divIcon({
  className: 'selected-location-marker',
  html: '<div style="width:20px;height:20px;border-radius:50%;background:#2563eb;border:4px solid white;box-shadow:0 2px 8px rgba(15,23,42,.35)"></div>',
  iconSize: [20, 20],
  iconAnchor: [10, 10],
})

function parseCoordinate(value) {
  const n = Number(value)
  return Number.isFinite(n) ? n : null
}

function formatCoordinate(value) {
  const n = Number(value)
  return Number.isFinite(n) ? n.toFixed(6) : ''
}

export async function geocodeAddress(address) {
  const query = address.trim()
  if (!query) {
    throw new Error('Informe um endereco para buscar.')
  }

  const params = new URLSearchParams({
    format: 'json',
    limit: '1',
    q: query,
  })
  const res = await fetch(`https://nominatim.openstreetmap.org/search?${params}`)
  if (!res.ok) {
    throw new Error('Nao foi possivel buscar o endereco.')
  }
  const results = await res.json()
  if (!Array.isArray(results) || results.length === 0) {
    throw new Error('Endereco nao encontrado.')
  }
  const [item] = results
  return {
    address: item.display_name || query,
    latitude: Number(item.lat),
    longitude: Number(item.lon),
  }
}

function MapClick({ onSelect }) {
  useMapEvents({
    click(event) {
      onSelect({ lat: event.latlng.lat, lng: event.latlng.lng })
    },
  })
  return null
}

function MapCenter({ center, zoom }) {
  const map = useMap()
  useEffect(() => {
    map.setView(center, zoom)
    setTimeout(() => map.invalidateSize(), 0)
  }, [center, map, zoom])
  return null
}

export default function LocationPicker({ open, value, onClose, onConfirm }) {
  const [address, setAddress] = useState('')
  const [center, setCenter] = useState(DEFAULT_CENTER)
  const [selected, setSelected] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!open) return
    const lat = parseCoordinate(value?.latitude)
    const lng = parseCoordinate(value?.longitude)
    const point = lat !== null && lng !== null ? { lat, lng } : null
    setAddress(value?.address || '')
    setSelected(point)
    setCenter(point ? [point.lat, point.lng] : DEFAULT_CENTER)
    setError('')
  }, [open, value])

  async function searchAddress(event) {
    event?.preventDefault()
    setLoading(true)
    setError('')
    try {
      const location = await geocodeAddress(address)
      const point = { lat: location.latitude, lng: location.longitude }
      setAddress(location.address)
      setSelected(point)
      setCenter([point.lat, point.lng])
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  function confirm() {
    if (!selected) return
    onConfirm({
      address,
      latitude: formatCoordinate(selected.lat),
      longitude: formatCoordinate(selected.lng),
    })
  }

  const zoom = selected ? 16 : 5

  return (
    <Modal open={open} title="Selecionar posicao da OS" onClose={onClose} size="wide">
      <div>
        <ErrorBox message={error} />
        <form className="grid grid-cols-1 md:grid-cols-[1fr_auto] gap-3 mb-3" onSubmit={searchAddress}>
          <div>
            <label className="label">Endereco</label>
            <input className="input" value={address} onChange={(e) => setAddress(e.target.value)} />
          </div>
          <div className="flex items-end">
            <button type="submit" className="btn-secondary w-full md:w-auto" disabled={loading}>
              {loading ? 'Buscando...' : 'Centralizar'}
            </button>
          </div>
        </form>

        <div className="overflow-hidden rounded-md border border-slate-200" style={{ height: 380 }}>
          <MapContainer center={center} zoom={zoom} style={{ height: '100%', width: '100%' }}>
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <MapCenter center={center} zoom={zoom} />
            <MapClick onSelect={(point) => setSelected(point)} />
            {selected && <Marker position={[selected.lat, selected.lng]} icon={selectedIcon} />}
          </MapContainer>
        </div>

        <div className="mt-3 text-sm text-slate-600">
          {selected ? (
            <span className="font-mono">{formatCoordinate(selected.lat)}, {formatCoordinate(selected.lng)}</span>
          ) : (
            <span>Nenhuma posicao selecionada.</span>
          )}
        </div>

        <div className="flex justify-end gap-2 mt-4">
          <button type="button" className="btn-secondary" onClick={onClose}>Cancelar</button>
          <button type="button" className="btn-primary" onClick={confirm} disabled={!selected}>Confirmar posicao</button>
        </div>
      </div>
    </Modal>
  )
}
