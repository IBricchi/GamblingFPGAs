// let q;
class GameSettings {
    refreshRate = 1000;

    state;

    lobby;

    availableMovesLength = 0;

    comm = {};
    intervalGet;

    me = {};
    players = [];
    winners = [];

    communityCards;
    currentPlayer;
    currentPot;
    currentRound;
    winningPots;
    winningReason;

    constructor(comm) {
        this.comm = comm;
        this.state = "login"
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
        this.state = "open"
        if (form.name != "OPEN" && !data.open) {
            removeForm();
            removeLobby();
            createOpenForm();
            removeTerminate();
        }
        if (form.name != "JOIN" && data.open && (data.players == null || !data.players.some(element => element.name == comm.username))) {
            removeForm();
            createLobby();
            createJoinForm();
            createTerminate();
        }
        if (form.name != "START" && data.players != null && data.players.some(element => element.name == comm.username)) {
            removeForm();
            createLobby();
            createStartForm();
            createTerminate();
        }
        if (data.open) {
            createLobby();
            lobby.updatePlayers(data.players);
        }
        if (data.active) {
            removeForm();
            removeLobby();
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
        this.state = "active";
        if (data.hasEnded) {
            clearInterval(this.intervalGet);
            let t = this;
            t.updateShowdown();
            this.intervalGet = setInterval(() => { t.updateShowdown() }, t.refreshRate);
            createTerminate();
        }
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
        this.leaglMoves = data.availableNextMoves;
        this.sbb = data.smallBlindValue;

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

        if (this.currentPlayer == this.me.order)
            updateLeaglMoves(this, data.availableNextMoves);
        else
            updateLeaglMoves(this, ["peek"]);

    }

    updateShowdown() {
        this.comm.getShowdown()
            .then(data => this.processShowdown(data))
            .catch(err => {
                console.warn(err);
                console.warn("GameSettings: Unable to obtain new data.")
            });
    }

    processShowdown(data) {
        this.state = "showdown"
        createContinueGame();
        if (data.newGameStarted) {
            removeContinueGame();
            clearInterval(this.intervalGet);
            let t = this;
            t.updateActive();
            this.intervalGet = setInterval(() => { t.updateActive() }, t.refreshRate);
            createTerminate();
        }
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
        this.players = [];
        this.winners = [];
        this.winningReason = data.winningReason;
        this.winningPots = data.winningMoneyAmounts;
        this.currentPot = data.potMoneyAmount;
        let k = 0;
        for (let i = 0; i < data.players.length; i++) {
            let player = data.players[i];
            player.order = i;
            player.winner = false;
            for (let j = 0; j < data.winners.length; j++) {
                if (player.name == data.winners[j].name) {
                    player.winnerOrder = k;
                    player.winner = true;
                    this.winners.push(player);
                    k++;
                }
            }
            if (player.name == comm.username) {
                this.me = player;
            }
            this.players.push(player);
        }
    }

}