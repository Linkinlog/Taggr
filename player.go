package main

import (
	"crypto/rand"
	"fmt"
)

func NewPlayer(name string) *Player {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return &Player{
		name: name,
		id:   fmt.Sprintf("%X", b),
	}
}

type Player struct {
	id       string
	name     string
	score    int
	x        int
	y        int
	infected bool
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) X() int {
	return p.x
}

func (p *Player) Y() int {
	return p.y
}

func (p *Player) Postition() (x int, y int) {
	return p.X(), p.Y()
}

func (p *Player) SetPosition(x, y int) {
	p.x = x
	p.y = y
}

func (p *Player) Score() int {
	return p.score
}

func (p *Player) AddScore() {
	p.score += 1
}

func (p *Player) DecreaseScore() {
	p.score -= 1
}

func (p *Player) Infected() bool {
	return p.infected
}

func (p *Player) Infect() {
	p.infected = true
}

func (p *Player) Heal() {
	p.infected = false
}

func (p *Player) String() string {
	return p.name
}
