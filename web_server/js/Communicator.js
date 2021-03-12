class Communicator {
    username = "";

    serverIP = "18.132.52.158:3000";
    modelRequest = {
        headers: {
            "Content-Type": "application/json; charset=UTF-8"
        },
        method: "GET"
    };

    checkCredUrl = this.serverIP + "/testCredentials";
    checkCredRequest = this.modelRequest;

    openGameUrl = this.serverIP + "/poker/openGameStatus";
    openGameRequest = this.modelRequest;

    activeGameUrl = "./testData.json";
    activeGameRequest = this.modelRequest;

    constructor(username, password) {
        this.username = username;
        let auth = 'Basic ' + btoa(username + ":" + password);
        this.checkCredRequest.headers.Authorization = auth;
        this.openGameRequest.headers.Authorization = auth;
        this.activeGameRequest.headers.Authorization = auth
    }

    async getCheckCred() {
        return fetch(this.checkCredUrl, this.checkCredRequest)
            .then(request => request.json())
            .then(data => { return data.valid })
            .catch(err => {
                console.warn(err);
                console.warn("Communicator: Unable to fetch data from ", this.checkCredUrl);
                return false;
            });
    }

    async getOpen() {
        return fetch(this.openGameUrl, this.openGameRequest)
            .then(request => request.json())
            .then(data => { return data })
            .catch(err => {
                console.warn(err);
                console.warn("Communicator: Unable to fetch data from ", this.openGameUrl);
                return false
            });
    }

    async getActive() {
        return fetch(this.activeGameUrl, this.activeGameRequest)
            .then(request => request.json())
            .then(data => { return data })
            .catch(err => {
                console.warn(err);
                console.warn("Communicator: Unable to fetch data from ", this.activeGameUrl);
                return false
            });
    }
}