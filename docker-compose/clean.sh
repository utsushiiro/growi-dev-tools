#!/bin/bash

docker-compose down
docker volume rm growi-dev-tools_es_data
docker volume rm growi-dev-tools_mysql_data
docker volume rm growi-dev-tools_mongo_configdb
docker volume rm growi-dev-tools_mongo_db
docker volume rm growi-dev-tools_sqlite_db
