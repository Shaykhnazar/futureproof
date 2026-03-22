import { useEffect, useState } from 'react';
import { Globe } from './components/Globe';
import { CareerAnalyzer } from './components/CareerAnalyzer';
import { useWebSocket } from './hooks';
import { useGlobeStore, useCareerStore } from './store';
import { citiesAPI, careersAPI } from './api';

function App() {
  const [activeTab, setActiveTab] = useState<'globe' | 'analyzer'>('globe');
  const { setCities, selectedCity } = useGlobeStore();
  const { setProfessions, setFutureProfessions } = useCareerStore();

  // Connect to WebSocket for real-time updates
  useWebSocket();

  // Load initial data
  useEffect(() => {
    const loadData = async () => {
      try {
        // Load cities
        const cities = await citiesAPI.getAllCities();
        setCities(cities);

        // Load professions
        const professions = await careersAPI.getAllProfessions();
        setProfessions(professions);

        // Load future professions
        const futureProfessions = await careersAPI.getFutureProfessions();
        setFutureProfessions(futureProfessions);
      } catch (error) {
        console.error('Failed to load initial data:', error);
      }
    };

    loadData();
  }, [setCities, setProfessions, setFutureProfessions]);

  return (
    <div className="min-h-screen bg-slate-950">
      {/* Header */}
      <header className="glass border-b border-white/10">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold bg-gradient-to-r from-primary-400 to-primary-600 bg-clip-text text-transparent">
                FutureProof
              </h1>
              <p className="text-sm text-gray-400">
                Global Career Intelligence Platform
              </p>
            </div>

            <nav className="flex gap-4">
              <button
                onClick={() => setActiveTab('globe')}
                className={`px-4 py-2 rounded-lg transition-colors ${
                  activeTab === 'globe'
                    ? 'bg-primary-600 text-white'
                    : 'text-gray-400 hover:text-white'
                }`}
              >
                Globe
              </button>
              <button
                onClick={() => setActiveTab('analyzer')}
                className={`px-4 py-2 rounded-lg transition-colors ${
                  activeTab === 'analyzer'
                    ? 'bg-primary-600 text-white'
                    : 'text-gray-400 hover:text-white'
                }`}
              >
                Analyzer
              </button>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main>
        {activeTab === 'globe' ? (
          <div className="relative">
            <Globe />

            {/* City Info Panel */}
            {selectedCity && (
              <div className="absolute top-6 right-6 w-80 rounded-xl overflow-hidden shadow-2xl border border-slate-700 bg-slate-900 text-white">
                {/* Score colour bar */}
                <div
                  className="h-1 w-full"
                  style={{
                    background:
                      selectedCity.score >= 85
                        ? '#10b981'
                        : selectedCity.score >= 70
                        ? '#f59e0b'
                        : '#ef4444',
                  }}
                />

                <div className="p-5">
                  {/* Header */}
                  <div className="flex items-start justify-between mb-1">
                    <h3 className="text-xl font-bold leading-tight">{selectedCity.name}</h3>
                    <button
                      onClick={() => useGlobeStore.getState().setSelectedCity(null)}
                      className="text-slate-400 hover:text-white ml-2 mt-0.5 leading-none text-lg"
                      aria-label="Close"
                    >
                      ✕
                    </button>
                  </div>
                  <p className="text-sm text-slate-400 mb-4">
                    {selectedCity.country} • {selectedCity.region}
                  </p>

                  {/* Score */}
                  <div className="flex items-center justify-between mb-4 pb-4 border-b border-slate-700">
                    <span className="text-sm text-slate-400">Opportunity Score</span>
                    <span
                      className="text-2xl font-bold"
                      style={{
                        color:
                          selectedCity.score >= 85
                            ? '#10b981'
                            : selectedCity.score >= 70
                            ? '#f59e0b'
                            : '#ef4444',
                      }}
                    >
                      {selectedCity.score}/100
                    </span>
                  </div>

                  {/* Stats grid */}
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

                  {/* Population */}
                  <p className="text-xs text-slate-500">
                    Population: {selectedCity.population.toLocaleString()}
                  </p>
                </div>
              </div>
            )}
          </div>
        ) : (
          <CareerAnalyzer />
        )}
      </main>

      {/* Footer */}
      <footer className="glass border-t border-white/10 mt-auto">
        <div className="container mx-auto px-6 py-4 text-center text-sm text-gray-400">
          <p>
            Navigate AI-driven career disruption with real-time intelligence
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
