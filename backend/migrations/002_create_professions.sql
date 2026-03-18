-- ============================================================
-- FILE: migrations/002_create_professions.sql
-- ============================================================

CREATE TABLE professions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug            VARCHAR(100) UNIQUE NOT NULL,
    title           VARCHAR(200) NOT NULL,
    category        VARCHAR(50)  NOT NULL,
    ai_risk_score   INTEGER CHECK (ai_risk_score BETWEEN 0 AND 100),
    avg_salary_usd  INTEGER,
    description     TEXT,
    is_future_job   BOOLEAN DEFAULT false,
    demand_index    INTEGER CHECK (demand_index BETWEEN 0 AND 100),
    growth_pct      DECIMAL(5,2),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE skills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) UNIQUE NOT NULL,
    category    VARCHAR(50),
    is_ai_proof BOOLEAN DEFAULT false
);

CREATE TABLE profession_skills (
    profession_id UUID REFERENCES professions(id) ON DELETE CASCADE,
    skill_id      UUID REFERENCES skills(id) ON DELETE CASCADE,
    importance    INTEGER CHECK (importance BETWEEN 1 AND 5),
    is_at_risk    BOOLEAN DEFAULT false,
    PRIMARY KEY (profession_id, skill_id)
);

CREATE TABLE career_transitions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_profession     UUID REFERENCES professions(id),
    to_profession       UUID REFERENCES professions(id),
    match_score         INTEGER CHECK (match_score BETWEEN 0 AND 100),
    transition_reason   TEXT,
    avg_reskill_months  INTEGER,
    UNIQUE (from_profession, to_profession)
);

CREATE TABLE city_top_professions (
    city_id         UUID REFERENCES cities(id) ON DELETE CASCADE,
    profession_id   UUID REFERENCES professions(id) ON DELETE CASCADE,
    rank            INTEGER NOT NULL CHECK (rank > 0),
    local_demand    INTEGER,
    snapshot_date   DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (city_id, profession_id, snapshot_date)
);
