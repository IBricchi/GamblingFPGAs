function createTerminate() {
    if (!terminate || !document.querySelector("#terminate")) {
        terminate = createButton("Terminate");
        terminate.id("terminate");
        terminate.mousePressed(() => { comm.terminate() });
    }
}

function removeTerminate() {
    if (terminate) {
        terminate.remove();
    }
}

function createContinueGame() {
    if (!continueGame || !document.querySelector("#continueGame")) {
        continueGame = createButton("Continue Game");
        continueGame.id("continueGame");
        continueGame.mousePressed(() => { comm.sendContinueGame() });
    }
}

function removeContinueGame() {
    if (continueGame) {
        continueGame.remove();
    }
}

function createLogin() {
    form = {
        name: "LOGIN",
        container: createElement("form"),
        usernameLabel: createElement("label", "Username: "),
        username: createInput(),
        passwordLabel: createElement("label", "Password: "),
        password: createInput("", "password"),
        br: createElement("br"),
        submit: createInput("Login", "button"),
        info: createElement("p", "Please wait a second to see changes, do not submit multiple times.")
    }
    form.submit.id("submit")
    form.usernameLabel.parent(form.container);
    form.username.parent(form.container);
    form.passwordLabel.parent(form.container);
    form.password.parent(form.container);
    form.br.parent(form.container);
    form.submit.parent(form.container);
    form.submit.mousePressed(checkCredentials);
}

function checkCredentials() {
    let username = form.username.value();
    let password = form.password.value();
    comm.getCheckCred(username, password)
        .then(b => {
            if (b) {
                removeForm();
                createOpenForm();
                comm.setAuthentication(username, password);
                game.start();
            } else {
                alert("Password Incorrect!");
            }
        });
}

function createOpenForm() {
    form = {
        name: "OPEN",
        container: createElement("form"),
        smLabel: createElement("label", "Starting Money: "),
        sm: createInput(2000, "number"),
        sbbLabel: createElement("label", "Small Blind Bet: "),
        sbb: createInput(5, "number"),
        br: createElement("br"),
        submit: createInput("Open Game", "button"),
        info: createElement("p", "Please wait a second to see changes, do not submit multiple times.")
    };
    form.submit.id("submit")
    form.smLabel.parent(form.container);
    form.sm.parent(form.container);
    form.sbbLabel.parent(form.container);
    form.sbb.parent(form.container);
    form.br.parent(form.container);
    form.submit.parent(form.container);
    form.submit.mousePressed(checkValues);
}

function checkValues() {
    let sm = form.sm.value();
    let sbb = form.sbb.value();
    if (isNaN(sm) || isNaN(sbb)) {
        alert("Invalid inpus, expecting numbers.");
    } else if (sm <= 0 || sbb <= 0) {
        alert("Starting monery must be greater than 0, and Small Blind Bet must be 0 or greater.");
    } else {
        comm.sendOpen(sm, sbb);
    }
}

function createJoinForm() {
    form = {
        name: "JOIN",
        container: createElement("form"),
        submit: createInput("Join", "button"),
        info: createElement("p", "Please wait a second to see changes, do not submit multiple times.")
    }
    form.submit.parent(form.container);
    form.submit.mousePressed(joinGame);
}

function joinGame() {
    comm.sendJoin();
}

function createStartForm() {
    form = {
        name: "START",
        container: createElement("form"),
        submit: createInput("Start", "button"),
        info: createElement("p", "Please wait a second to see changes, do not submit multiple times.")
    }
    form.submit.parent(form.container);
    form.submit.mousePressed(startGame);
}

function startGame() {
    comm.sendStart();
}

function removeForm() {
    Object.keys(form)
        .forEach(key => {
            if (key != "name")
                form[key].remove();
        })
    form = {};
}