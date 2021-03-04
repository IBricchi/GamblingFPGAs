#!/bin/bash

apt-get update
apt-get -y install --no-install-recommends python3 python3-pip python3-venv make
python3 -m venv env
./env/bin/pip3 install -r requirements.txt
mkdir -p __pycache__
echo "*" > __pycache__/.gitignore
echo "*" > env/.gitignore