## Server Setup

```bash
 ./setup.sh
```

Make sure to reboot the machine after executing the script

## Http Server

### Setup

Default port is `3000`.
Default database name is `serverDB.db`.

```bash
  ./run_server.sh
```

Custom port and/or database name:
```bash
  ./run_server.sh {optional custom port} {optional custom database name}
```

### Accessing Static Test Data

Open `http://localhost:3000/test´ using your browser or using the command line:

```bash
  curl http://localhost:3000/test
```

### Sending Test Data To Server

```bash
  curl --header "Content-Type: application/json; charset=UTF-8" \
  --request POST \
  --data '{"info":"test data","data":[100,200,300,400,500]}' \
  http://localhost:3000/test/dynamic
```
The received data will be inserted into the server's sqlite3 database.

Incorrect data formats will return an error code.
