export interface Profession {
  id: string;
  slug: string;
  title: string;
  category: string;
  ai_risk_score: number;
  avg_salary_usd: number;
  description: string;
  is_future_job: boolean;
  demand_index: number;
  growth_pct: number;
  updated_at: string;
}

export interface CareerTransition {
  id: string;
  from_profession: Profession;
  to_profession: Profession;
  match_score: number;
  transition_reason: string;
  avg_reskill_months: number;
}

export interface AnalysisRequest {
  profession_slug: string;
  location: string;
  years_exp: number;
  current_skills: string[];
}

export interface PivotSuggestion {
  target_profession: string;
  target_slug: string;
  match_score: number;
  reason: string;
  time_to_transition: string;
}

export interface AnalysisResult {
  profession_slug: string;
  profession_title: string;
  ai_risk_score: number;
  risk_level: 'Low' | 'Medium' | 'High';
  summary: string;
  threats: string[];
  opportunities: string[];
  recommended_pivots: PivotSuggestion[];
  timeline: string;
  skills_to_learn: string[];
  generated_at: string;
}
