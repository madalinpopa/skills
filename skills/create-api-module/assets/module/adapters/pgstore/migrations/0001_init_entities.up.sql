BEGIN;

CREATE SCHEMA IF NOT EXISTS {{MODULE_NAME}};

CREATE TABLE {{MODULE_NAME}}.{{ENTITIES}}
(
    {{ENTITY}}_uuid uuid NOT NULL,
    name varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ({{ENTITY}}_uuid)
);

COMMIT;
