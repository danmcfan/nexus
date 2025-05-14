CREATE TABLE property (
    pk_property_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    is_demo BOOLEAN NOT NULL DEFAULT FALSE,
    fk_point_of_contact_id TEXT,
    fk_manager_id TEXT,
    fk_client_id TEXT NOT NULL
);