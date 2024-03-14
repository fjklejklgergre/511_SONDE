# Programme de surveillance du système

Ce programme est conçu pour collecter des informations système telles que le nom d'hôte, l'adresse IP, l'utilisation de la RAM, l'utilisation du CPU et l'espace disque, et les stocker dans un fichier JSON.

## Fonctionnalités

- Collecte des informations système telles que le nom d'hôte, l'adresse IP, l'utilisation de la RAM, l'utilisation du CPU et l'espace disque.
- Convertit les informations système au format JSON.
- Met à jour les données JSON dans un fichier spécifié.
- Écoute les connexions entrantes via TCP/IP.

## Utilisation

### Prérequis

- Langage de programmation Go installé.

### Installation

1. Clonez ce dépôt.
2. Accédez au répertoire du projet.
3. Exécutez `go build` pour compiler l'exécutable.

### Exécution

1. Lancez l'exécutable généré après la compilation du programme.
2. Le programme commencera à écouter les connexions entrantes sur le port 8080.

### Remarques

- Assurez-vous d'avoir les autorisations nécessaires pour écrire dans le fichier JSON spécifié (`/var/www/html/received_data.json`).