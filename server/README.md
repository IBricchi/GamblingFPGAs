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

### Accessing static test data

Open `http://localhost:3000/test´ using your browser or using the command line:

```bash
  curl http://localhost:3000/test
```
