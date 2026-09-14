-- name: Create{{ENTITY_PASCAL}} :exec
INSERT INTO {{MODULE_NAME}}.{{ENTITIES}} (
    {{ENTITY}}_uuid, name, created_at, updated_at
) VALUES ($1, $2, $3, $4);

-- name: Get{{ENTITY_PASCAL}}ByUUID :one
SELECT {{ENTITY}}_uuid, name, created_at, updated_at
FROM {{MODULE_NAME}}.{{ENTITIES}}
WHERE {{ENTITY}}_uuid = $1;
