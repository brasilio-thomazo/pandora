#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 -U $POSTGRES_USER -d cerberus <<-EOSQL
	insert into "permissions" ("name", "description", "permissions")
	values
		('tudo', 'pode realizar todas as operações', '["READ", "WRITE", "UPDATE", "DELETE"]'),
		('criar e editar', 'pode criar e editar registros', '["READ", "WRITE", "UPDATE"]'),
		('criar e deletar', 'pode criar e deletar registros', '["READ", "WRITE", "DELETE"]'),
		('criar', 'pode criar registros', '["READ", "WRITE"]'),
		('editar e deletar', 'pode editar e deletar registros', '["READ", "UPDATE", "DELETE"]'),
		('editar', 'pode editar registros', '["READ", "UPDATE"]'),
		('deletar', 'pode deletar registros', '["READ", "DELETE"]'),
		('visualizar', 'pode visualizar registros', '["READ"]') on conflict do nothing;

	insert into "domains" ("name", "description", "paths")
	values
		('ALL', 'acesso a todos os domínios', '["/"]'),
		('AUTH', 'acesso ao domínio de autenticação', '["/permission", "/role", "/group", "/user"]') on conflict do nothing;
EOSQL
