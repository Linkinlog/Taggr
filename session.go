package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func NewSession(size int, ctx context.Context) *Session {
	field := make([][]*Player, size)
	for i := 0; i < size; i++ {
		field[i] = make([]*Player, size)
		for j := 0; j < size; j++ {
			field[i][j] = nil
		}
	}

	s := &Session{
		field:   field,
		size:    size,
		players: make(map[string]*Player),

		Ws:          NewWebSocketHandler(),
		fieldMutex:  &sync.Mutex{},
		playerMutex: &sync.Mutex{},
	}

	bots := []string{"maggie bot", "bandit bot", "billy bot", "bobby bot", "benny bot", "bandy bot", "bendy bot", "benky bot", "bency bot", "benjy bot", "maggie bot 2", "bandit bot 2", "billy bot 2", "bobby bot 2", "benny bot 2", "bandy bot 2", "bendy bot 2", "bekny bot 2", "bency bot 2", "benjy bot 2", "maggie bot 3", "bandit bot 3", "billy bot 3", "bobby bot 3", "benny bot 3", "bandy bot 3", "bendy bot 3", "bekny bot 3", "bency bot 3", "benjy bot 3"}

	// 2 bots per 10 size
	for i := 0; i < size/5; i++ {
		if i >= len(bots) {
			break
		}
		bot := NewPlayer(bots[i])

		s.AddPlayer(bot)
		s.PlacePlayerRandomly(bot)

		go func() {
			ticker := time.NewTicker(time.Duration(time.Duration(rand.Intn(1000))*time.Millisecond) + 250*time.Millisecond)
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					switch rand.Intn(3) {
					case 0:
						s.MovePlayerDown(bot)
						s.MovePlayerLeft(bot)
					case 1:
						s.MovePlayerUp(bot)
						s.MovePlayerRight(bot)
					case 2:
						s.MovePlayerLeft(bot)
					case 3:
						s.MovePlayerRight(bot)
					}
				}
			}
		}()
	}

	return s
}

type Session struct {
	field   [][]*Player
	players map[string]*Player
	size    int

	Ws          *WebSocketHandler
	fieldMutex  *sync.Mutex
	playerMutex *sync.Mutex
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
	s.fieldMutex.Lock()
	defer s.fieldMutex.Unlock()
	if !s.CheckPlayerPosition(p) {
		panic("hey we need to fix moving players")
	}
	if s.field[x][y] != nil {
		if p.Infected() && !s.field[x][y].Infected() {
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
		} else if s.field[x][y].Infected() && !p.Infected() {
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
			Score:    p.Score(),
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
	for !placed {
		fmt.Println("trying to place player")
		fmt.Println(p.Name())
		x := rand.Intn(s.size)
		y := rand.Intn(s.size)

		if !s.PlacePlayer(p, x, y) {
			continue
		}
		placed = true
		fmt.Println("placed player")
		fmt.Println(p.Name())
	}

	if !placed {
		panic("hey we need to fix placing players")
	}
}

func (s *Session) PlacePlayer(p *Player, x, y int) (placed bool) {
	s.fieldMutex.Lock()
	defer s.fieldMutex.Unlock()
	placed = false
	if s.field[x][y] == nil {
		p.SetPosition(x, y)
		s.field[x][y] = p
		placed = true
		s.Ws.Broadcast(&Event{
			Data: &PlayerMoveResponse{
				Name:     p.Name(),
				Score:    p.Score(),
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
	s.playerMutex.Lock()
	defer s.playerMutex.Unlock()
	if len(s.players) == 0 {
		p.Infect()
	}
	s.players[p.ID()] = p
}

func (s *Session) RemovePlayer(p *Player) {
	s.playerMutex.Lock()
	defer s.playerMutex.Unlock()
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
