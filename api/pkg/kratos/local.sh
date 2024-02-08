#!/bin/sh

# starting up all services including kratos
docker-compose --profile kratos up -d

# function to check if the docker container is running
container_is_running() {
  docker inspect -f '{{.State.Running}}' "$1" 2>/dev/null
}

# function to wait for a docker container to start
wait_for_container() {
  container_name=$1
  while ! container_is_running $container_name; do
    sleep 1
  done
  # extra delay to wait, otherwise gettinng error while connecting to database while migrating
  sleep 2
  echo "===== Container $container_name is running"
}

wait_for_container "golang-api-postgresdb"

echo "===== Running the migrations..."
# running the migrations
make -C ../../. migrate-up

echo "===== Starting golang api server..."
# starting the golang backend server
make -C ../../. start-api
