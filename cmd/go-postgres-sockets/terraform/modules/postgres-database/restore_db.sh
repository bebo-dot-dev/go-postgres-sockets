#!/usr/bin/env bash

#
# A bash shell script that performs database restore using the /postgres/notifications_db_backup.sql script
#
# This shell script is called by local-exec in ./modules/postgres-database terraform module main.tf config and
# required variables are passed to this script as positional arguments
#

CLOUD_SQL_PROXY_BIN=$1;
CONNECTION_STRING=$2;
CLOUD_SQL_PROXY_PORT=$3;
POSTGRES_USERNAME=$4;
POSTGRES_PASSWORD=$5;

# get the working directory of this script to correctly workout the location of the sql backup file
THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
SQL_BACKUP_FILE=$THIS_SCRIPT_DIR/../../../../../postgres/notifications_db_backup.sql

# start cloud_sql_proxy as a background process
CLOUD_SQL_PROXY_CMD="$CLOUD_SQL_PROXY_BIN -instances=$CONNECTION_STRING=tcp:0.0.0.0:$CLOUD_SQL_PROXY_PORT"
echo $CLOUD_SQL_PROXY_CMD
$CLOUD_SQL_PROXY_CMD &

# wait 5 seconds for cloud_sql_proxy to start
sleep 5

# db restore with psql via cloud_sql_proxy
psql "host=127.0.0.1 port=$CLOUD_SQL_PROXY_PORT sslmode=disable dbname=$POSTGRES_USERNAME user=$POSTGRES_USERNAME password=$POSTGRES_PASSWORD" < $SQL_BACKUP_FILE

# kill cloud_sql_proxy
CSP_PID=$(pgrep cloud_sql_proxy)
kill $CSP_PID