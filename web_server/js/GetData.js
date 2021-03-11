class GetData {
    requestSettings = true;
    constructor() {
        fetch("./setup.json")
            .then(response => response.json())
            .then(data => this.requestSettings = data)
            .catch(_ => {
                this.requestSettings = false;
                console.error("Unable to fetch setup.json file.");
            });
    }

    async get() {
        if (this.requestSettings === false) {
            console.error("GetData is in error state.");
            return {};
        } else if (this.requestSettings === true) {
            console.log("GetData is still fetching settings.");
            return {}
        } else {
            return fetch(this.requestSettings.data_url, {
                    // body: setup.data,
                    headers: this.requestSettings.headers,
                    method: this.requestSettings.request
                })
                .then(request => request.json())
                .then(data => { return data })
                .catch(_ => {
                    console.warn("Unable to fetch data from ", this.requestSettings.data_url);
                    return {}
                })
        }
    }
}