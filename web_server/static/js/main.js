/// <reference path="../p5/p5.js" />

data_path = "/data/"

async function get_data(){
    let response = await fetch(data_path);
    let data = await response.json();
    
    document.querySelector("#data").innerHTML = JSON.stringify(data, null, "<br>");
}

setInterval(get_data, 1000);