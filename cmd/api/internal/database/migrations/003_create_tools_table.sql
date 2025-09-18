-- Create tools table with UUID primary key
CREATE TABLE IF NOT EXISTS tools (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name TEXT NOT NULL,
	status tool_status NOT NULL DEFAULT 'IN_OFFICE',
	current_user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
	last_checked_out_at TIMESTAMP WITH TIME ZONE NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tools_status ON tools(status);

CREATE INDEX IF NOT EXISTS idx_tools_created_at ON tools(created_at DESC);

DROP TRIGGER IF EXISTS update_tools_updated_at ON tools;
CREATE TRIGGER update_tools_updated_at
	BEFORE UPDATE ON tools 
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();

-- function for correctly setting the last_checked_out_time on insert or update
CREATE OR REPLACE FUNCTION set_last_checked_out_at()
RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.current_user_id IS NOT NULL AND NEW.last_checked_out_at IS NULL THEN
            NEW.last_checked_out_at := NOW();
        END IF;
    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW.current_user_id IS NOT NULL AND OLD.current_user_id IS NULL THEN
            NEW.last_checked_out_at := NOW();
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_last_checked_out_at ON tools;
CREATE TRIGGER update_last_checked_out_at
    BEFORE INSERT OR UPDATE ON tools
    FOR EACH ROW
    EXECUTE FUNCTION set_last_checked_out_at();
