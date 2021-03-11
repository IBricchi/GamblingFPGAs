/// <reference path="p5/p5.global-mode.d.ts" />
/// <reference path="GetData.js" />
/// <reference path="GameSettings.js" />
/// <reference path="drawHelpers.js" />

let getData;
let game;

let cardTypes = [
    "AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS",
    "AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD",
    "AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC",
    "AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH",
    "back"
]
let cards = {};

function preload() {
    cardTypes.forEach(name => {
        cards[name] = loadImage("./src/cards/" + name + ".png");
    })
}

function setup() {
    // setup canvas
    createCanvas(800, 600);

    // setup drawing params
    rectMode(CENTER);
    imageMode(CENTER);

    // startup data game
    getData = new GetData();
    game = new GameSettings(getData);
    game.start();
}

function draw() {
    const bgcolor = color(29, 117, 63);
    background(bgcolor);

    // draw table
    drawTable();

}