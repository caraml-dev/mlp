--  All secret storage types must be added here
-- 'internal' is temporarily added to handle existing secret stored in the database
CREATE TYPE secret_storage_type AS ENUM ('vault', 'internal');
CREATE TYPE secret_storage_scope AS ENUM ('project', 'global');

CREATE TABLE IF NOT EXISTS secret_storages (
    id                          SERIAL PRIMARY KEY,
    project_id                  integer REFERENCES projects(id),
    name                        VARCHAR(64) NOT NULL,
    type                        secret_storage_type NOT NULL,
    scope                       secret_storage_scope NOT NULL,
    config                      JSONB,
    created_at                  timestamp NOT NULL default current_timestamp,
    updated_at                  timestamp NOT NULL default current_timestamp,
    UNIQUE (project_id, name)
);

-- Create an 'internal' secret storage
INSERT INTO secret_storages (name, type, scope) VALUES ('internal', 'internal', 'global');

-- Add a foreign key constraint to secrets table
ALTER TABLE secrets ADD COLUMN secret_storage_id integer;
ALTER TABLE secrets ADD CONSTRAINT fk_secret_storage_id FOREIGN KEY (secret_storage_id) REFERENCES secret_storages(id) ON DELETE CASCADE;

-- Update existing secrets to use 'internal' secret storage during migration
UPDATE secrets 
SET secret_storage_id = (SELECT id FROM secret_storages WHERE name = 'internal') 
WHERE secret_storage_id IS NULL;
