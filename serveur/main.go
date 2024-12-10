package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	// Définir un port par défaut
	defaultPort := "8080"

	// Demander à l'utilisateur un port via le terminal
	fmt.Printf("Entrez le port du serveur [par defaut: %s] : ", defaultPort)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Supprimer les espaces ou sauts de ligne

	// Utiliser le port par défaut si aucun input n'est fourni
	port := defaultPort
	if input != "" {
		if _, err := strconv.Atoi(input); err == nil {
			port = input
		} else {
			log.Fatalf("Erreur : '%s' n'est pas un port valide.", input)
		}
	}

	// Obtenir l'adresse IP locale
	localIP := getLocalIP()

	// Écoute sur le port spécifié
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Erreur lors de l'écoute :", err)
	}
	defer listener.Close()

	// Afficher l'IP et le port du serveur
	log.Printf("Serveur démarré. Adresse : %s:%s\n", localIP, port)
	log.Println("En attente de connexions...")

	startServer(listener)
}
