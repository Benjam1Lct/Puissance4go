### README pour le Client

---

## Client Puissance 4

Le client est l’interface utilisateur permettant de jouer au Puissance 4 en réseau. Il communique avec le serveur pour synchroniser les actions des joueurs.

### Fonctionnalités

- **Connexion** :
    - Le joueur entre l’adresse du serveur pour se connecter.
    - Vérification de l’état des connexions.

- **Choix des Couleurs** :
    - Navigation via les flèches pour sélectionner une couleur.
    - Validation avec la touche Entrée.

- **Partie** :
    - Contrôle des pions avec les flèches gauche et droite.
    - Placement des pions avec la touche Entrée.

- **Résultats et Redémarrage** :
    - Résultats affichés en fin de partie.
    - Synchronisation avec l’autre joueur pour redémarrer.

### Installation

1. Lancer le client : à la racine du répertoire /client
   ```bash
   go run .
   ```

### Améliorations Possible

- Prévisualisation des choix de l’adversaire lors de la sélection des couleurs.
