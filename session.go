package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

func NewSession(size int) *Session {
	field := make([][]*Player, size)
	for i := 0; i < size; i++ {
		field[i] = make([]*Player, size)
		for j := 0; j < size; j++ {
			field[i][j] = nil
		}
	}

	ws := NewWebSocketHandler()

	return &Session{
		field:   field,
		size:    size,
		players: make(map[string]*Player),

		Ws: ws,
		m:  &sync.Mutex{},
	}
}

type Session struct {
	field   [][]*Player
	players map[string]*Player
	size    int

	Ws *WebSocketHandler
	m  *sync.Mutex
}

func (s *Session) GetPlayer(id string) *Player {
	return s.players[id]
}

func (s *Session) CheckPlayerPosition(p *Player) (found bool) {
	found = false

	x, y := p.Postition()
	if s.field[x][y] == p {
		found = true
	}

	return found
}

func (s *Session) MovePlayerDown(p *Player) {
	x, y := p.Postition()
	if x+1 < s.size {
		s.MovePlayer(p, x+1, y)
	}
}

func (s *Session) MovePlayerUp(p *Player) {
	x, y := p.Postition()
	if x-1 >= 0 {
		s.MovePlayer(p, x-1, y)
	}
}

func (s *Session) MovePlayerLeft(p *Player) {
	x, y := p.Postition()
	if y-1 >= 0 {
		s.MovePlayer(p, x, y-1)
	}
}

func (s *Session) MovePlayerRight(p *Player) {
	x, y := p.Postition()
	if y+1 < s.size {
		s.MovePlayer(p, x, y+1)
	}
}

func (s *Session) MovePlayer(p *Player, x, y int) {
	if !s.CheckPlayerPosition(p) {
		panic("hey we need to fix moving players")
	}
	if s.field[x][y] != nil {
		if p.Infected() {
			p.AddScore()
			event := &Event{
				Action: "score",
				Data: &PlayerResponse{
					Name:     p.Name(),
					Score:    p.Score(),
					Infected: p.Infected(),
					X:        p.X(),
					Y:        p.Y(),
				},
			}

			s.Ws.Broadcast(event)
			s.field[x][y].Infect()
			s.Ws.Broadcast(&Event{
				Data: &PlayerResponse{
					Name:     s.field[x][y].Name(),
					Score:    s.field[x][y].Score(),
					Infected: s.field[x][y].Infected(),
					X:        s.field[x][y].X(),
					Y:        s.field[x][y].Y(),
				},
				Action: "infect",
			})
		} else if s.field[x][y].Infected() {
			p.Infect()
			s.Ws.Broadcast(&Event{
				Data: &PlayerResponse{
					Name:     p.Name(),
					Score:    p.Score(),
					Infected: p.Infected(),
					X:        p.X(),
					Y:        p.Y(),
				},
				Action: "infect",
			})
			s.field[x][y].AddScore()
			s.Ws.Broadcast(&Event{
				Data: &PlayerResponse{
					Name:     s.field[x][y].Name(),
					Score:    s.field[x][y].Score(),
					Infected: s.field[x][y].Infected(),
					X:        s.field[x][y].X(),
					Y:        s.field[x][y].Y(),
				},
				Action: "score",
			})
		}
		return
	}
	s.field[p.x][p.y] = nil
	s.field[x][y] = p
	s.Ws.Broadcast(&Event{
		Data: &PlayerMoveResponse{
			Name:     p.Name(),
			X:        p.X(),
			Y:        p.Y(),
			NewX:     x,
			NewY:     y,
			Infected: p.Infected(),
		},
		Action: "move",
	})
	p.SetPosition(x, y)
}

func (s *Session) PlacePlayerRandomly(p *Player) {
	placed := false
	for i := 0; i < 10; i++ {
		x := rand.Intn(s.size)
		y := rand.Intn(s.size)

		if s.PlacePlayer(p, x, y) {
			placed = true
			break
		}
	}

	if !placed {
		panic("hey we need to fix placing players")
	}
}

func (s *Session) PlacePlayer(p *Player, x, y int) (placed bool) {
	placed = false
	if s.field[x][y] == nil {
		p.SetPosition(x, y)
		s.field[x][y] = p
		s.AddPlayer(p)
		placed = true
		s.Ws.Broadcast(&Event{
			Data: &PlayerMoveResponse{
				Name:     p.Name(),
				Infected: p.Infected(),
				X:        p.X(),
				Y:        p.Y(),
				NewX:     x,
				NewY:     y,
			},
			Action: "place",
		})
	}
	return placed
}

func (s *Session) AddPlayer(p *Player) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.players) == 0 {
		p.Infect()
	}
	s.players[p.ID()] = p
}

func (s *Session) RemovePlayer(p *Player) {
	s.m.Lock()
	defer s.m.Unlock()
	if s.CheckPlayerPosition(p) {
		x, y := p.Postition()
		delete(s.players, p.ID())
		s.field[x][y] = nil

		// TODO i hate this
		infected := rand.Intn(len(s.players))
		i := 0
		for _, p := range s.players {
			if i == infected {
				p.Infect()
			}
			i++
		}

		return
	}
}

func (s *Session) Players() []*Player {
	var players []*Player
	for _, p := range s.players {
		players = append(players, p)
	}

	return players
}

func (s *Session) Field() [][]string {
	field := make([][]string, s.size)
	for i := 0; i < s.size; i++ {
		field[i] = make([]string, s.size)
		for j := 0; j < s.size; j++ {
			if s.field[i][j] != nil {
				field[i][j] = s.field[i][j].Name()
			} else {
				field[i][j] = ""
			}
		}
	}

	return field
}

func (s *Session) FieldSize() int {
	return s.size
}

func (s *Session) String() string {
	var buf strings.Builder
	buf.WriteString("Players: ")
	for _, p := range s.players {
		buf.WriteString(p.String())
		buf.WriteString(" ")
	}
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("PlayerCount: %d\n", len(s.players)))
	for i := 0; i < len(s.field); i++ {
		buf.WriteString("[")
		for j := 0; j < len(s.field[i]); j++ {
			if s.field[i][j] != nil {
				buf.WriteString(s.field[i][j].String())
			} else {
				buf.WriteString(" ")
			}

			if j < len(s.field[i])-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString("]")
		buf.WriteString("\n")
	}

	return buf.String()
}
