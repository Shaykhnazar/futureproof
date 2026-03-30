import { useEffect } from 'react';
import { Routes, Route, Link, useLocation, Navigate } from 'react-router-dom';
import { Globe } from './components/Globe';
import { CareerAnalyzer } from './components/CareerAnalyzer';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { FutureJobsPage } from './pages/FutureJobsPage';
import { CareerCoachPage } from './pages/CareerCoachPage';
import { ProfilePage } from './pages/ProfilePage';
import { useWebSocket } from './hooks';
import { useGlobeStore, useCareerStore, useAuthStore } from './store';
import { citiesAPI, careersAPI } from './api';

function GlobePage() {
  const { selectedCity, setSelectedCity } = useGlobeStore();

  return (
    <div className="relative">
      <Globe />
      {selectedCity && (
        <div className="absolute top-6 right-6 w-80 rounded-xl overflow-hidden shadow-2xl border border-slate-700 bg-slate-900 text-white">
          <div
            className="h-1 w-full"
            style={{
              background:
                selectedCity.score >= 85 ? '#10b981' : selectedCity.score >= 70 ? '#f59e0b' : '#ef4444',
            }}
          />
          <div className="p-5">
            <div className="flex items-start justify-between mb-1">
              <h3 className="text-xl font-bold leading-tight">{selectedCity.name}</h3>
              <button
                onClick={() => setSelectedCity(null)}
                className="text-slate-400 hover:text-white ml-2 mt-0.5 text-lg leading-none"
              >
                ✕
              </button>
            </div>
            <p className="text-sm text-slate-400 mb-4">
              {selectedCity.country} • {selectedCity.region}
            </p>
            <div className="flex items-center justify-between mb-4 pb-4 border-b border-slate-700">
              <span className="text-sm text-slate-400">Opportunity Score</span>
              <span
                className="text-2xl font-bold"
                style={{
                  color: selectedCity.score >= 85 ? '#10b981' : selectedCity.score >= 70 ? '#f59e0b' : '#ef4444',
                }}
              >
                {selectedCity.score}/100
              </span>
            </div>
            <div className="grid grid-cols-2 gap-3 text-sm mb-4">
              {[
                { label: 'Job Growth', value: `${selectedCity.job_growth_pct}%` },
                { label: 'AI Investment', value: `${selectedCity.ai_investment}/100` },
                { label: 'Talent Demand', value: `${selectedCity.talent_demand}/100` },
                { label: 'Cost of Living', value: `${selectedCity.cost_of_living}/100` },
              ].map(({ label, value }) => (
                <div key={label} className="bg-slate-800 rounded-lg px-3 py-2">
                  <p className="text-slate-400 text-xs mb-0.5">{label}</p>
                  <p className="font-semibold text-white">{value}</p>
                </div>
              ))}
            </div>
            <p className="text-xs text-slate-500">
              Population: {selectedCity.population.toLocaleString()}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}

function Header() {
  const { isAuthenticated, user } = useAuthStore();
  const location = useLocation();

  const navLink = (to: string, label: string) => (
    <Link
      to={to}
      className={`px-3 py-2 rounded-lg text-sm transition-colors ${
        location.pathname === to
          ? 'bg-primary-600 text-white'
          : 'text-slate-400 hover:text-white'
      }`}
    >
      {label}
    </Link>
  );

  return (
    <header className="bg-slate-950 border-b border-slate-800 sticky top-0 z-50">
      <div className="container mx-auto px-6 py-3 flex items-center justify-between">
        <Link to="/" className="flex flex-col">
          <span className="text-xl font-bold bg-gradient-to-r from-primary-400 to-primary-600 bg-clip-text text-transparent">
            FutureProof
          </span>
          <span className="text-xs text-slate-500">Global Career Intelligence</span>
        </Link>

        <nav className="flex items-center gap-1">
          {navLink('/', 'Globe')}
          {navLink('/analyzer', 'Analyzer')}
          {navLink('/future-jobs', 'Future Jobs')}
          {navLink('/coach', 'AI Coach')}
          {isAuthenticated ? (
            navLink('/profile', user?.name?.split(' ')[0] || 'Profile')
          ) : (
            <>
              {navLink('/login', 'Sign In')}
              <Link
                to="/register"
                className="ml-1 px-3 py-2 rounded-lg text-sm bg-primary-600 hover:bg-primary-700 text-white transition-colors"
              >
                Sign Up
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}

// Pages that need the full-width globe get no header padding
const FULL_SCREEN_ROUTES = ['/'];

export default function App() {
  const { setCities } = useGlobeStore();
  const { setProfessions, setFutureProfessions } = useCareerStore();
  const { checkAuth } = useAuthStore();
  const location = useLocation();

  useWebSocket();

  useEffect(() => {
    checkAuth().catch(() => {});
  }, []);

  useEffect(() => {
    citiesAPI.getAllCities().then(setCities).catch(console.error);
    careersAPI.getAllProfessions().then(setProfessions).catch(console.error);
    careersAPI.getFutureProfessions().then(setFutureProfessions).catch(console.error);
  }, []);

  const isFullScreen = FULL_SCREEN_ROUTES.includes(location.pathname);

  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col">
      {!['/', '/login', '/register'].includes(location.pathname) || location.pathname === '/' ? (
        <Header />
      ) : null}

      <main className={isFullScreen ? '' : 'flex-1'}>
        <Routes>
          <Route path="/" element={<GlobePage />} />
          <Route path="/analyzer" element={<CareerAnalyzer />} />
          <Route path="/future-jobs" element={<FutureJobsPage />} />
          <Route path="/coach" element={<CareerCoachPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  );
}
