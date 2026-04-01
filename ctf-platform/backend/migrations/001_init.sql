CREATE TYPE user_role AS ENUM ('player', 'admin');

CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(64)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    role          user_role    NOT NULL DEFAULT 'player',
    is_disabled   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS teams (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(128) NOT NULL UNIQUE,
    invite_code VARCHAR(64)  NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_members (
    user_id INT NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    team_id INT NOT NULL REFERENCES teams(id)  ON DELETE CASCADE,
    PRIMARY KEY (user_id, team_id)
);

CREATE TABLE IF NOT EXISTS challenges (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    category    VARCHAR(64)  NOT NULL,
    points      INT          NOT NULL DEFAULT 0,
    flag_hash   TEXT         NOT NULL,
    is_visible  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS submissions (
    id           SERIAL PRIMARY KEY,
    user_id      INT         NOT NULL REFERENCES users(id)       ON DELETE CASCADE,
    challenge_id INT         NOT NULL REFERENCES challenges(id)  ON DELETE CASCADE,
    is_correct   BOOLEAN     NOT NULL DEFAULT FALSE,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_submissions_user        ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_challenge   ON submissions(challenge_id);
CREATE INDEX IF NOT EXISTS idx_submissions_correct     ON submissions(user_id, challenge_id) WHERE is_correct = TRUE;
