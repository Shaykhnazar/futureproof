package models

import (
	"time"

	"github.com/google/uuid"
)

// Profession represents a career or job role
type Profession struct {
	ID           uuid.UUID `json:"id"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	AIRiskScore  int       `json:"ai_risk_score"`
	AvgSalaryUSD int       `json:"avg_salary_usd"`
	Description  string    `json:"description"`
	IsFutureJob  bool      `json:"is_future_job"`
	DemandIndex  int       `json:"demand_index"`
	GrowthPct    float64   `json:"growth_pct"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Skill represents a professional skill
type Skill struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	IsAIProof  bool      `json:"is_ai_proof"`
}

// ProfessionSkill links professions to skills
type ProfessionSkill struct {
	ProfessionID uuid.UUID `json:"profession_id"`
	SkillID      uuid.UUID `json:"skill_id"`
	Importance   int       `json:"importance"` // 1-5 scale
	IsAtRisk     bool      `json:"is_at_risk"`
}

// CareerTransition represents a career pivot path
type CareerTransition struct {
	ID               uuid.UUID `json:"id"`
	FromProfession   uuid.UUID `json:"from_profession"`
	ToProfession     uuid.UUID `json:"to_profession"`
	MatchScore       int       `json:"match_score"`        // 0-100
	TransitionReason string    `json:"transition_reason"`
	AvgReskillMonths int       `json:"avg_reskill_months"`
}

// CareerTransitionWithDetails includes profession details
type CareerTransitionWithDetails struct {
	ID                uuid.UUID  `json:"id"`
	FromProfession    Profession `json:"from_profession"`
	ToProfession      Profession `json:"to_profession"`
	MatchScore        int        `json:"match_score"`
	TransitionReason  string     `json:"transition_reason"`
	AvgReskillMonths  int        `json:"avg_reskill_months"`
}

// SavedCareer represents a user's saved career
type SavedCareer struct {
	UserID       uuid.UUID `json:"user_id"`
	ProfessionID uuid.UUID `json:"profession_id"`
	Notes        string    `json:"notes"`
	SavedAt      time.Time `json:"saved_at"`
}
