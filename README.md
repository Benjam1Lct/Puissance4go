# Puissance 4 en Réseau

Ce projet implémente une version réseau du jeu Puissance 4, permettant à deux joueurs de s'affronter en temps réel via un serveur central. Le projet utilise la bibliothèque graphique **Ebitengine** pour l'interface utilisateur et propose une gestion robuste des échanges réseau pour synchroniser les actions des joueurs.

---

## Contributeurs

| Nom                  | Rôle                | Contact                          |
|----------------------|---------------------|----------------------------------|
| **Loïc Jezequel**    | Développeur | loig.jezequel@univ-nantes.fr    |
| **Benjamin LECOMTE** | Développeur  | benjamin.lecomte@etu.univ-nantes.fr |
| **Cheikh KEITA**     | Développeur  | cheikh.keita@etu.univ-nantes.fr |


## Fonctionnalités Principales

### 1. **Connexion Réseau**
- Les joueurs se connectent via un serveur central qui gère les échanges.
- Le serveur attribue un **ID unique** à chaque joueur pour identifier leurs actions.
- Synchronisation des étapes clés :
   - Connexion des joueurs.
   - Prêt pour commencer le jeu.
   - Redémarrage des parties.

### 2. **Choix des Couleurs**
- Chaque joueur sélectionne une couleur pour ses pions.
- Le serveur empêche les deux joueurs de choisir la même couleur.

### 3. **Déroulement de la Partie**
- Les joueurs placent leurs pions en temps réel, synchronisés via le serveur.
- Le serveur transmet les positions des pions à l'adversaire après chaque coup.
- Gestion des actions du jeu avec des **messages structurés** (JSON).

### 4. **Résultats et Redémarrage**
- Une fois la partie terminée, le serveur affiche les résultats aux deux joueurs.
- Possibilité de redémarrer une nouvelle partie après synchronisation des joueurs.

### 5. **Extensions (Optionnelles)**
- Visualisation en temps réel des choix de couleurs de l’adversaire.
- Gestion asynchrone des communications réseau avec des **goroutines**.
- Historique des coups disponible pour chaque joueur en fin de partie.
- Choix du joueur qui commence par pierre/feuille/ciseaux
- Possibilité de voir un replay de la partie
- Affichage du cursor de l'adversaire au-dessus de la grille sur l'écran de la partie
- Refonte totale de l'UI avec affiche d'erreur améliorée
- Plusieurs choix de theme
- Ajout d'un chat en direct in game
- Affichage d'information (nombre joueur connéctés, nombre partie gagné)
- Gestion de la déconnexion impromptu d'un client

---

## Pré-requis

### Dépendances :
- **Langage** : Go (version 1.18 ou supérieure).
- **Bibliothèque graphique** : [Ebitengine](https://ebitengine.org).

---

## Installation

### Étapes :

1. **Cloner le projet** :
   ```bash
   git clone https://gitlab.univ-nantes.fr/pub/but/but2/r3.05/r3.05.projet/r3.05.projet.groupe4.eq07.keita-cheikh_lecomte-benjamin.git
   cd r3.05.projet.groupe4.eq07.keita-cheikh_lecomte-benjamin
   ```

2. **Construire le serveur du projet** :
   ```bash
   cd serveur/
   go build
   ./nom_du_binaire
   ```

2. **Construire le client du projet** :
   ```bash
   cd client/
   go build
   ./nom_du_binaire
   ```

---

## Architecture du Projet

### 1. **Serveur**
- Le serveur est le cœur de la communication entre les deux joueurs.
- Il gère :
   - **Connexions réseau** : accepte et identifie les joueurs.
   - **Synchronisation** : garantit que les deux joueurs progressent simultanément.
   - **Échanges** : transmet les actions des joueurs à l'adversaire.
- Le serveur est implémenté en Go avec des goroutines et des canaux pour gérer plusieurs connexions de manière asynchrone.

### 2. **Client**
- Chaque joueur lance une instance du client pour interagir avec le jeu.
- Le client utilise **Ebitengine** pour afficher l'interface graphique et gérer les interactions.
- Communication réseau via le protocole TCP/IP et des messages JSON structurés.

---

## Protocole de Communication

### 1. **Messages JSON**
Les données échangées entre le serveur et les clients suivent un format structuré pour chaque action du jeu.

Exemple de message JSON pour un déplacement :
```json
{
  "type": "move",
  "payload": {
    "x": 3,
    "y": 2
  }
}
```

- **`type`** : Type d'action (ex. : "move", "color", "ready").
- **`payload`** : Données spécifiques à l'action (coordonnées, couleur, message de chat, etc.).

### 2. **Étapes principales**
- **Connexion** : Le serveur attribue un ID unique à chaque joueur et synchronise leurs connexions.
- **Sélection des couleurs** : Chaque joueur choisit une couleur, et le serveur empêche les doublons.
- **Déroulement du jeu** : Les actions des joueurs (déplacement des pions) sont transmises en temps réel.
- **Fin de partie et rematch** : Une fois la partie terminée, les joueurs peuvent demander un rematch.

---

## Améliorations Possible

1. **Mode spectateur** :
   - Permettre à un troisième utilisateur de regarder les parties en cours.

2. **Chat enrichi** :
   - Ajouter des emojis ou des réactions dans les messages de chat.

3. **Améliorations graphiques** :
   - Ajout d’animations pour les actions des joueurs.