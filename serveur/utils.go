package main

import (
	"log"
	"net"
)

// getLocalIP récupère l'adresse IP locale de la machine.
// Elle établit une connexion UDP à une adresse publique (8.8.8.8) sans envoyer de données,
// ce qui permet de déterminer l'adresse IP locale utilisée pour atteindre Internet.
// Si une erreur se produit, elle retourne "inconnue".
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println("Erreur lors de la récupération de l'adresse IP :", err)
		return "inconnue"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
