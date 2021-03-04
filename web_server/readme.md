## web-server setup

Initial setup for installing the required dependencies.

```bash
 ./setup.sh
```

## Run server

### Setup

To run the server simply run.

```bash
  ./run.sh
```

Confugurations for the server can be changed in setup.json file. This file contains setup data for port the servers should listen on, what headers to use in curl request, data to use in curl request, and what url to use for curl requrest

### Web Files

Main web page is hosted on:
Currently only displays json back to user
http://127.0.0.1:port/

Mirror of json file is hosted on:
http://127.0.0.1:port/data/

