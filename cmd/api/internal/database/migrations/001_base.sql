-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- Create a trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$ 
	BEGIN NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ language plpgsql;

-- Soft delete function
CREATE OR REPLACE FUNCTION soft_delete()
RETURNS TRIGGER AS $$
BEGIN
    NEW.deleted_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- pg_type enums
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tool_status') THEN
        CREATE TYPE tool_status AS ENUM ('IN_OFFICE','CHECKED_OUT','MAINTENANCE','LOST');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('EMPLOYEE','ADMIN','MANAGER');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
        CREATE TYPE event_type AS ENUM (
            'TOOL_CHECKED_OUT','TOOL_CHECKED_IN','TOOL_MARKED_LOST','TOOL_MARKED_MAINTENANCE','TOOL_MARKED_IN_OFFICE'
        );
    END IF;
END$$;