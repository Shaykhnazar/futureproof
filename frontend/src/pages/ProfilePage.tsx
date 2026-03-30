import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore, useCareerStore } from '../store';
import { authAPI, careersAPI } from '../api';
import type { Profession } from '../types';

export function ProfilePage() {
  const { user, logout, checkAuth } = useAuthStore();
  const { savedCareers, setSavedCareers } = useCareerStore();
  const navigate = useNavigate();
  const [isSaving, setIsSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [yearsExp, setYearsExp] = useState(user?.profile?.years_exp || 0);
  const [skills, setSkills] = useState((user?.profile?.skills || []).join(', '));

  useEffect(() => {
    if (!user) {
      checkAuth().catch(() => navigate('/login'));
    } else {
      setYearsExp(user.profile?.years_exp || 0);
      setSkills((user.profile?.skills || []).join(', '));
    }
  }, [user]);

  useEffect(() => {
    careersAPI.getSavedCareers()
      .then((careers: Profession[]) => setSavedCareers(careers))
      .catch(() => {});
  }, []);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      await authAPI.updateProfile({
        years_exp: yearsExp,
        skills: skills.split(',').map((s) => s.trim()).filter(Boolean),
      });
      await checkAuth();
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } finally {
      setIsSaving(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (!user) return null;

  return (
    <div className="max-w-3xl mx-auto px-6 py-10 space-y-8">
      <h2 className="text-3xl font-bold">Your Profile</h2>

      {/* Account info (read-only) */}
      <div className="bg-slate-900 border border-slate-700 rounded-xl p-6 space-y-4">
        <h3 className="font-semibold text-lg">Account</h3>
        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <p className="text-xs text-slate-500 mb-1">Name</p>
            <p className="text-white font-medium">{user.name}</p>
          </div>
          <div>
            <p className="text-xs text-slate-500 mb-1">Email</p>
            <p className="text-white font-medium">{user.email}</p>
          </div>
        </div>
      </div>

      {/* Career profile form */}
      <form onSubmit={handleSave} className="bg-slate-900 border border-slate-700 rounded-xl p-6 space-y-5">
        <h3 className="font-semibold text-lg">Career Profile</h3>

        <div>
          <label className="block text-sm text-slate-400 mb-1.5">Years of Experience</label>
          <input
            type="number"
            min={0}
            max={50}
            value={yearsExp}
            onChange={(e) => setYearsExp(parseInt(e.target.value) || 0)}
            className="w-full bg-slate-800 border border-slate-600 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-primary-500"
          />
        </div>

        <div>
          <label className="block text-sm text-slate-400 mb-1.5">Current Skills (comma-separated)</label>
          <input
            type="text"
            value={skills}
            onChange={(e) => setSkills(e.target.value)}
            placeholder="e.g. Python, SQL, Project Management"
            className="w-full bg-slate-800 border border-slate-600 rounded-lg px-4 py-3 text-white placeholder-slate-500 focus:outline-none focus:border-primary-500"
          />
        </div>

        <div className="flex items-center gap-4">
          <button
            type="submit"
            disabled={isSaving}
            className="bg-primary-600 hover:bg-primary-700 disabled:opacity-50 text-white px-6 py-2.5 rounded-lg font-medium transition-colors"
          >
            {isSaving ? 'Saving…' : 'Save Profile'}
          </button>
          {saved && <span className="text-emerald-400 text-sm">Saved!</span>}
        </div>
      </form>

      {/* Saved careers */}
      <div className="bg-slate-900 border border-slate-700 rounded-xl p-6">
        <h3 className="font-semibold text-lg mb-4">Saved Careers ({savedCareers.length})</h3>
        {savedCareers.length === 0 ? (
          <p className="text-slate-400 text-sm">
            No saved careers yet. Use the Career Analyzer to explore options.
          </p>
        ) : (
          <div className="space-y-3">
            {savedCareers.map((career) => (
              <div key={career.id} className="flex items-center justify-between bg-slate-800 rounded-lg px-4 py-3">
                <div>
                  <p className="font-medium text-white">{career.title}</p>
                  <p className="text-sm text-slate-400 capitalize">{career.category}</p>
                </div>
                <span
                  className={`text-sm font-medium ${
                    career.ai_risk_score >= 70
                      ? 'text-red-400'
                      : career.ai_risk_score >= 40
                      ? 'text-yellow-400'
                      : 'text-emerald-400'
                  }`}
                >
                  Risk {career.ai_risk_score}/100
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      <button
        onClick={handleLogout}
        className="text-red-400 hover:text-red-300 text-sm font-medium transition-colors"
      >
        Sign out
      </button>
    </div>
  );
}
