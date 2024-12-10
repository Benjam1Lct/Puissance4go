
### README pour le Serveur

---

## Serveur Puissance 4

Le serveur est le cœur du projet. Il gère la communication entre les clients et assure la synchronisation des étapes du jeu.

### Fonctionnalités

- **Gestion des Connexions** :
    - Accepte deux connexions clients.
    - Synchronise les clients pour démarrer le jeu simultanément.

- **Synchronisation des Phases** :
    - Indique à chaque client quand passer à la phase suivante.
    - Gère les actions de jeu en temps réel.

- **Partage des Données** :
    - Transmet les choix des joueurs (position des pions, résultats).
    - Maintient l’état global du jeu.

### Installation

1. Démarrer le serveur : à la racine du repertoir /serveur
   ```bash
   go run .
   ```
2. Le serveur écoute sur le port spécifié dans la console et attend deux connexions clients.

---

### Protocole

1. **Connexion** :
    - Les clients envoient un message de connexion.
    - Le serveur répond une fois les deux clients connectés.

2. **Choix des Couleurs** :
    - Le serveur reçoit les choix des couleurs et les valide.

3. **Partie** :
    - Le serveur transmet les positions des pions joués à chaque client.

4. **Résultats et Redémarrage** :
    - Le serveur synchronise les clients pour recommencer une partie.

---

### Extensions

- Gestion via goroutines pour chaque client.
- Transmettre en temps réel les sélections de couleurs.
