package main

// DefaultPort définit le port par défaut utilisé par le serveur.
const DefaultPort = "8080"

// SelectedPayload représente la charge utile pour un message indiquant une sélection d'élément.
type SelectedPayload struct {
	Selected string `json:"selected"` // Élément sélectionné
}

// Message est une structure générique pour échanger des données entre le serveur et les clients.
// Chaque message a un type (Type) et une charge utile (Payload) spécifique à ce type.
type Message struct {
	Type    string      `json:"type"`    // Le type de message, par exemple : "move", "color", "chat", etc.
	Payload interface{} `json:"payload"` // Les données spécifiques au message, de type générique (interface{})
}

// MovePayload représente la charge utile d'un message de type "move".
// Elle contient les coordonnées de déplacement.
type MovePayload struct {
	X int `json:"x"` // Coordonnée X du déplacement
	Y int `json:"y"` // Coordonnée Y du déplacement
}

// ColorPayload représente la charge utile d'un message de type "color".
// Elle contient la couleur sélectionnée par un joueur.
type ColorPayload struct {
	Color int `json:"color"` // Couleur sélectionnée (représentée par un entier)
}

// Coordinate représente une position dans un espace 2D, associée à un joueur (ID).
type Coordinate struct {
	ID int // ID du joueur
	X  int // Coordonnée X
	Y  int // Coordonnée Y
}

// ChatMessage représente la charge utile d'un message de type "chat".
// Elle contient le texte envoyé par un joueur dans le chat.
type ChatMessage struct {
	Text string `json:"text"` // Texte du message envoyé par le joueur
}
