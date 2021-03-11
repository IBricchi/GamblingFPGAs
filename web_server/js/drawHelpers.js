function drawTable() {
    drawCloth();

    const cardw = 55;
    const cardh = 80;

    const slotp = 10;
    const slotw = cardw + slotp;
    const sloth = cardh + slotp;
    const slotm = 20;
    const slotc = color(255);
    const deckc = color(237, 226, 21);

    // draw places for center cards
    const centre_start_x = width / 2 - 2 * slotw - 2 * slotm;
    const centre_start_y = 80;
    for (let i = 0; i < 5; i++) {
        drawCardSlot(centre_start_x + i * slotw + i * slotm, centre_start_y, slotw, sloth, slotc);
    }

    // draw deck
    const deck_x = width / 2;
    const deck_y = centre_start_y + 2 * slotm + slotw;
    drawCardSlot(deck_x, deck_y, sloth, slotw, deckc);
    drawCard(deck_x, deck_y, cardw, cardh, "back", PI / 2);


    // draw places for your cards
    const personal_start_x = width / 2 - slotw / 2 - slotm / 2;
    const personal_start_y = height - 80;
    for (let i = 0; i < 2; i++) {
        drawCardSlot(personal_start_x + i * slotw + i * slotm, personal_start_y, slotw, sloth, slotc);
    }
}

function drawCloth() {
    const c = color(39, 161, 86);
    const cloth_edge = 100;
    fill(c);
    noStroke();
    beginShape();
    vertex(cloth_edge, 0);
    vertex(width - cloth_edge, 0);
    vertex(width, cloth_edge);
    vertex(width, height - cloth_edge);
    vertex(width - cloth_edge, height);
    vertex(cloth_edge, height);
    vertex(0, height - cloth_edge);
    vertex(0, cloth_edge);
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