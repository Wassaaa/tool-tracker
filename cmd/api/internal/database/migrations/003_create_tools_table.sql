-- Create tools table with UUID primary key
CREATE TABLE IF NOT EXISTS tools (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name TEXT NOT NULL,
	status tool_status NOT NULL DEFAULT 'IN_OFFICE',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tools_status ON tools(status);

CREATE INDEX IF NOT EXISTS idx_tools_created_at ON tools(created_at DESC);

CREATE TRIGGER update_tools_updated_at
	BEFORE UPDATE ON tools 
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();