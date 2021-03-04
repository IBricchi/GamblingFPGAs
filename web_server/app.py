import json
import time
import requests
import threading

from flask import Flask, render_template
app = Flask(__name__)

data = {}
setup = {}

# setup web ports
@app.route('/')
def index():
    return render_template("index.html")

@app.route('/data/')
def data_src():
    global data
    return data

# setup thread
def curl_data():
    time.sleep(10)

    global data
    global setup
    data_url = setup["data_url"]
    headers = setup["headers"]
    data = setup["data"]
    
    while True:
        response = requests.get(data_url, headers=headers, data=data)
        data = json.loads(response.content.decode('utf-8'))
        time.sleep(1)
    
if __name__ == '__main__':
    # get setup data
    with open("setup.json", "r") as f:
        setup = json.loads(f.read())

    # run trhead
    t = threading.Thread(target=curl_data, daemon=True)
    t.start()

    # run flask
    app.run(port = setup["port"])