class Lobby extends HTMLDivElement {
    currentCount;

    title;
    playerList;

    constructor() {
        super();
        this.title = document.createElement("h4");
        this.title.innerHTML = "Lobby";
        this.appendChild(this.title);
        this.playerList = document.createElement("div");
        this.appendChild(this.playerList);
        this.currentCount = 0;
        this.id = "lobby";
    }

    clearPlayers() {
        this.playerList.innerHTML = "";
    }

    createPlayer(name) {
        let playerName = document.createElement("span");
        playerName.innerHTML = name;
        return playerName;
    }

    updatePlayers(players) {
        if (players == null) {
            return;
        }
        if (this.currentCount != players.length) {
            this.clearPlayers();
            for (let i = 0; i < players.length; i++) {
                this.playerList.appendChild(this.createPlayer(players[i]
                    .name));
            }
        }
        this.currentCount = players.length;
    }
}
customElements.define('d-lobby', Lobby, { extends: 'div' });

let lobby = false;

function createLobby() {
    if (lobby == false) {
        lobby = document.createElement("div", { is: "d-lobby" });
        document.body.append(lobby);
    }
}

function removeLobby() {
    if (lobby != false) {
        lobby.remove();
    }
    lobby = false;
}