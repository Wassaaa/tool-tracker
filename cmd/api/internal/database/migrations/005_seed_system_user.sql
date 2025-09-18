-- Seed a system user used for automated actions and default actor attribution
-- Use a fixed UUID so application code can reference it safely.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM users WHERE id = '00000000-0000-0000-0000-000000000001'
    ) THEN
        INSERT INTO users (id, name, email, role)
        VALUES ('00000000-0000-0000-0000-000000000001', 'System', 'system@local', 'ADMIN');
    END IF;
END$$;
