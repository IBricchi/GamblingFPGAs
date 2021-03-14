/// <reference path="p5/p5.global-mode.d.ts" />
/// <reference path="Communicator.js" />
/// <reference path="GameSettings.js" />
/// <reference path="drawHelpers.js" />
/// <reference path="forms.js" />

let canvas;

let form = {};
let terminate = null;

let comm;
let game;
let ds;

let iconTypes = [
    "user"
];
let icons = {};

let cardTypes = [
    "As", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s", "Ts", "Js", "Qs", "Ks",
    "Ad", "2d", "3d", "4d", "5d", "6d", "7d", "8d", "9d", "Td", "Jd", "Qd", "Kd",
    "Ac", "2c", "3c", "4c", "5c", "6c", "7c", "8c", "9c", "Tc", "Jc", "Qc", "Kc",
    "Ah", "2h", "3h", "4h", "5h", "6h", "7h", "8h", "9h", "Th", "Jh", "Qh", "Kh",
    "back"
];
let cards = {};

let inputs = {};

function preload() {
    cardTypes.forEach(name => {
        cards[name] = loadImage("./src/cards/" + name + ".png");
    })
    iconTypes.forEach(name => {
        icons[name] = loadImage("./src/icons/" + name + ".png");
    })
}

function setup() {
    // startup data game
    comm = new Communicator("player1", "player1");
    game = new GameSettings(comm);
    ds = new DS();

    // setup canvas
    canvas = createCanvas(ds.canvasw, ds.canvash);
    canvas.hide();

    // setup drawing params
    rectMode(CENTER);
    imageMode(CENTER);
    textAlign(CENTER);

    // create login page
    createLogin();

    // stop looping
    noLoop();
}

function draw() {
    // if game is active
    background(ds.tablec);

    // draw table
    drawTable(game);

    // draw player table
    drawPlayerTable(game);
}