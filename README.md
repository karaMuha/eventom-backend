# Eventom
## About this app
With eventom you can create and register for events
## How to run the app
**Prerequisites:**
- Docker installed (I am using docker desktop version 4.3 on Windows 11)

**Installation steps:**
- Clone the code with `git clone https://github.com/karaMuha/eventom-backend.git`
- make sure docker is running
- Inside the root directory of the project run the command `make setup` (this will generate a private key as .pem which is used to sign and verify jwt and create the folder db-data/postgres to persist data from the postgres container)
- run `make start`

The entry point of this app is `main.go`. On start up the app will try to connect to the postgres container `dbServer.go in package server`. Since postgres might need some time to be ready to accept requests, this app will try to establish a connection in an interval of 5 seconds for 10 times at max and crash if a connection to postgres cannot be established. After a connection to postgres has been established successfully, the http server will be initialized `httpServer.go in package server`. The http server initializes the logic layers (repositories, services, controllers and middleware) and the routes. Then the server starts and listens on the specified port (see `docker-compose.yaml`)

## ToDos
- Finish registration cancellation logic
- Provide tests for events and registrations logic
- implement purchasable events (using [kara-bank](https://github.com/karaMuha/kara-bank) for payment)
- implement RBAC