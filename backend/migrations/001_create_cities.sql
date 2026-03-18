-- ============================================================
-- FILE: migrations/001_create_cities.sql
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE cities (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    country     VARCHAR(100) NOT NULL,
    region      VARCHAR(100) NOT NULL,
    lat         DECIMAL(9,6) NOT NULL,
    lng         DECIMAL(9,6) NOT NULL,
    population  INTEGER,
    timezone    VARCHAR(50),
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE city_scores (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id         UUID REFERENCES cities(id) ON DELETE CASCADE,
    score           INTEGER CHECK (score BETWEEN 0 AND 100),
    job_growth_pct  DECIMAL(5,2),
    remote_score    INTEGER CHECK (remote_score BETWEEN 0 AND 100),
    ai_investment   INTEGER CHECK (ai_investment BETWEEN 0 AND 100),
    talent_demand   INTEGER CHECK (talent_demand BETWEEN 0 AND 100),
    cost_of_living  INTEGER CHECK (cost_of_living BETWEEN 0 AND 100),
    snapshot_date   DATE NOT NULL DEFAULT CURRENT_DATE,
    source          VARCHAR(50) DEFAULT 'manual',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_city_scores_city_date ON city_scores(city_id, snapshot_date DESC);
CREATE UNIQUE INDEX idx_city_scores_city_snapshot ON city_scores(city_id, snapshot_date);
