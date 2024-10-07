let socket = new WebSocket("ws://tag.test/api/ws/testing");
console.log("Websocket started.");

socket.onOpen = () => {
  socket.send("Hello, World!");
  console.log("Client started.");
};

socket.onclose = (event) => {
  console.log("Socket closed: ", event);
};

socket.onError = (error) => {
  console.log("Socket Error: ", error);
};

socket.onmessage = (msg) => {
  console.log(msg);
  let j = JSON.parse(msg.data);

  if (j.action === "init") {
    console.log(j);
    createGrid(j.data.field, j.data.players);
  } else if (j.action === "move" || j.action === "place") {
    console.log(j);
    updatePlayer(j.data);
  } else if (j.action === "infect") {
    console.log(j);
    infectPlayer(j.data);
  }
};

function createGrid(field, players) {
  const gridElement = document.getElementById("game-grid");
  gridElement.innerHTML = "";

  field.forEach((row, x) => {
    row.forEach((_, y) => {
      const div = document.createElement("div");
      div.classList.add("grid-cell");
      div.setAttribute("data-x", x);
      div.setAttribute("data-y", y);

      if (players) {
        // TODO feels like we can do this better
        const player = players.find((p) => p.x === x && p.y === y);
        if (player) {
          div.innerText = player.name;
          div.classList.add("player");
          if (player.infected) {
            div.classList.add("infected");
          }
        }
      }

      gridElement.appendChild(div);
    });
  });
}
function updatePlayer(data) {
  const oldCell = document.querySelector(
    `[data-x="${data.x}"][data-y="${data.y}"]`,
  );
  if (oldCell) {
    oldCell.innerText = "";
    oldCell.classList.remove("player", "infected");
  }

  const newCell = document.querySelector(
    `[data-x="${data.newX}"][data-y="${data.newY}"]`,
  );
  if (newCell) {
    newCell.innerText = data.name;
    newCell.classList.add("player");
    if (data.infected) {
      newCell.classList.add("infected");
    }
  }
}

function infectPlayer(data) {
  const cell = document.querySelector(
    `[data-x="${data.x}"][data-y="${data.y}"]`,
  );
  if (cell) {
    cell.classList.add("infected");
  }
}
