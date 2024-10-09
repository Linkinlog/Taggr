package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

func ServeHTTP(g *Game) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("GET /api/ws/{sessionID}", websocket.Handler(func(ws *websocket.Conn) {
		sessionID := ws.Request().PathValue("sessionID")

		sess := g.GetSession(sessionID)

		if sess == nil {
			ws.Close()
			return
		}

		sess.Ws.AddSocket(ws)

		fmt.Println("new websocket connection")

		var players []PlayerResponse

		for _, player := range sess.Players() {
			players = append(players, PlayerResponse{
				Name:     player.Name(),
				Score:    player.Score(),
				Infected: player.Infected(),
				X:        player.X(),
				Y:        player.Y(),
			})
		}

		responder := &FieldResponse{
			Players:   players,
			Field:     sess.Field(),
			FieldHTML: sess.FieldHTML(),
			Size:      sess.FieldSize(),
		}

		event := &Event{
			Data:   responder,
			Action: "init",
		}

		eventBytes, err := event.Response()
		if err != nil {
			ws.Close()
			return
		}

		fmt.Println("sending init event")
		if _, err := ws.Write(eventBytes); err != nil {
			fmt.Println("error sending init event")
			fmt.Println(err)
			ws.Close()
			return
		}

		sess.Ws.Open(ws.Request().Context())
	}))

	mux.HandleFunc("GET /api/games", func(w http.ResponseWriter, r *http.Request) {
		var games []string
		for name := range g.sessions {
			games = append(games, name)
		}

		if err := json.NewEncoder(w).Encode(games); err != nil {
			http.Error(w, "error encoding games", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /api/games", func(w http.ResponseWriter, r *http.Request) {
		size := 10

		sizeInput := r.FormValue("size")
		if sizeInput != "" {
			s, err := strconv.Atoi(sizeInput)
			if err == nil {
				size = s
			}
		}

		name := g.AddSession(size)

		resp := struct {
			Name string `json:"name"`
		}{
			Name: name,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /api/games/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		sess := g.GetSession(id)
		if sess == nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		resp := &FieldResponse{
			Field: sess.Field(),
			FieldHTML: sess.FieldHTML(),
			Size:  sess.FieldSize(),
		}

		players := sess.Players()
		for _, player := range players {
			resp.Players = append(resp.Players, PlayerResponse{
				Name:     player.Name(),
				Score:    player.Score(),
				Infected: player.Infected(),
				X:        player.X(),
				Y:        player.Y(),
			})
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /api/games/{id}/players", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		sess := g.GetSession(id)
		if sess == nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		resp := []PlayerResponse{}

		players := sess.Players()
		for _, player := range players {
			resp = append(resp, PlayerResponse{
				Name:     player.Name(),
				Score:    player.Score(),
				Infected: player.Infected(),
				X:        player.X(),
				Y:        player.Y(),
			})
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("POST /api/games/{id}/players", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "missing name", http.StatusBadRequest)
			return
		}

		sess := g.GetSession(id)
		if sess == nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		for _, player := range sess.Players() {
			if player.Name() == name {
				http.Error(w, "player already exists", http.StatusBadRequest)
				return
			}
		}

		player := NewPlayer(name)
		sess.AddPlayer(player)
		sess.PlacePlayerRandomly(player)

		resp := struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID:   player.ID(),
			Name: player.Name(),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /api/games/{id}/players/{player}/move/{direction}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		playerID := r.PathValue("player")
		if playerID == "" {
			http.Error(w, "missing player", http.StatusBadRequest)
			return
		}

		direction := r.PathValue("direction")
		if direction == "" {
			http.Error(w, "missing direction", http.StatusBadRequest)
			return
		}

		sess := g.GetSession(id)
		if sess == nil {
			http.Error(w, "session not found", http.StatusNotFound)
			return
		}

		player := sess.GetPlayer(playerID)
		if player == nil {
			http.Error(w, "player not found", http.StatusNotFound)
			return
		}

		switch direction {
		case "up":
			sess.MovePlayerUp(player)
		case "down":
			sess.MovePlayerDown(player)
		case "left":
			sess.MovePlayerLeft(player)
		case "right":
			sess.MovePlayerRight(player)
		default:
			http.Error(w, "invalid direction", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":420", mux)
}
