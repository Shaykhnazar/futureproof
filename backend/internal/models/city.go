package models

import (
	"time"

	"github.com/google/uuid"
)

// City represents a global city
type City struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Country    string    `json:"country"`
	Region     string    `json:"region"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Population int       `json:"population"`
	Timezone   string    `json:"timezone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CityScore represents opportunity metrics for a city
type CityScore struct {
	ID            uuid.UUID `json:"id"`
	CityID        uuid.UUID `json:"city_id"`
	Score         int       `json:"score"` // Overall score 0-100
	JobGrowthPct  float64   `json:"job_growth_pct"`
	RemoteScore   int       `json:"remote_score"`
	AIInvestment  int       `json:"ai_investment"`
	TalentDemand  int       `json:"talent_demand"`
	CostOfLiving  int       `json:"cost_of_living"`
	SnapshotDate  time.Time `json:"snapshot_date"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
}

// CityWithScore combines city and its latest score
type CityWithScore struct {
	City
	Score        int     `json:"score"`
	JobGrowthPct float64 `json:"job_growth_pct"`
	RemoteScore  int     `json:"remote_score"`
	AIInvestment int     `json:"ai_investment"`
	TalentDemand int     `json:"talent_demand"`
	CostOfLiving int     `json:"cost_of_living"`
}

// CityTopProfession represents top professions in a city
type CityTopProfession struct {
	CityID       uuid.UUID `json:"city_id"`
	ProfessionID uuid.UUID `json:"profession_id"`
	Rank         int       `json:"rank"`
	LocalDemand  int       `json:"local_demand"`
	SnapshotDate time.Time `json:"snapshot_date"`
}
