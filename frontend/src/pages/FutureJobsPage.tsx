import { useCareerStore } from '../store';

const riskColor = (score: number) =>
  score >= 70 ? 'text-red-400' : score >= 40 ? 'text-yellow-400' : 'text-emerald-400';

const riskLabel = (score: number) =>
  score >= 70 ? 'High Risk' : score >= 40 ? 'Medium Risk' : 'Low Risk';

export function FutureJobsPage() {
  const { futureProfessions, professions } = useCareerStore();

  const futureJobs = futureProfessions.length > 0 ? futureProfessions : [];
  const topRiskyJobs = [...professions]
    .sort((a, b) => b.ai_risk_score - a.ai_risk_score)
    .slice(0, 5);

  return (
    <div className="max-w-5xl mx-auto px-6 py-10">
      <div className="mb-10">
        <h2 className="text-3xl font-bold mb-2">Future Jobs of 2030+</h2>
        <p className="text-slate-400">Emerging roles with the highest demand growth in the AI era</p>
      </div>

      {/* Future / Emerging roles */}
      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-5 mb-14">
        {futureJobs.map((job) => (
          <div key={job.id} className="bg-slate-900 border border-slate-700 rounded-xl p-5 hover:border-primary-500 transition-colors">
            <div className="flex items-start justify-between mb-3">
              <h3 className="font-semibold text-white leading-tight">{job.title}</h3>
              <span className="text-xs bg-primary-600/20 text-primary-400 px-2 py-0.5 rounded-full whitespace-nowrap ml-2">
                Future
              </span>
            </div>

            <p className="text-slate-400 text-sm mb-4 line-clamp-2">{job.description}</p>

            <div className="space-y-1.5 text-sm">
              {job.avg_salary_usd > 0 && (
                <div className="flex justify-between">
                  <span className="text-slate-500">Avg Salary</span>
                  <span className="text-white font-medium">${job.avg_salary_usd.toLocaleString()}</span>
                </div>
              )}
              {job.growth_pct > 0 && (
                <div className="flex justify-between">
                  <span className="text-slate-500">Growth</span>
                  <span className="text-emerald-400 font-medium">+{job.growth_pct}%</span>
                </div>
              )}
              <div className="flex justify-between">
                <span className="text-slate-500">AI Risk</span>
                <span className={`font-medium ${riskColor(job.ai_risk_score)}`}>
                  {riskLabel(job.ai_risk_score)}
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Most at-risk jobs */}
      <div>
        <h3 className="text-xl font-bold mb-4 text-red-400">Most At-Risk from AI Automation</h3>
        <div className="bg-slate-900 border border-slate-700 rounded-xl overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-slate-700 text-slate-400">
                <th className="text-left px-5 py-3 font-medium">Profession</th>
                <th className="text-left px-5 py-3 font-medium">Category</th>
                <th className="text-right px-5 py-3 font-medium">AI Risk Score</th>
              </tr>
            </thead>
            <tbody>
              {topRiskyJobs.map((job, i) => (
                <tr key={job.id} className={i < topRiskyJobs.length - 1 ? 'border-b border-slate-800' : ''}>
                  <td className="px-5 py-3 text-white font-medium">{job.title}</td>
                  <td className="px-5 py-3 text-slate-400 capitalize">{job.category}</td>
                  <td className="px-5 py-3 text-right">
                    <span className={`font-bold ${riskColor(job.ai_risk_score)}`}>
                      {job.ai_risk_score}/100
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
