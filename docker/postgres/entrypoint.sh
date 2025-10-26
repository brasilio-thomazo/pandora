#!/bin/bash
set -e

hba_file="${PGDATA}/pg_hba.conf"
pg_ident_file="${PGDATA}/pg_ident.conf"
config_file="${PGDATA}/postgresql.conf"

if [ ! -d "$PGDATA" ]; then
	mkdir -p "$PGDATA"
	chown -R postgres:postgres "$PGDATA"
fi

POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-postgres}"
POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_DB="${POSTGRES_DB:-postgres}"
POSTGRES_REPLICATION_USER="${POSTGRES_REPLICATION_USER:-replicator}"
POSTGRES_REPLICATION_PASSWORD="${POSTGRES_REPLICATION_PASSWORD:-replicator}"
POSTGRES_SLOT_NAME="${POSTGRES_SLOT_NAME:-replication_slot}"
POSTGRES_REPLICATION_MODE="${POSTGRES_REPLICATION_MODE:-master}"
POSTGRES_PRIMARY_PORT="${POSTGRES_PRIMARY_PORT:-5432}"

init_db() {
	if [ -f "$PGDATA/PG_VERSION" ]; then
		echo "database already initialized"
		return
	fi
	if [ "$POSTGRES_REPLICATION_MODE" = "master" ]; then
		echo "initializing database"
		gosu postgres bash -c "initdb --auth=trust --auth-host=scram-sha-256 --auth-local=trust -D $PGDATA --pwfile=<(echo -n '$POSTGRES_PASSWORD')"
		adjust_configs
	fi
}

adjust_configs() {
	sed -i "s/#wal_level = replica/wal_level = replica/" $config_file
	sed -i "s/#max_wal_senders = 10/max_wal_senders = 10/" $config_file
	sed -i "s/#wal_keep_segments = 10/wal_keep_segments = 10/" $config_file
	sed -i "s/#max_replication_slots = 10/max_replication_slots = 10/" $config_file
	sed -i "s/#hot_standby_feedback = off/hot_standby_feedback = on/" $config_file
	sed -i "s/#hot_standby = on/hot_standby = on/" $config_file

	if ! grep -Fxq "host all all 0.0.0.0/0" $hba_file; then
		echo "enabling scram-sha-256 authentication for remote connections"
		echo "host all all 0.0.0.0/0 scram-sha-256" | tee -a $hba_file
	fi

	if ! grep -Fxq "host replication all 0.0.0.0/0" $hba_file; then
		echo "enabling scram-sha-256 authentication for replication connections"
		echo "host replication all 0.0.0.0/0 scram-sha-256" | tee -a $hba_file
	fi

	if [ ! -f $pg_ident_file ]; then
		touch $pg_ident_file
		chown postgres:postgres $pg_ident_file
	fi
}

waiting_for_primary() {
	until pg_isready -h "$POSTGRES_PRIMARY_HOST" -p "$POSTGRES_PRIMARY_PORT" -U "$POSTGRES_REPLICATION_USER" -d postgres; do
		echo "waiting for primary to be ready"
		sleep 3
	done
}

backup_db() {
	if [ -f "$PGDATA/PG_VERSION" ]; then
		echo "database already initialized"
		return
	fi
	waiting_for_primary
	echo "initializing database from primary"
	pg_basebackup -h "$POSTGRES_PRIMARY_HOST" -p "$POSTGRES_PRIMARY_PORT" -U "$POSTGRES_REPLICATION_USER" -D "$PGDATA" -v
	chmod 700 "$PGDATA"
	chown -R postgres:postgres "$PGDATA"
}

temp_server() {
	echo "starting temp server"
	gosu postgres bash -c "pg_ctl -D $PGDATA -o '-c listen_addresses=\"localhost\"' start"
	if [ $? != 0 ]; then
		echo "failed to start temp server"
		exit 1
	fi
}

