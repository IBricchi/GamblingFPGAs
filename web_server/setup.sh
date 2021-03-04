#!/bin/bash

# update
apt-get update
apt-get -y install --no-install-recommends python3 python3-pip python3-venv curl

# setup python
python3 -m venv env
./env/bin/pip3 install -r requirements.txt

# download js dependencies
mkdir -p static/p5
curl -o static/p5/p5.js https://cdn.jsdelivr.net/npm/p5@1.2.0/lib/p5.js
curl -o static/p5/p5.min.js https://cdn.jsdelivr.net/npm/p5@1.2.0/lib/p5.min.js

# gitignore created setups
echo "*" > env/.gitignore
mkdir -p __pycache__
echo "*" > __pycache__/.gitignore
echo "*" > static/p5/.gitignore