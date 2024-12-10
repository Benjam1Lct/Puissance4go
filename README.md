## Puissance 4 en Réseau

Ce projet implémente une version réseau du jeu Puissance 4. Deux joueurs peuvent s'affronter en temps réel via un serveur. Le projet utilise la bibliothèque **Ebitengine** pour les graphismes et propose plusieurs fonctionnalités comme le choix de couleurs, la synchronisation des actions via le réseau, et des améliorations supplémentaires pour une expérience fluide.

### Fonctionnalités Principales

1. **Connexion Réseau** :
    - Les joueurs se connectent via un serveur central.
    - Synchronisation des étapes du jeu : connexion, choix des couleurs, déroulement de la partie, affichage des résultats.

2. **Choix des Couleurs** :
    - Chaque joueur choisit une couleur pour ses pions.
    - Possibilité d’éviter les choix de couleurs identiques entre les joueurs.

3. **Déroulement de la Partie** :
    - Contrôle des pions en temps réel avec synchronisation réseau.
    - Affichage des actions de l'adversaire en direct.

4. **Résultats** :
    - Affichage des résultats après chaque partie.
    - Redémarrage possible après synchronisation des joueurs.

5. **Extensions (Optionnelles)** :
    - Visualisation des choix de couleurs de l’adversaire.
    - Gestion des communications réseau via des goroutines.

---

### Pré-requis

- **Système** : Ubuntu ou une autre distribution Linux.
- **Dépendances** :
    - `libgl1-mesa-dev`
    - `xorg-dev`
- **Langage** : Go
- **Bibliothèque** : Ebitengine

---

### Installation

1. Cloner le projet depuis Gitlab de l’Université.
2. Installer les dépendances :
   ```bash
   sudo apt install libgl1-mesa-dev xorg-dev
   ```
3. Construire le projet :
   ```bash
   go build
   ```
4. Lancer le jeu :
   ```bash
   ./nom_du_binaire
   ```

---

### Contributeurs

- Nom : Loïc Jezequel
- Contact : loig.jezequel@univ-nantes.fr

- Nom : Benjamin LECOMTE
- Contact : benjamin.lecomte@etu.univ-nantes.fr

- Nom : Cheikh KEITA
- Contact : cheikh.keita@etu.univ-nantes.fr