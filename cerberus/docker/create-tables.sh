#!/bin/bash
set -e
if ! psql -U $POSTGRES_USER -d $POSTGRES_DB -Atc "SELECT 1 FROM pg_database WHERE datname = 'cerberus'" | grep -qw 1; then
	createdb -U $POSTGRES_USER cerberus --owner=$POSTGRES_USER
fi

psql -v ON_ERROR_STOP=1 -U $POSTGRES_USER -d cerberus <<-EOSQL
	create or replace function unix_timestamp () returns bigint as 'select extract(epoch from now())' language sql;
	create table if not exists "permissions" (
	        "id" uuid primary key default gen_random_uuid (),
	        "name" varchar(30) not null,
	        "description" varchar(255) null,
	        "permissions" jsonb not null,
	        "disabled" boolean not null default false,
	        "created_at" bigint not null default unix_timestamp (),
	        "updated_at" bigint not null default unix_timestamp (),
	        "deleted_at" bigint null
	    );
	create unique index if not exists "idx_permissions_name" on "permissions" (lower("name")) where "deleted_at" is null;

	create table if not exists "domains" (
	        "id" uuid primary key default gen_random_uuid (),
	        "name" varchar(30) not null,
	        "description" varchar(255) null,
	        "paths" jsonb not null,
	        "disabled" boolean not null default false,
	        "created_at" bigint not null default unix_timestamp (),
	        "updated_at" bigint not null default unix_timestamp (),
	        "deleted_at" bigint null
	    );
	create unique index if not exists "idx_domains_name" on "domains" (lower("name")) where "deleted_at" is null;

	create table if not exists "roles" (
	        "id" uuid primary key default gen_random_uuid (),
	        "domain_id" uuid not null,
	        "permission_id" uuid not null,
	        "name" varchar(30) not null,
	        "description" varchar(255) null,
	        "disabled" boolean not null default false,
	        "created_at" bigint not null default unix_timestamp (),
	        "updated_at" bigint not null default unix_timestamp (),
	        "deleted_at" bigint null,
	        foreign key (domain_id) references domains (id) on update cascade on delete cascade,
	        foreign key (permission_id) references permissions (id) on update cascade on delete cascade
	    );
	create unique index if not exists "idx_roles_name"  on "roles" (lower("name"), "domain_id", "permission_id") where "deleted_at" is null;
EOSQL