create_postgres_user() {
	psql -U postgres -d postgres -Atc "SELECT 1 FROM pg_roles WHERE rolname = '$POSTGRES_USER'" | grep -qw 1
	if [ $? = 0 ]; then
		echo "postgres user already exists, updating password"
		psql -U postgres -c "ALTER USER $POSTGRES_USER WITH SUPERUSER ENCRYPTED PASSWORD '$POSTGRES_PASSWORD'"
	elif [ -n "$POSTGRES_USER" ] && [ -n "$POSTGRES_PASSWORD" ]; then
		echo "creating postgres user"
		psql -U postgres -c "CREATE USER $POSTGRES_USER WITH SUPERUSER ENCRYPTED PASSWORD '$POSTGRES_PASSWORD'"
	fi
}

create_replication_user() {
	if [ -n "$POSTGRES_REPLICATION_USER" ] && [ -n "$POSTGRES_REPLICATION_PASSWORD" ]; then
		echo "creating or updating replication user"
		if psql -U postgres -d postgres -Atc "SELECT 1 FROM pg_roles WHERE rolname = '$POSTGRES_REPLICATION_USER'" | grep -qw 1; then
			echo "replication user already exists, updating password"
			psql -U postgres -c "ALTER ROLE $POSTGRES_REPLICATION_USER WITH REPLICATION LOGIN ENCRYPTED PASSWORD '$POSTGRES_REPLICATION_PASSWORD'"
		else
			echo "creating replication user"
			psql -U postgres -c "CREATE USER $POSTGRES_REPLICATION_USER WITH REPLICATION LOGIN ENCRYPTED PASSWORD '$POSTGRES_REPLICATION_PASSWORD'"
		fi
	fi
}

create_slot() {
	if [ -n "$POSTGRES_SLOT_NAME" ]; then
		echo "checking for existing slot $POSTGRES_SLOT_NAME"
		if ! psql -U postgres -d postgres -Atc "SELECT 1 FROM pg_replication_slots WHERE slot_name = '$POSTGRES_SLOT_NAME'" | grep -qw 1; then
			echo "slot does not exists creating slot $POSTGRES_SLOT_NAME"
			psql -U postgres -c "SELECT pg_create_physical_replication_slot('$POSTGRES_SLOT_NAME')"
		fi
	fi
}

execute_init_scripts() {
	for f in /docker-entrypoint-initdb.d/*; do
		if [ ! -f "$f" ]; then
			continue
		fi
		echo "$0: processing $f"
		case "$f" in
		*.sh)
			echo "$0: running $f"
			. "$f"
			;;
		*.sql)
			echo "$0: running $f"
			psql -U postgres -f "$f"
			;;
		*.sql.gz)
			echo "$0: running $f"
			gunzip -c "$f" | psql -U postgres
			;;
		*) echo "$0: ignoring $f" ;;
		esac
	done
}

make_password_file() {
	if [ -z "$POSTGRES_REPLICATION_USER" ] || [ -z "$POSTGRES_REPLICATION_PASSWORD" ]; then
		echo "POSTGRES_REPLICATION_USER and POSTGRES_REPLICATION_PASSWORD must be set when POSTGRES_REPLICATION_MODE is set to slave"
		exit 1
	fi
	echo "creating password file"
	echo "${POSTGRES_PRIMARY_HOST}:${POSTGRES_PRIMARY_PORT}:replication:${POSTGRES_REPLICATION_USER}:${POSTGRES_REPLICATION_PASSWORD}" >"$HOME/.pgpass"
	chmod 0600 "$HOME/.pgpass"
}

init_replica() {
	if [ "$POSTGRES_REPLICATION_MODE" = "slave" ]; then
		if [ -z "$POSTGRES_PRIMARY_HOST" ] || [ -z "$POSTGRES_PRIMARY_PORT" ]; then
			echo "POSTGRES_PRIMARY_HOST and POSTGRES_PRIMARY_PORT must be set when POSTGRES_REPLICATION_MODE is set to slave"
			exit 1
		fi
		make_password_file
		backup_db
		adjust_configs
	fi
}

if [ "$POSTGRES_REPLICATION_MODE" = "slave" ]; then
	init_replica
	waiting_for_primary
fi

if [ "$POSTGRES_REPLICATION_MODE" = "master" ]; then
	init_db
	adjust_configs
	temp_server
	create_postgres_user
	create_replication_user
	create_slot
	execute_init_scripts
	echo "stoping temp server"
	gosu postgres bash -c "pg_ctl -D $PGDATA stop"
fi

echo "setup complete, starting postgres server"
exec gosu postgres "$@"
