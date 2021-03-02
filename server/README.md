## Server Setup

```bash
 ./setup.sh
```

Make sure to reboot the machine after executing the script

## Http Server

### Setup

Default port 3000:

```bash
  cd cmd/server
  go run main.go
```

Using custom port:

```bash
  cd cmd/server
  go run main.go -httpPort {your Port number}
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

Incorrect data format will return an error code.
