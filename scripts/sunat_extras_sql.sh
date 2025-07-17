#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname peru_collection <<-EOSQL
    CREATE TABLE IF NOT EXISTS sunat_ruc_extras(
        ruc character varying(11) not null,
        type_id integer not null,
        position integer not null,
        content character varying(500),
        primary key (ruc, type_id, position)
    );
    TRUNCATE TABLE sunat_ruc_extras;
    ALTER TABLE sunat_ruc_extras  SET UNLOGGED;
    \copy sunat_ruc_extras FROM '/tmp/buenos_contribuyentes.csv' DELIMITER '|';
    \copy sunat_ruc_extras FROM '/tmp/agentes_retencion.csv' DELIMITER '|';
    \copy sunat_ruc_extras FROM '/tmp/padron_reducido_local_anexo.csv' DELIMITER '|';
EOSQL