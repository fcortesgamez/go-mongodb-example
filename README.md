# go-mongodb-example

**Note:** This is all work in progress...

Example of a service written in Golang with a Mongo DB and Gorilla Mux

## Checkout the project

```shell
git clone git@github.com:fcortesgamez/go-mongodb-example.git $GOPATH/src/github.com/fcortesgamez/go-mongodb-example
```

## Setup

Recommended setup steps:

* Install Docker
* Create a docker machine for the project (optional)
* Run MongoDB as a Docker container

### Create Docker machine

```shell
docker-machine create --driver virtualbox go-mongodb-example-machine
eval $(docker-machine env go-mongodb-example-machine)
```

### Run MongoDB as a Docker containers

docker run --name mongodb -d -p 27017:27017 mongo:latest

### Build the webshop web service

make clean build
