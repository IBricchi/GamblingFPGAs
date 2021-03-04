import json
import time
import requests
import threading

from flask import Flask, render_template
app = Flask(__name__)


data = ""

# setup web ports
@app.route('/')
def index():
    return "hello world"

@app.route('/j/')
def data_src():
    global data
    return data

# setup thread
def curl_data():
    time.sleep(10)

    global data
    headers = {
        'Accept': 'application/json',
    }
    
    while True:
        response = requests.get('https://random-data-api.com/api/beer/random_beer', headers=headers)
        data = response.content.decode('utf-8')
        time.sleep(1)
    
if __name__ == '__main__':
    # run trhead
    t = threading.Thread(target=curl_data, daemon=True)
    t.start()

    # run flask
    app.run()