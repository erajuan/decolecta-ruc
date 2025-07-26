#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname sunat_collection <<-EOSQL
 -- Drop table if exists for idempotency
    DROP TABLE IF EXISTS sunat_ruc_extenses;

    -- Create as UNLOGGED for faster bulk load
    CREATE UNLOGGED TABLE sunat_ruc_extenses(
        ruc character varying(11) not null primary key,
        content character varying(512) not null
    );
    \copy sunat_ruc_extenses FROM '/tmp/ruc2utf8.csv' DELIMITER '|';
EOSQL
