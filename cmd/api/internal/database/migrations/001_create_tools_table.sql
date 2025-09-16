-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tools table with UUID primary key
CREATE TABLE IF NOT EXISTS tools (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL,
	status VARCHAR(50) NOT NULL DEFAULT 'IN_OFFICE',
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tools_status ON tools(status);

CREATE INDEX IF NOT EXISTS idx_tools_created_at ON tools(created_at DESC);

-- Create a trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$ 
	BEGIN NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_tools_updated_at
	BEFORE UPDATE ON tools 
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();