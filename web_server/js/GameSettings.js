// let q;
class GameSettings {
    refreshRate = 1000;

    comm = {};
    intervalGet;

    me = {};
    players = [];

    constructor(comm) {
        this.comm = comm;
    }

    start() {
        let t = this;
        t.updateOpen();
        this.intervalGet = setInterval(() => { t.updateOpen() }, t.refreshRate);
    }

    updateOpen() {
        this.comm.getOpen()
            .then(data => this.processOpen(data))
            .catch(err => {
                console.warn(err);
                console.warn("GameSettings: Unable to obtain new data.")
            });
    }

    processOpen(data) {
        // data.active = true;
        if (form.name != "OPEN" && !data.open) {
            removeForm();
            createOpenForm();
            removeTerminate();
        }
        if (form.name != "JOIN" && data.open && (data.players == null || !data.players.some(element => element.name == comm.username))) {
            removeForm();
            createJoinForm();
            createTerminate();
        }
        if (form.name != "START" && data.players != null && data.players.some(element => element.name == comm.username)) {
            removeForm();
            createStartForm();
            createTerminate();
        }
        if (data.active) {
            if (form != {})
                removeForm();
            canvas.show();
            loop();
            clearInterval(this.intervalGet);
            let t = this;
            t.updateActive();
            this.intervalGet = setInterval(() => { t.updateActive() }, t.refreshRate);
            createTerminate();
        }
    }

    updateActive() {
        this.comm.getActive()
            .then(data => this.processActive(data))
            .catch(err => {
                console.warn(err);
                console.warn("GameSettings: Unable to obtain new data.")
            });
    }

    processActive(data) {
        if (!data.active) {
            canvas.hide();
            noLoop();
            clearInterval(this.intervalGet);
            let t = this;
            t.updateOpen();
            this.intervalGet = setInterval(() => { t.updateOpen() }, t.refreshRate);
            createTerminate();
        }
        this.communityCards = data.communityCards;
        this.currentPlayer = data.currentPlayer;
        this.currentPot = 0;
        this.me = {};
        this.players = [];
        for (let i = 0; i < data.players.length; i++) {
            let player = data.players[i];
            this.currentPot += player.totalMoneyBetAmount;
            player.order = i;
            if (player.name == comm.username) {
                this.me = player;
            }
            this.players.push(player);
        }
        this.currentRound = data.currentRound;
    }

}