# Eventom
## About this app
With eventom you can create and register for events. The purpose of this project was to get familiar with golang and its tooling for api development.
## How to run the app
**Prerequisites:**
- Docker installed (I am using docker desktop version 4.3 on Windows 11)

**Installation steps:**
- Clone the code with `git clone https://github.com/karaMuha/eventom-backend.git`
- make sure docker is running
- Inside the root directory of the project run the command `make setup` (this will generate a private key as .pem which is used to sign and verify jwt and create the folder db-data/postgres to persist data from the postgres container)
- run `make start`

The entry point of this app is `main.go`. On start up the app will try to connect to the postgres container `dbServer.go in package server`. Since postgres might need some time to be ready to accept requests, this app will try to establish a connection in an interval of 5 seconds for 10 times at max and crash if a connection to postgres cannot be established. After a connection to postgres has been established successfully, the http server will be initialized `httpServer.go in package server`. The http server initializes the logic layers (repositories, services, controllers and middleware) and the routes. Then the server starts and listens on the specified port (see `docker-compose.yaml`)

## Usage
- POST /signup -> signup as a user with your email and a password
```
{
    "email": "test@test.com",
    "password": "test123"
}
```
- POST /login -> login with you signed up user to get a jwt
```
{
    "email": "test@test.com",
    "password": "test123"
}
```
- (protected) POST /events -> create an event with an event name, location, date, and max capacity
```
{
    "name": "Test",
    "location": "Köln",
    "date": "1994-10-27T21:00:00Z",
    "max_capacity": 3
}
```
- GET /events/{id} -> get event with given event id
- GET /events?page={number>=1}&page_size=[10, 15, 20, 25] -> list all events. You can search, filter and sort results using query parameters
  - name -> provide parts of the event name to search for it
  - location -> filter for event location
  - capacity -> filter for minimum free capacity
  - column -> sort by column [id, event_name, event_description, event_date, max_capacity, amount_registrations]
  - order -> set order [DESC, ASC]
  - e.g. /events?page=1&page_size=10&location=Köln&capacity=4&sort=event_date&order=DESC
- (protected) PUT /events/{id} -> update event with given event id. User can only update events created by himself
- (protected) DELETE /events/{id} -> delete event with given event id. User can only delete events created by himself

- (protected) POST /registrations -> register for an event. Provide event id in request body, user id will be extraced from jwt
```
{
    "event_id": {id}
}
```
- GET /registrations -> list all registration (will be refactored to list all registrations of logged in user)
- (protected) DELETE /registrations/{id} -> cancel registration with given registration id. User can only cancel his own registrations

## ToDos
- finish registration cancellation logic
- cancel registrations when event is deleted
- provide tests for events and registrations logic
- provide tests for transaction handler
- implement purchasable events (using [kara-bank](https://github.com/karaMuha/kara-bank) for payment and expand filtering capability to include price filtering)
- implement RBAC