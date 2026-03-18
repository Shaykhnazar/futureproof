import { useState } from 'react';
import { useCareerAnalysis } from '../hooks';
import { useCareerStore } from '../store';
import type { AnalysisRequest } from '../types';

export function CareerAnalyzer() {
  const { professions } = useCareerStore();
  const { analyzeCareer, isLoading, error } = useCareerAnalysis();
  const { analysis } = useCareerStore();

  const [formData, setFormData] = useState<AnalysisRequest>({
    profession_slug: '',
    location: '',
    years_exp: 0,
    current_skills: [],
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await analyzeCareer(formData);
  };

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'Low':
        return 'text-success-500';
      case 'Medium':
        return 'text-warning-500';
      case 'High':
        return 'text-danger-500';
      default:
        return 'text-gray-400';
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h2 className="text-3xl font-bold mb-6">AI Career Risk Analyzer</h2>

      <form onSubmit={handleSubmit} className="card mb-8">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">Profession</label>
            <select
              className="w-full p-3 bg-slate-800 rounded-lg border border-slate-700"
              value={formData.profession_slug}
              onChange={(e) =>
                setFormData({ ...formData, profession_slug: e.target.value })
              }
              required
            >
              <option value="">Select a profession</option>
              {professions.map((prof) => (
                <option key={prof.id} value={prof.slug}>
                  {prof.title}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Location</label>
            <input
              type="text"
              className="w-full p-3 bg-slate-800 rounded-lg border border-slate-700"
              placeholder="e.g., San Francisco"
              value={formData.location}
              onChange={(e) =>
                setFormData({ ...formData, location: e.target.value })
              }
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">
              Years of Experience
            </label>
            <input
              type="number"
              className="w-full p-3 bg-slate-800 rounded-lg border border-slate-700"
              value={formData.years_exp}
              onChange={(e) =>
                setFormData({
                  ...formData,
                  years_exp: parseInt(e.target.value),
                })
              }
              min="0"
            />
          </div>

          <button type="submit" className="btn-primary w-full" disabled={isLoading}>
            {isLoading ? 'Analyzing...' : 'Analyze Career'}
          </button>

          {error && <p className="text-danger-500 text-sm">{error}</p>}
        </div>
      </form>

      {analysis && (
        <div className="card">
          <h3 className="text-2xl font-bold mb-4">{analysis.profession_title}</h3>

          <div className="grid grid-cols-2 gap-4 mb-6">
            <div className="p-4 bg-slate-800 rounded-lg">
              <p className="text-sm text-gray-400">AI Risk Score</p>
              <p className="text-3xl font-bold">{analysis.ai_risk_score}/100</p>
            </div>
            <div className="p-4 bg-slate-800 rounded-lg">
              <p className="text-sm text-gray-400">Risk Level</p>
              <p className={`text-3xl font-bold ${getRiskColor(analysis.risk_level)}`}>
                {analysis.risk_level}
              </p>
            </div>
          </div>

          <p className="text-gray-300 mb-6">{analysis.summary}</p>

          <div className="grid md:grid-cols-2 gap-6 mb-6">
            <div>
              <h4 className="font-bold mb-3 text-danger-500">Threats</h4>
              <ul className="space-y-2">
                {analysis.threats.map((threat, i) => (
                  <li key={i} className="text-sm text-gray-300">
                    • {threat}
                  </li>
                ))}
              </ul>
            </div>

            <div>
              <h4 className="font-bold mb-3 text-success-500">Opportunities</h4>
              <ul className="space-y-2">
                {analysis.opportunities.map((opp, i) => (
                  <li key={i} className="text-sm text-gray-300">
                    • {opp}
                  </li>
                ))}
              </ul>
            </div>
          </div>

          <div className="mb-6">
            <h4 className="font-bold mb-3">Recommended Career Pivots</h4>
            <div className="space-y-3">
              {analysis.recommended_pivots.map((pivot, i) => (
                <div key={i} className="p-4 bg-slate-800 rounded-lg">
                  <div className="flex justify-between items-start mb-2">
                    <h5 className="font-semibold">{pivot.target_profession}</h5>
                    <span className="text-primary-500 font-bold">
                      {pivot.match_score}% match
                    </span>
                  </div>
                  <p className="text-sm text-gray-400">{pivot.reason}</p>
                  <p className="text-xs text-gray-500 mt-2">
                    Transition time: {pivot.time_to_transition}
                  </p>
                </div>
              ))}
            </div>
          </div>

          <div>
            <h4 className="font-bold mb-3">Skills to Learn</h4>
            <div className="flex flex-wrap gap-2">
              {analysis.skills_to_learn.map((skill, i) => (
                <span
                  key={i}
                  className="px-3 py-1 bg-primary-600/20 text-primary-400 rounded-full text-sm"
                >
                  {skill}
                </span>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
