# Eventom
## About this app
With eventom you can create and register for events
## How to run the app
**Prerequisites:**
- Docker installed (I am running docker desktop version 4.3 on Windows 11)

**Installation steps:**
- Clone the code with `git clone https://github.com/karaMuha/eventom-backend.git`
- Inside the root directory of the project create the directory `db-data`
- Inside `db-data` create the directory `postgres` (this directory will be used to persist the data from the postgres container)
- Generate a private key and save it as `id_rsa_priv.pem` inside the root directory (will be used to create and verify JWT)
- Run `docker-compose up -d`

The entry point of this app is `main.go`. On start up the app will try to connect to the postgres container `dbServer.go in package server`. Since postgres might need some time to be ready to accept requests, this app will try to establish a connection in an interval of 5 seconds for 10 times at max and crash if a connection to postgres cannot be established. After a connection to postgres has been established successfully, the http server will be initialized `httpServer.go in package server`. The http server initializes the logic layers (repositories, services, controllers and middleware) and the routes. Then the server starts and listens on the specified port (see `docker-compose.yaml`)

## ToDos
- Finish registrations logic
- Implement capacity feature for events/registrations
- Provide tests for events and registrations logic
- Figure out why docker occasionally failes to create reaper for testcontainer and how to fix / work around