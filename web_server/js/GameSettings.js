class GameSettings {
    comm = {};
    intervalGet;

    me = {};
    players = [];

    constructor(comm) {
        this.comm = comm;
    }

    start() {
        let t = this;
        // actual get data from source function
        this.intervalGet = setInterval(() => { t.update() }, 1000);
    }

    update() {
        this.comm.getActive()
            .then(data => this.processData(data))
            .catch(_ => console.warn("GameSettings: Unable to obtain new data."));
    }

    processData(data) {
        this.communityCards = data.communityCards;
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