export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  auth_provider: string;
  created_at: string;
}

export interface UserProfile {
  user_id: string;
  current_job_id?: string;
  city_id?: string;
  years_exp: number;
  education?: string;
  target_job_id?: string;
  skills: string[];
  updated_at: string;
}

export interface UserWithProfile {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  auth_provider: string;
  created_at: string;
  profile?: UserProfile;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  name: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}
