export interface City {
  id: string;
  name: string;
  country: string;
  region: string;
  lat: number;
  lng: number;
  population: number;
  timezone: string;
  created_at: string;
  updated_at: string;
}

export interface CityWithScore extends City {
  score: number;
  job_growth_pct: number;
  remote_score: number;
  ai_investment: number;
  talent_demand: number;
  cost_of_living: number;
}

export interface GlobeUpdate {
  city_id: string;
  city_name: string;
  new_score: number;
  job_growth_pct: number;
  update_reason: string;
}

export interface WebSocketMessage {
  type: string;
  timestamp: string;
  data: unknown;
}
