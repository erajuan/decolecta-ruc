#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname decolecta_rucs <<-EOSQL
    -- Drop table if exists for idempotency
    DROP TABLE IF EXISTS sunat_rucs;

    -- Create as UNLOGGED for faster bulk load
    CREATE UNLOGGED TABLE sunat_rucs(
        ruc character varying(11) not null primary key,
        content character varying(300) not null
    );

    \copy sunat_rucs FROM '/tmp/padron_reducido_ruc_utf8.csv' DELIMITER '|';
EOSQL
