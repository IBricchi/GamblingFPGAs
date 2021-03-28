// draw settings
class DS {
    constructor() {
        this.canvasw = 1200; // canvas width
        this.canvash = 600; // canvas height

        this.tablew = 700; // table width
        this.tableh = this.canvash; // table height
        this.tablec = color(29, 117, 63); // table color
        this.ttexts = 60; // table text size

        this.clothm = 100; // cloth margin (for corners)
        this.clothc = color(39, 161, 86); // cloth color

        this.playertw = this.canvasw - this.tablew; // player table width
        this.playerth = this.canvash; // player table height
        this.playertc = color(255); // player table background color

        this.playercw = this.playertw; // player card width
        this.playerch = this.playerth / 6; // player card height
        this.playerco = color(0); // player card outline color
        this.playercto = color(252, 186, 3); // player card outline color if turn
        this.playerc = [
            color(255, 0, 0),
            color(0, 255, 0),
            color(0, 0, 255),
            color(255, 255, 0),
            color(255, 0, 255),
            color(0, 255, 255)
        ]; // possible player colors
        this.picons = this.playerch - 30; // player icon size
        this.pcardl = 10; // player label size
        this.pnames = this.playerch / 3; // size of name text
        this.ptexts = this.playerch / 6; // size of rest of player text
        this.pfoldc = color(0, 0, 0, 150) // tint on folded players

        this.cardw = 62; // card width
        this.cardh = 90; // card height

        this.slotp = 10; // card slot padding (padding around the card)
        this.slotw = this.cardw + this.slotp; // card slot width
        this.sloth = this.cardh + this.slotp; // card slot height
        this.slotm = 20; // card slot margin
        this.slotc = color(255); // card slot colour
        this.deckc = color(237, 226, 21); // deck card slot colour

        this.winnerbh = this.canvash - 3 * this.cardh;
        this.winnerbw = 50;
        this.winnerts = 30;
    }
};

function drawTable(game) {
    drawCloth();

    // draw places for center cards
    const centre_start_x = ds.tablew / 2 - 2 * ds.slotw - 2 * ds.slotm;
    const centre_start_y = 80;
    for (let i = 0; i < 5; i++) {
        drawCardSlot(centre_start_x + i * ds.slotw + i * ds.slotm, centre_start_y, ds.slotw, ds.sloth, ds.slotc);
    }

    const comm_on_round = [0, 0, 3, 4, 5];
    for (let i = 0; i < comm_on_round[game.currentRound]; i++) {
        drawCard(centre_start_x + i * ds.slotw + i * ds.slotm, centre_start_y, ds.cardw, ds.cardh, game.communityCards[i]);
    }

    // draw deck
    const deck_x = ds.tablew / 2;
    const deck_y = centre_start_y + 2 * ds.slotm + ds.slotw;
    drawCardSlot(deck_x, deck_y, ds.sloth, ds.slotw, ds.deckc);
    drawCard(deck_x, deck_y, ds.cardw, ds.cardh, "back", PI / 2);


    // draw places for your cards
    const personal_start_x = ds.tablew / 2 - ds.slotw / 2 - ds.slotm / 2;
    const personal_start_y = ds.tableh - 80;

    let myCards = ["back", "back"];
    if (game.me.showCardsMe)
        myCards = game.me.hand;
    for (let i = 0; i < 2; i++) {
        drawCardSlot(personal_start_x + i * ds.slotw + i * ds.slotm, personal_start_y, ds.slotw, ds.sloth, ds.slotc);
        drawCard(personal_start_x + i * ds.slotw + i * ds.slotm, personal_start_y, ds.cardw, ds.cardh, myCards[i]);
    }

    // draw text stuff
    noStroke();
    fill(color(0));
    textSize(ds.ttexts);
    text("Round: " + game.currentRound + "\nPot: £" + game.currentPot, ds.tablew / 2, ds.tableh * 6 / 11);
}

