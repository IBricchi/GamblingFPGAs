/// <reference path="./p5/p5.js" />

async function run_setup() {
    let response = await fetch("./setup.json");
    let setup = response.json();
    return setup;
}

async function get_data(setup) {
    console.log(setup);

    let response = await fetch(setup.data_url);
    let data = await response.json();

    document.querySelector("#data")
        .innerHTML = JSON.stringify(data, null, "<br>");
}

run_setup()
    .then(setup => {
        setInterval(get_data, 1000, setup);
    })
    .catch(_ => {
        console.error("Unable to find setup data");
    })