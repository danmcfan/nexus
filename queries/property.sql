-- name: CreateProperty :one
INSERT INTO property (pk_property_id, name, address, is_demo, fk_point_of_contact_id, fk_manager_id, fk_client_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListProperties :many
SELECT property.*,
  point_of_contact.pk_user_id AS point_of_contact_id,
  point_of_contact.first_name AS point_of_contact_first_name,
  point_of_contact.last_name AS point_of_contact_last_name,
  manager.pk_user_id AS manager_id,
  manager.first_name AS manager_first_name,
  manager.last_name AS manager_last_name,
  client.pk_client_id AS client_id,
  client.name AS client_name
FROM property
LEFT JOIN user AS point_of_contact ON property.fk_point_of_contact_id = point_of_contact.pk_user_id
LEFT JOIN user AS manager ON property.fk_manager_id = manager.pk_user_id
LEFT JOIN client ON property.fk_client_id = client.pk_client_id
ORDER BY client_name, property.pk_property_id;

-- name: ListPropertiesWithFilter :many
SELECT property.*,
  point_of_contact.pk_user_id AS point_of_contact_id,
  point_of_contact.first_name AS point_of_contact_first_name,
  point_of_contact.last_name AS point_of_contact_last_name,
  manager.pk_user_id AS manager_id,
  manager.first_name AS manager_first_name,
  manager.last_name AS manager_last_name,
  client.pk_client_id AS client_id,
  client.name AS client_name
FROM property
LEFT JOIN user AS point_of_contact ON property.fk_point_of_contact_id = point_of_contact.pk_user_id
LEFT JOIN user AS manager ON property.fk_manager_id = manager.pk_user_id
LEFT JOIN client ON property.fk_client_id = client.pk_client_id
WHERE 
  property.pk_property_id = ?
  OR property.name LIKE ?
  OR property.address LIKE ?
ORDER BY client_name, property.pk_property_id;

-- name: GetProperty :one
SELECT property.*,
  point_of_contact.pk_user_id AS point_of_contact_id,
  point_of_contact.first_name AS point_of_contact_first_name,
  point_of_contact.last_name AS point_of_contact_last_name,
  manager.pk_user_id AS manager_id,
  manager.first_name AS manager_first_name,
  manager.last_name AS manager_last_name,
  client.pk_client_id AS client_id,
  client.name AS client_name
FROM property
LEFT JOIN user AS point_of_contact ON property.fk_point_of_contact_id = point_of_contact.pk_user_id
LEFT JOIN user AS manager ON property.fk_manager_id = manager.pk_user_id
LEFT JOIN client ON property.fk_client_id = client.pk_client_id
WHERE pk_property_id = ?;

-- name: UpdateProperty :one
UPDATE property
SET name = ?, address = ?, is_demo = ?, fk_point_of_contact_id = ?, fk_manager_id = ?, fk_client_id = ?
WHERE pk_property_id = ?
RETURNING *;

-- name: DeleteProperty :exec
DELETE FROM property
WHERE pk_property_id = ?;