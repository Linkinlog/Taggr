package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func NewGame() *Game {
	return &Game{
		sessions: make(map[string]*Session),
	}
}

type Game struct {
	sessions map[string]*Session
}

var (
	namePart1 = []string{"Super", "Mega", "Ultra", "Hyper", "Mighty", "Powerful", "Speedy", "Extremely", "Fantastic", "Incredible"}
	namePart2 = []string{"Motorized", "Flying", "Jumping", "Rolling", "Bouncing", "Exploding", "Falling", "Crawling", "Swimming", "Running"}
	namePart3 = []string{"Turtles", "Balls", "Cubes", "Spheres", "Pyramids", "Cylinders", "Cones", "Planes", "Helicopters", "Cars"}
	namePart4 = []string{"of-Doom", "of-Death", "of-Destruction", "of-Despair", "of-Darkness", "of-Desolation", "of-Dreams", "of-Delight", "of-Daring", "of-Daring"}
)

func (g *Game) AddSession(size int) (name string) {
	if size > 100 {
		size = 100
	}
	added := false
	for i := 0; i < 100; i++ {
		name = fmt.Sprintf("%s-%s-%s-%s",
			namePart1[rand.Intn(len(namePart1))],
			namePart2[rand.Intn(len(namePart2))],
			namePart3[rand.Intn(len(namePart3))],
			namePart4[rand.Intn(len(namePart4))],
		)
		if _, ok := g.sessions[name]; !ok {
			g.sessions[name] = NewSession(size)
			added = true
			break
		}
	}

	if !added {
		panic("hey we need to fix adding sessions")
	}

	return name
}

func (g *Game) GetSession(name string) *Session {
	return g.sessions[name]
}

func (g *Game) Players() []*Player {
	var players []*Player
	for _, sess := range g.sessions {
		players = append(players, sess.Players()...)
	}

	return players
}

func (g *Game) String() string {
	var buf strings.Builder
	for name, sess := range g.sessions {
		buf.WriteString("\n")
		buf.WriteString("Name: " + name)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("Size: %d", sess.FieldSize()))
		buf.WriteString("\n")
		buf.WriteString(sess.String())
		buf.WriteString("\n")
	}
	buf.WriteString(fmt.Sprintf("Total: %d", len(g.sessions)))

	return buf.String()
}
