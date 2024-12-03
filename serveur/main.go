package main

import (
	"log"
	"net"
)

func main() {
	// Obtenir l'adresse IP locale
	localIP := getLocalIP()

	// Écoute sur le port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Erreur lors de l'écoute :", err)
	}
	defer listener.Close()

	// Afficher l'IP et le port du serveur
	log.Printf("Serveur démarré. Adresse : %s:8080\n", localIP)
	log.Println("En attente de connexions...")

	startServer(listener)
}
