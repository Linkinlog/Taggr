package main

import "encoding/json"

type Responder interface {
	Data() ([]byte, error)
}

type Event struct {
	Data   Responder `json:"data"`
	Action string    `json:"action"`
}

func (e *Event) Response() ([]byte, error) {
	return json.Marshal(e)
}

type PlayerResponse struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	Infected bool   `json:"infected"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

func (p *PlayerResponse) Data() ([]byte, error) {
	return json.Marshal(p)
}

type PlayerMoveResponse struct {
	Name     string `json:"name"`
	Infected bool   `json:"infected"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	NewX     int    `json:"newX"`
	NewY     int    `json:"newY"`
}

func (p *PlayerMoveResponse) Data() ([]byte, error) {
	return json.Marshal(p)
}

type FieldResponse struct {
	Players []PlayerResponse `json:"players"`
	Field   [][]string       `json:"field"`
	Size    int              `json:"size"`
}

func (f *FieldResponse) Data() ([]byte, error) {
	return json.Marshal(f)
}
