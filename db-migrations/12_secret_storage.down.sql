ALTER TABLE secrets DROP COLUMN secret_storage_id;
DROP TABLE IF EXISTS secret_storages;
DROP TYPE IF EXISTS secret_storage_type;
DROP TYPE IF EXISTS secret_storage_scope;