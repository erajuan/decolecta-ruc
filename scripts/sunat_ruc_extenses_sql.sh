#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname peru_collection <<-EOSQL
    create table if not exists sunat_ruc_extenses(
        ruc character varying(11) not null primary key,
        content character varying(512) not null
    );
    TRUNCATE TABLE sunat_ruc_extenses;
    ALTER TABLE sunat_ruc_extenses  SET UNLOGGED;
    \copy sunat_ruc_extenses FROM '/tmp/ruc2utf8.csv' DELIMITER '|';
EOSQL