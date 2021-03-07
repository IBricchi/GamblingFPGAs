#!/bin/bash

# download js dependencies
mkdir -p js/p5
curl https://cdn.jsdelivr.net/npm/p5@1.2.0/lib/p5.js -o js/p5/p5.js
curl https://cdn.jsdelivr.net/npm/p5@1.2.0/lib/p5.min.js -o js/p5/p5.min.js
echo "*" > js/p5/.gitignore