async function fetchGames() {
  try {
    const response = await fetch("/api/games");
    const gameSessions = await response.json();

    const gamesContainer = document.getElementById("games");
    gamesContainer.innerHTML = "";

    gameSessions.forEach((sessionName) => {
      const gameLink = document.createElement("a");
      gameLink.href = "#";
      gameLink.textContent = sessionName;
      gameLink.onclick = () => openWebSocket(sessionName);

      const listItem = document.createElement("p");
      listItem.appendChild(gameLink);
      gamesContainer.appendChild(listItem);
    });
  } catch (error) {
    console.error("Error fetching game sessions:", error);
  }
}

var lastSession = null;
function openWebSocket(sessionName) {
  if (lastSession) {
    lastSession.close();
  }
  const socket = new WebSocket(`ws://tag.test/api/ws/${sessionName}`);
  lastSession = socket;

  socket.onopen = () => {
    console.log("WebSocket connection established.");
  };

  socket.onmessage = (msg) => {
    let j = JSON.parse(msg.data);

    if (j.action === "init") {
      const gridElement = document.getElementById("game-grid");
      gridElement.outerHTML = j.data.fieldHTML;

      const gameInfoDiv = document.getElementById("game-info");
      gameInfoDiv.innerHTML = "";

      const gameSize = j.data.size;

	 const gameSizeElement = document.createElement("p");
	 gameSizeElement.textContent = `Game Size: ${gameSize} x ${gameSize} (Total Cells: ${gameSize * gameSize})`;
	 gameInfoDiv.appendChild(gameSizeElement);

      const players = j.data.players;
      if (!players || players.length === 0) {
        return;
      }

      players.forEach((player) => {
        const playerInfo = document.createElement("p");
        playerInfo.textContent = `${player.name} - Score: ${player.score} - Position: (${player.x}, ${player.y})`;
        playerInfo.id = `player-${player.name}`;
        if (player.infected) {
          playerInfo.textContent += " - Infected";
          playerInfo.classList.add("infected");
        }
        gameInfoDiv.appendChild(playerInfo);
      });
    } else if (j.action === "move" || j.action === "place") {
      updatePlayer(j.data);
    } else if (j.action === "infect") {
      infectPlayer(j.data);
    }
  };

  socket.onclose = (event) => {
    console.log("WebSocket closed:", event);
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };
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

  player = document.getElementById(`player-${data.name}`);
  if (!player) {
    const gameInfoDiv = document.getElementById("game-info");
    player = document.createElement("p");
    player.id = `player-${data.name}`;
    gameInfoDiv.appendChild(player);
  }
  player.textContent = `${data.name} - Score: ${data.score || 0} - Position: (${data.newX + 1}, ${data.newY + 1})`;
  if (data.infected) {
    player.textContent += " - Infected";
    player.classList.add("infected");
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

window.onload = fetchGames;