function drawShowdown(game) {
    drawCloth();

    // draw places for center cards
    const centre_start_x = ds.tablew / 2 - 2 * ds.slotw - 2 * ds.slotm;
    const centre_start_y = 80;
    for (let i = 0; i < 5; i++) {
        drawCardSlot(centre_start_x + i * ds.slotw + i * ds.slotm, centre_start_y, ds.slotw, ds.sloth, ds.slotc);
    }

    for (let i = 0; i < 5; i++) {
        drawCard(centre_start_x + i * ds.slotw + i * ds.slotm, centre_start_y, ds.cardw, ds.cardh, game.communityCards[i]);
    }

    // draw money bar
    let winnerMargins = drawMoneyBar(game.winners, game.winningPots, game.currentPot, ds.cardh, ds.cardh * 2);
    for (let i = 0; i < game.winners.length; i++) {
        drawWinnerPlayerData(game.winners[i], game.winningPots[i], ds.cardh + 3 / 2 * ds.winnerbw, ds.cardh * 2 + winnerMargins[i] - ds.playerch / 2);
    }

    // draw winning reason
    fill(color(0));
    textSize(ds.winnerts);
    text(game.winningReason, ds.tablew / 2, ds.tableh - ds.winnerts);
}

function drawPlayerTable(game) {
    push();
    translate(ds.tablew, 0);
    fill(ds.playertc);
    rect(ds.playertw / 2, ds.playerth / 2, ds.playertw, ds.playerth);

    for (let i = 0; i < game.players.length; i++) {
        drawPlayerCard(game, game.players[i], 0, ds.playerch * i);
    }

    // outline who's turn it is
    strokeWeight(5);
    stroke(ds.playercto);
    noFill();
    rect(ds.playercw / 2, ds.playerch * game.currentPlayer + ds.playerch / 2, ds.playercw, ds.playerch);

    pop();
}

function drawCloth() {
    fill(ds.clothc);
    noStroke();
    beginShape();
    vertex(ds.clothm, 0);
    vertex(ds.tablew - ds.clothm, 0);
    vertex(ds.tablew, ds.clothm);
    vertex(ds.tablew, ds.tableh - ds.clothm);
    vertex(ds.tablew - ds.clothm, ds.tableh);
    vertex(ds.clothm, ds.tableh);
    vertex(0, ds.tableh - ds.clothm);
    vertex(0, ds.clothm);
    endShape();

}

function drawCardSlot(x, y, w, h, c) {
    push();
    translate(-w / 2, -h / 2);

    const sw = 5;
    strokeWeight(sw);
    noFill();
    stroke(c);
    // draw left side
    beginShape();
    vertex(x + w / 3, y);
    vertex(x, y);
    vertex(x, y + h);
    vertex(x + w / 3, y + h);
    endShape();
    // draw right side
    beginShape();
    vertex(x + 2 * w / 3, y);
    vertex(x + w, y);
    vertex(x + w, y + h);
    vertex(x + 2 * w / 3, y + h);
    endShape();

    pop();
}

function drawCard(x, y, w, h, v, t = 0) {
    push();
    translate(x, y);
    rotate(t);
    image(cards[v], 0, 0, w, h);
    pop();
}

function drawMoneyBar(winners, pots, pot, x, y) {
    push();
    translate(x, y);
    rectMode(CORNER);

    let out = [];
    let avaiabley = 0;
    noStroke();
    for (let i = 0; i < winners.length; i++) {
        fill(ds.playerc[winners[i].order]);
        rect(0, avaiabley, ds.winnerbw, ds.winnerbh * pots[i] / pot);
        out.push(avaiabley + ds.winnerbh * pots[i] / pot / 2);
        avaiabley += ds.winnerbh * pots[i] / pot;
    }

    pop();

    outMargins = [];
    return out;
}

