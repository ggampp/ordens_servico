import { useEffect, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts'
import client, { apiError } from '../api/client'
import { Spinner, ErrorBox } from '../components/ui'

function Stat({ label, value, color }) {
  return (
    <div className="card flex flex-col">
      <span className="text-3xl font-bold" style={{ color }}>{value}</span>
      <span className="text-sm text-slate-500 mt-1">{label}</span>
    </div>
  )
}

export default function Dashboard() {
  const [data, setData] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    client
      .get('/dashboard')
      .then((res) => setData(res.data))
      .catch((err) => setError(apiError(err)))
  }, [])

  if (error) return <ErrorBox message={error} />
  if (!data) return <Spinner />

  const byEmployee = (data.orders_by_employee || []).map((e) => ({
    name: e.employee_name,
    total: e.count,
  }))

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">Dashboard</h1>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4 mb-6">
        <Stat label="OS Abertas" value={data.open_orders} color="#dc2626" />
        <Stat label="Atribuídas" value={data.assigned_orders} color="#ea580c" />
        <Stat label="Em Atendimento" value={data.in_progress_orders} color="#2563eb" />
        <Stat label="Concluídas" value={data.completed_orders} color="#16a34a" />
        <Stat label="Empregados Ativos" value={data.active_employees} color="#0f172a" />
      </div>

      <div className="card">
        <h2 className="font-medium mb-4">Ordens por Responsável</h2>
        {byEmployee.length === 0 ? (
          <p className="text-sm text-slate-400">Sem dados.</p>
        ) : (
          <ResponsiveContainer width="100%" height={320}>
            <BarChart data={byEmployee} margin={{ bottom: 40 }}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" angle={-25} textAnchor="end" interval={0} height={60} fontSize={12} />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Bar dataKey="total" fill="#2563eb" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        )}
      </div>
    </div>
  )
}
