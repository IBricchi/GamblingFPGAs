## Server Setup

Initial setup for installing the required dependencies.

```bash
 ./setup.sh
```

Make sure to reboot the machine after executing the script

## Http Server

### Run Server

Default port is `3000`.
Default database name is `serverDB.db`.

```bash
  ./run_server.sh
```

Custom port and/or database name:
```bash
  ./run_server.sh {optional custom port} {optional custom database name}
```

### Create New User

Default database name is `serverDB.db`.
This must be the same database that is used for the server!

Every poker player must also be a user.

```bash
  ./run_credentials.sh
```

Custom database name:
```bash
  ./run_credentials.sh {optional custom database name}
```

### Accessing Static Test Data

#### Without password protection

Open `http://localhost:3000/public/test` using your browser or using the command line:

```bash
  curl http://localhost:3000/public/test
```

#### With password protection

Open `http://{username}:{password}@localhost:3000/test` using your browser or using the command line:

```bash
  curl http://{username}:{password}@localhost:3000/test
```

### Sending Test Data To Server

#### Without password protection

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  --data '{"info":"test data","data":[100,200,300,400,500]}' \
  http://localhost:3000/public/test/dynamic
```
The received data will be inserted into the server's sqlite3 database.

Incorrect data formats will return an error code.

#### With password protection

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  --data '{"info":"test data","data":[100,200,300,400,500]}' \
  http://{username}:{password}@localhost:3000/test/dynamic
```
The received data will be inserted into the server's sqlite3 database.

Incorrect data formats will return an error code.

## Poker

### Opening a game

This is the first phase of every game.

Players can join a game during the open phase.

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  --data '{"initialPlayerMoney":2000,"smallBlindValue":5}' \
  http://test:test@localhost:3000/poker/openGame
```

### Joining a game

This is only possible during the open phase.

A player must be a registered user to join (See "Create New User" section above).
The player name will correspond to the username specified in the URL (http basic auth).

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  http://player1:player1@localhost:3000/poker/joinGame
```

### Starting a game

A game can be started once at least two players have joined.
This will end the open phase and start the active phase.

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  http://test:test@localhost:3000/poker/startGame
```

### Terminating a game

All game state can be reset using a terminate request.
A terminate request can always be called.

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  http://test:test@localhost:3000/poker/terminateGame
```

### Getting game status/state

Use ```curl -i``` for additional information.

```bash
  curl http://test:test@localhost:3000/poker/openGameStatus
```

```bash
  curl -i http://test:test@localhost:3000/poker/activeGameStatus
```