function drawWinnerPlayerData(winner, pot, x, y) {
    push();
    translate(x, y);

    // used to dynamically manage how much of the x axis has been used up
    let avaiablex = 0;

    // draw cards
    drawCard(10 + ds.cardw / 2, ds.playerch / 2, ds.cardw, ds.cardh, winner.hand[0]);
    drawCard(10 + ds.cardw * 3 / 2, ds.playerch / 2, ds.cardw, ds.cardh, winner.hand[1]);
    avaiablex += 2 * ds.cardw;

    // remaining space left
    translate(avaiablex, 0);
    let remainingw = ds.playercw - avaiablex;

    // draw name
    noStroke();
    fill(0);
    textSize(ds.cardh / 2);
    text(winner.name + ": £" + pot, remainingw / 2, ds.playerch / 2 + ds.cardh / 16);


    pop();
}

function drawPlayerCard(game, player, x, y, turn) {
    push();
    translate(x, y);

    // used to dynamically manage how much of the x axis has been used up
    let avaiablex = 0;

    // draw background
    stroke(ds.playerco);
    noFill();
    rect(ds.playercw / 2, ds.playerch / 2, ds.playercw, ds.playerch);

    // draw two cards
    let card1 = "back";
    let card2 = "back";
    if (player.hand != null && player.name != game.me.name) {
        card1 = player.hand[0];
        card2 = player.hand[1];
    }
    if (player.showCardsMe && player.name == game.me.name) {
        card1 = player.hand[0];
        card2 = player.hand[1];
    }
    drawCard(10 + ds.cardw / 2, ds.playerch / 2, ds.cardw, ds.cardh, card1);
    drawCard(10 + ds.cardw * 3 / 2, ds.playerch / 2, ds.cardw, ds.cardh, card2);
    avaiablex += 2 * ds.cardw;

    // draw icon
    tint(ds.playerc[player.order]);
    image(icons.user, avaiablex + ds.playerch / 2, ds.playerch / 2, ds.picons, ds.picons);
    avaiablex += ds.picons;

    // draw labels if necessary
    textSize(ds.pcardl);
    if (player.dealer) {
        text("D", avaiablex, ds.cardh * 4 / 5);
    } else if (player.smallBlind) {
        text("SB", avaiablex, ds.cardh * 4 / 5);
    } else if (player.bigBlind) {
        text("BB", avaiablex, ds.cardh * 4 / 5);
    }

    // remaining space left
    translate(avaiablex, 0);
    let remainingw = ds.playercw - avaiablex;

    // draw name
    noStroke();
    fill(0);
    textSize(ds.pnames);
    text(player.name + ": £" + player.moneyAvailableAmount, remainingw / 2, ds.pnames);

    // draw bet and money
    textSize(ds.ptexts);
    text("Bet Total/Last: £" + player.totalMoneyBetAmount + "/" + player.lastBetAmount, remainingw / 2, ds.playerch - ds.ptexts / 2);

    // draw failed attempts
    if (player.name == game.me.name)
        text("FA: " + player.failedPeekAttemptsCurrentGame, ds.picons / 3, ds.playerch - ds.ptexts / 2);

    translate(-avaiablex, 0);

    if (player.hasFolded) {
        noStroke();
        fill(ds.pfoldc);
        rect(ds.playercw / 2, ds.playerch / 2, ds.playercw, ds.playerch);
    }

    pop();
}

function updateLeaglMoves(game, moves) {
    if (game.availableMovesLength != moves.length) {
        document.querySelector("#mnb")
            .innerHTML = "Minimum Bet Ammount: £" + game.mnb
        document.querySelector("#sbb")
            .innerHTML = "Small Blind Bet: £" + game.sbb;
        let ml = document.querySelector("#avaiableMovesList");
        ml.innerHTML = "";
        moves.forEach(move => {
            let mli = document.createElement("li");
            mli.innerText = move;
            ml.append(mli);
        });
        availableMovesLength = moves.length;
    }
}