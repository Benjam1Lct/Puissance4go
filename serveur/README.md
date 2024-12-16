
# Serveur Puissance 4

Le serveur est le cœur du projet. Il gère la communication entre les clients et assure la synchronisation des étapes du jeu Puissance 4, avec une gestion robuste des messages échangés.

---

## Fonctionnalités

### 1. **Gestion des Connexions**
- Le serveur écoute sur un port par défaut (**`:8080`**) ou un autre port spécifié.
- Accepte jusqu’à **deux clients** simultanément.
- Les clients sont identifiés par un ID unique, attribué lors de leur connexion.
- Le serveur synchronise les connexions pour garantir que les deux joueurs soient prêts avant de commencer la partie.

### 2. **Synchronisation des Phases de Jeu**
- **Sélection des Couleurs** : Chaque joueur choisit une couleur, et le serveur notifie les autres joueurs de leurs choix.
- **Déplacement des Pions** : Les positions jouées par un joueur sont transmises en temps réel à l’autre joueur.
- **Prêt pour Redémarrer** : Le serveur gère les signaux de redémarrage envoyés par les joueurs et coordonne la préparation d’une nouvelle partie.

### 3. **Gestion en Temps Réel**
- Utilise des goroutines pour gérer les connexions des clients simultanément.
- Un mutex protège les accès concurrents aux ressources partagées (comme les connexions et l’état du jeu).
- Les données des messages sont sérialisées/désérialisées en JSON pour un échange standardisé.

### 4. **Protocole de Communication**
- Échanges de messages structurés entre le serveur et les clients, avec des types spécifiques :
    - **`ready`** : Indique que le joueur est prêt à jouer.
    - **`move`** : Représente un déplacement d’un pion.
    - **`color`** : Sélection de couleur par un joueur.
    - **`chat`** : Messages texte envoyés par les joueurs.
    - **`restartReady`** : Signal que le joueur est prêt à redémarrer.
    - **`require_history`** : Demande l’historique des actions de la partie.

---

## Installation et Lancement

### Pré-requis
- **Go** version 1.18 ou supérieure.

### Étapes d'installation

1. Clonez le dépôt :
   ```bash
   git clone https://gitlab.univ-nantes.fr/pub/but/but2/r3.05/r3.05.projet/r3.05.projet.groupe4.eq07.keita-cheikh_lecomte-benjamin.git
   cd r3.05.projet.groupe4.eq07.keita-cheikh_lecomte-benjamin/serveur
   ```

2. Lancez le serveur depuis le répertoire racine :
   ```bash
   go run .
   ```

3. Par défaut, le serveur écoute sur le port **`:8080`** (modifiable dans le code via la constante `DefaultPort`).

---

## Protocole de Communication

### 1. **Connexion**
- Les clients se connectent au serveur via une connexion TCP.
- Une fois les deux clients connectés, le serveur attribue un ID unique à chaque joueur et leur notifie leur statut.

### 2. **Échanges Structurés (JSON)**
- **Message Type** : Le champ `Type` du message indique la nature de l’action.
- **Payload** : Charge utile associée au message, contenant les données spécifiques.

Exemple de message JSON envoyé par le serveur :
```json
{
  "type": "move",
  "payload": {
    "x": 2,
    "y": 3
  }
}
```

---

## Structure des Données

### **Structures principales**
- **`Message`** : Structure générique pour tous les messages échangés.
- **`MovePayload`** : Contient les coordonnées d’un déplacement.
- **`ColorPayload`** : Représente la sélection de couleur d’un joueur.
- **`ChatMessage`** : Contient un message texte envoyé par un joueur.
- **`Coordinate`** : Représente la position d’un pion joué par un joueur.

### **Canaux et Synchronisation**
- **`restartReadyChannel`** : Gère les signaux de redémarrage envoyés par les joueurs.
- **`restartControlChannel`** : Utilisé pour arrêter et redémarrer proprement les processus de synchronisation.

---

## Extensions

- **Gestion des Messages** :
    - Les types de messages sont traités par la fonction `processMessage`, qui gère chaque action en fonction de son type.
    - Les actions incluent la mise à jour des positions, la synchronisation des couleurs, et la diffusion de messages de chat.

- **Gestion du Redémarrage** :
    - Le serveur utilise des canaux pour attendre que tous les joueurs soient prêts avant de redémarrer une partie.

- **Robustesse** :
    - Les erreurs de réseau ou de sérialisation JSON sont loguées et gérées pour éviter les plantages du serveur.
    - Les accès concurrents sont protégés grâce au mutex (`sync.Mutex`).

---

## Exemple de Fonctionnalité

### 1. **Récupération de l’IP Locale**
La fonction `getLocalIP` permet au serveur de récupérer son adresse IP locale pour l’afficher ou la partager :
```go
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println("Erreur lors de la récupération de l'adresse IP :", err)
		return "inconnue"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
```

---

## Logs et Débogage

- Les logs sont utilisés pour suivre les événements clés :
    - Connexion/déconnexion des clients.
    - Synchronisation des étapes (choix des couleurs, déplacements, redémarrage).
    - Erreurs de réseau ou de traitement des messages.

---

Avec ce guide, vous avez une vue complète du fonctionnement et des fonctionnalités du serveur Puissance 4, ainsi que les étapes pour le configurer et l’utiliser efficacement. 