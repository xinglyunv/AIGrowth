-- Workspace and team membership foundation.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS current_space_id UUID;

CREATE TABLE IF NOT EXISTS spaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(120) NOT NULL,
    slug VARCHAR(80) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS space_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (space_id, user_id),
    CHECK (role IN ('owner', 'admin', 'member', 'viewer'))
);

CREATE TABLE IF NOT EXISTS space_invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    invited_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '7 days'),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (space_id, email, status),
    CHECK (role IN ('admin', 'member', 'viewer'))
);

CREATE INDEX IF NOT EXISTS idx_spaces_owner ON spaces(owner_id);
CREATE INDEX IF NOT EXISTS idx_space_members_user ON space_members(user_id);
CREATE INDEX IF NOT EXISTS idx_space_members_space ON space_members(space_id);
CREATE INDEX IF NOT EXISTS idx_space_invitations_email ON space_invitations(email);

CREATE OR REPLACE FUNCTION create_personal_space_for_user()
RETURNS TRIGGER AS $$
DECLARE
    new_space_id UUID;
BEGIN
    INSERT INTO spaces (owner_id, name, slug)
    VALUES (NEW.id, COALESCE(NULLIF(NEW.company_name, ''), NEW.username || '''s Workspace'),
            LEFT('user-' || REPLACE(NEW.id::text, '-', ''), 80))
    RETURNING id INTO new_space_id;

    INSERT INTO space_members (space_id, user_id, role)
    VALUES (new_space_id, NEW.id, 'owner');

    UPDATE users SET current_space_id = new_space_id WHERE id = NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS users_create_personal_space ON users;
CREATE TRIGGER users_create_personal_space
    AFTER INSERT ON users
    FOR EACH ROW EXECUTE FUNCTION create_personal_space_for_user();

INSERT INTO spaces (owner_id, name, slug)
SELECT u.id, COALESCE(NULLIF(u.company_name, ''), u.username || '''s Workspace'),
       LEFT('user-' || REPLACE(u.id::text, '-', ''), 80)
FROM users u
WHERE NOT EXISTS (SELECT 1 FROM spaces s WHERE s.owner_id = u.id);

INSERT INTO space_members (space_id, user_id, role)
SELECT s.id, s.owner_id, 'owner'
FROM spaces s
WHERE NOT EXISTS (
    SELECT 1 FROM space_members sm WHERE sm.space_id = s.id AND sm.user_id = s.owner_id
);

UPDATE users u
SET current_space_id = s.id
FROM spaces s
WHERE s.owner_id = u.id AND u.current_space_id IS NULL;

CREATE TRIGGER spaces_updated_at
    BEFORE UPDATE ON spaces
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER space_members_updated_at
    BEFORE UPDATE ON space_members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
