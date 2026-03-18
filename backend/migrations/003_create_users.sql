-- ============================================================
-- FILE: migrations/003_create_users.sql
-- ============================================================

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(100),
    password_hash   VARCHAR(255),
    avatar_url      TEXT,
    auth_provider   VARCHAR(20) DEFAULT 'email',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE user_profiles (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE PRIMARY KEY,
    current_job_id  UUID REFERENCES professions(id),
    city_id         UUID REFERENCES cities(id),
    years_exp       INTEGER DEFAULT 0,
    education       VARCHAR(50),
    target_job_id   UUID REFERENCES professions(id),
    skills          TEXT[] DEFAULT '{}',
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE saved_careers (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    profession_id   UUID REFERENCES professions(id) ON DELETE CASCADE,
    notes           TEXT,
    saved_at        TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, profession_id)
);

CREATE TABLE ai_analyses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    profession_slug VARCHAR(100) NOT NULL,
    location        VARCHAR(100),
    request_hash    VARCHAR(64) UNIQUE NOT NULL,
    result          JSONB NOT NULL,
    model_used      VARCHAR(50),
    tokens_used     INTEGER,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_ai_analyses_hash   ON ai_analyses(request_hash);
CREATE INDEX idx_ai_analyses_slug   ON ai_analyses(profession_slug);
CREATE INDEX idx_ai_analyses_created ON ai_analyses(created_at DESC);
