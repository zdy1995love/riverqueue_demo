-- Initialize River Queue schema
-- This script is automatically run when the PostgreSQL container starts
-- if you use the pg.sh script

-- Create the riverqueue schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS riverqueue;

-- Note: River will automatically create its tables on first run
-- This includes river_job, river_leader, river_migration, etc.

-- You can verify the schema was created by connecting to the database:
-- docker exec -it riverqueue-postgres psql -U riverqueue -d river_demo
-- Then run: \dn

-- After River creates the tables, you can view them with:
-- \dt riverqueue.*
