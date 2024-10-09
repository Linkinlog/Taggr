# Taggr

Taggr is a simple, highly concurrent, event driven, multiplayer infection game where players navigate a square grid. The objective varies depending on the player's infection status: either infect other players or avoid being infected.

## Features

- Multiple Sessions: Players can create and join different game sessions.
- **Multiplayer Gameplay**: Engage with multiple players in real-time.
- **Infection Mechanics**: Players can either infect others or avoid getting infected based on their status.
- **Highly Concurrent**: Built with Go, ensuring efficient handling of multiple connections.
- **WebSocket Communication**: Real-time updates and interactions between the server and clients.
- **Vanilla JavaScript**: Client-side implementation using plain JavaScript for simplicity and performance.
- **Event-Driven Design**: Efficiently manages game state changes and player interactions.
- **Docker Support**: Easily set up and run the game locally using Docker.
- **Caddy for Reverse Proxying**: Simplifies the server deployment and management.

## Getting Started

### Prerequisites

- Docker installed on your machine.
- A modern web browser to run the client.

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/linkinlog/taggr.git
   cd taggr
   ```

2. **Build and run the Docker container**:
   ```bash
   docker-compose up --build
   ```

3. **Access the game**:
    - Set your DNS to point `*.test` to `127.0.0.1` or simply modify your host file to add `tag.test 127.0.0.1`

### Playing the Game

TODO - For now, check out `http.go` for the routing, you can create games, players, and move around, through HTTP requests.

## Technologies Used

- **Go**: Server-side logic and concurrent handling.
- **WebSockets**: Real-time communication between server and clients.
- **JavaScript**: Client-side interactions.
- **Caddy**: For reverse proxying and serving the application.
- **Docker**: For local development and deployment.

## Contributing

Contributions are welcome! If you have suggestions or improvements, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
