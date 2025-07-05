#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname decolecta_rucs <<-EOSQL
    CREATE TABLE IF NOT EXISTS sunat_rucs(
        ruc character varying(11) not null primary key,
        content character varying(300) not null
    );
    TRUNCATE TABLE sunat_rucs;
    ALTER TABLE sunat_rucs  SET UNLOGGED;
    \copy sunat_rucs FROM '/tmp/padron_reducido_ruc_utf8.csv' DELIMITER '|';
EOSQL
