create table if not exists sunat_ruc_extras(
	ruc character varying(11) not null,
	type_id integer not null,
	position integer not null,
	content character varying(500),
	primary key (ruc, type_id, position)
);

create index sunat_ruc_extras_ruc_idx on sunat_ruc_extras(ruc);

create table sunat_ruc_extenses(
	ruc character varying(11) not null primary key,
	content character varying(512) not null
);
