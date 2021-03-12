class GetData {
    requestSettings = true;
    username = "player1";
    password = "player1";
    constructor() {
        fetch("./setup.json")
            .then(response => response.json())
            .then(data => {
                this.requestSettings = data;
                this.requestSettings.headers.Authorization = 'Basic ' + btoa(this.username + ":" + this.password);
                // this.requestSettings.mode = "no-cors";
            })
            .catch(err => {
                this.requestSettings = false;
                console.error(err);
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
            return fetch(this.requestSettings.url, this.requestSettings)
                .then(request => request.json())
                .then(data => { return data })
                .catch(err => {
                    console.warn(err);
                    console.warn("Unable to fetch data from ", this.requestSettings.url);
                    return {}
                })
        }
    }
}