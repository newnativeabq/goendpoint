# Endpoint

Simple configurable endpoint to collect sensor data and photos on the home network and toss them into a cassandra db somehow.


## Requirements

> go >= 1.15

> Docker or Podman

## Docker

Setting up Cassandra.

1. Start cassandra and forward the port to the network

> podman run -d --network host  --name cassandra_db cassandra:latest 

> podman run -it --network host --rm cassandra cqlsh casandra_db

**Cassandra** did not work on Fedora 32.  The CQLSH container was not creating properly.  Trying postgresql for now.

Setting up PostgreSQL
1. Start postgres container

> podman run --name postgres_db -p 5432:5432 -e POSTGRES_PASSWORD=[dbpassword] -d postgres

2. Enter the container and create a database

> podman exec -it postgres_db psql -U postgres

> In the db shell: create database sensor_db


### Build

Building the program in go

> go build -o SensorEndpoint *.go