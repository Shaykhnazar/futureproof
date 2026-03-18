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
              <div className="absolute top-6 right-6 w-80 card">
                <h3 className="text-xl font-bold mb-2">{selectedCity.name}</h3>
                <p className="text-sm text-gray-400 mb-4">
                  {selectedCity.country} • {selectedCity.region}
                </p>

                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-400">Opportunity Score</span>
                    <span className="text-2xl font-bold text-primary-500">
                      {selectedCity.score}/100
                    </span>
                  </div>

                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div>
                      <span className="text-gray-400">Job Growth</span>
                      <p className="font-semibold">{selectedCity.job_growth_pct}%</p>
                    </div>
                    <div>
                      <span className="text-gray-400">AI Investment</span>
                      <p className="font-semibold">{selectedCity.ai_investment}/100</p>
                    </div>
                    <div>
                      <span className="text-gray-400">Talent Demand</span>
                      <p className="font-semibold">{selectedCity.talent_demand}/100</p>
                    </div>
                    <div>
                      <span className="text-gray-400">Cost of Living</span>
                      <p className="font-semibold">{selectedCity.cost_of_living}/100</p>
                    </div>
                  </div>

                  <div className="pt-3 border-t border-white/10">
                    <p className="text-xs text-gray-500">
                      Population: {selectedCity.population.toLocaleString()}
                    </p>
                  </div>
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
