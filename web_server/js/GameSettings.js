class GameSettings {
    getData;
    intervalGet;

    constructor(getData) {
        this.getData = getData;
    }

    start() {
        let t = this;
        this.intervalGet = setInterval(() => { t.update() }, 1000);
    }

    update() {
        this.getData.get()
            .then(data => this.processData(data))
            .catch(_ => console.warn("GameSettings: Unable to obtain new data."));
    }

    processData(data) {
        console.log(data);
    }

}