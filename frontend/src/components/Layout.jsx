import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { ROLE_LABELS } from '../constants'

export default function Layout() {
  const { user, logout, canManage } = useAuth()
  const navigate = useNavigate()

  const links = [
    { to: '/map', label: 'Mapa' },
    ...(canManage ? [{ to: '/dashboard', label: 'Dashboard' }] : []),
    ...(canManage ? [{ to: '/employees', label: 'Empregados' }] : []),
    { to: '/service-orders', label: 'Ordens de Serviço' },
  ]

  function handleLogout() {
    logout()
    navigate('/login')
  }

  return (
    <div className="min-h-full flex flex-col">
      <header className="bg-brand-700 text-white shadow">
        <div className="max-w-7xl mx-auto px-4 flex items-center justify-between h-14">
          <div className="flex items-center gap-6">
            <span className="font-semibold tracking-tight">🛠️ Ordens de Serviço</span>
            <nav className="hidden md:flex gap-1">
              {links.map((l) => (
                <NavLink
                  key={l.to}
                  to={l.to}
                  className={({ isActive }) =>
                    `px-3 py-1.5 rounded text-sm ${isActive ? 'bg-brand-600' : 'hover:bg-brand-600/60'}`
                  }
                >
                  {l.label}
                </NavLink>
              ))}
            </nav>
          </div>
          <div className="flex items-center gap-3 text-sm">
            <span className="hidden sm:inline">
              {user?.name} · {ROLE_LABELS[user?.role]}
            </span>
            <button onClick={handleLogout} className="rounded bg-brand-600 px-3 py-1.5 hover:bg-brand-500">
              Sair
            </button>
          </div>
        </div>
        {/* Mobile nav */}
        <nav className="md:hidden flex gap-1 px-4 pb-2 overflow-x-auto">
          {links.map((l) => (
            <NavLink
              key={l.to}
              to={l.to}
              className={({ isActive }) =>
                `px-3 py-1 rounded text-sm whitespace-nowrap ${isActive ? 'bg-brand-600' : 'bg-brand-600/40'}`
              }
            >
              {l.label}
            </NavLink>
          ))}
        </nav>
      </header>
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  )
}
