package main

const DefaultPort = ":8080"

type Message struct {
	Type    string      `json:"type"`    // Le type de message, par ex : "move", "color", etc.
	Payload interface{} `json:"payload"` // Les données spécifiques au message
}

type MovePayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ColorPayload struct {
	Color int `json:"color"`
}
