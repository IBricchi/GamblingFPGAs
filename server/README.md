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
