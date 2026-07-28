CREATE UNIQUE INDEX IF NOT EXISTS idx_competitors_project_name_unique ON competitors(project_id, name);
