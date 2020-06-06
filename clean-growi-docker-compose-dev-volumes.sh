#!/bin/sh

if (( $# != 1 )); then
  echo "Error: specify only one argument" 1>&2
  exit 1
fi

if [[ `basename $1` != "growi-docker-compose-dev" ]]; then
  echo "Error: specify a valid path to growi-docker-compose-dev directory" 1>&2
  exit 1
fi

cd $1
docker-compose down
docker volume rm growi-docker-compose-dev_es_data
docker volume rm growi-docker-compose-dev_es_plugins
docker volume rm growi-docker-compose-dev_growi_data
docker volume rm growi-docker-compose-dev_mariadb_data
docker volume rm growi-docker-compose-dev_mongo_configdb
docker volume rm growi-docker-compose-dev_mongo_db
docker volume rm growi-docker-compose-dev_sqlite_db
