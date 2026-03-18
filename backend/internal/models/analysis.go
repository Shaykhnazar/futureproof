package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AIAnalysis represents a cached AI career analysis
type AIAnalysis struct {
	ID             uuid.UUID       `json:"id"`
	UserID         *uuid.UUID      `json:"user_id"`
	ProfessionSlug string          `json:"profession_slug"`
	Location       string          `json:"location"`
	RequestHash    string          `json:"request_hash"`
	Result         json.RawMessage `json:"result"`
	ModelUsed      string          `json:"model_used"`
	TokensUsed     int             `json:"tokens_used"`
	CreatedAt      time.Time       `json:"created_at"`
}

// AnalysisRequest represents an AI analysis request
type AnalysisRequest struct {
	ProfessionSlug string `json:"profession_slug"`
	Location       string `json:"location"`
	YearsExp       int    `json:"years_exp"`
	CurrentSkills  []string `json:"current_skills"`
}

// AnalysisResult represents the structured AI analysis response
type AnalysisResult struct {
	ProfessionSlug    string         `json:"profession_slug"`
	ProfessionTitle   string         `json:"profession_title"`
	AIRiskScore       int            `json:"ai_risk_score"`
	RiskLevel         string         `json:"risk_level"` // "Low", "Medium", "High"
	Summary           string         `json:"summary"`
	Threats           []string       `json:"threats"`
	Opportunities     []string       `json:"opportunities"`
	RecommendedPivots []PivotSuggestion `json:"recommended_pivots"`
	Timeline          string         `json:"timeline"`
	SkillsToLearn     []string       `json:"skills_to_learn"`
	GeneratedAt       time.Time      `json:"generated_at"`
}

// PivotSuggestion represents a career pivot recommendation
type PivotSuggestion struct {
	TargetProfession string  `json:"target_profession"`
	TargetSlug       string  `json:"target_slug"`
	MatchScore       int     `json:"match_score"`
	Reason           string  `json:"reason"`
	TimeToTransition string  `json:"time_to_transition"`
}

// WebSocketMessage represents real-time updates
type WebSocketMessage struct {
	Type      string          `json:"type"` // "city_update", "job_update", etc.
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// GlobeUpdate represents a real-time globe visualization update
type GlobeUpdate struct {
	CityID       uuid.UUID `json:"city_id"`
	CityName     string    `json:"city_name"`
	NewScore     int       `json:"new_score"`
	JobGrowthPct float64   `json:"job_growth_pct"`
	UpdateReason string    `json:"update_reason"`
}
