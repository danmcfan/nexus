CREATE TABLE metadata (
    fk_record_id TEXT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    fk_created_by_id TEXT NOT NULL,
    fk_updated_by_id TEXT
);