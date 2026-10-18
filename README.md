# 🎬 CinéCam Studio 🎥

Interface Windows de création de mouvements de caméra pour **Minecraft Bedrock**, avec un thème noir et lavande.

## Fonctionnalités

- Capture d'une position en jeu ou saisie manuelle des coordonnées et angles.
- Organisation, sauvegarde et rechargement des séquences dans le navigateur.
- Réglages de durée, interpolation, boucle et affichage du HUD.
- Lecture et arrêt de la caméra via un pont HTTP compatible CinéCam.
- Interface, polices et ressources intégrées dans un exécutable Go.

## Compiler sous Windows

Prérequis : Go 1.24 ou supérieur et Microsoft Edge ou Google Chrome.

```powershell
go build -trimpath -ldflags "-s -w -H=windowsgui" -o CineCamStudio.exe .
```

Pour garder les messages de diagnostic dans le terminal :

```powershell
go build -o CineCamStudio.exe .
```

## Configurer votre serveur

Cette application est un **client**. Elle nécessite un pont HTTP CinéCam et le comportement de caméra correspondant sur un serveur Minecraft Bedrock. Le pont et le pack serveur ne sont pas inclus dans ce dépôt.

1. Copiez `cinecam.example.json` vers **`cinecam.local.json`**.
2. Remplacez l'adresse d'exemple par celle de votre propre pont CinéCam.
3. Placez ce fichier à côté de `CineCamStudio.exe`.
4. Lancez l'exécutable, puis choisissez votre joueur.

Vous pouvez aussi définir la variable d'environnement `CINECAM_API_BASE`, qui prend priorité sur le fichier :

```powershell
$env:CINECAM_API_BASE = "https://example.com/votre-pont-cinecam/"
.\CineCamStudio.exe
```

Le pont doit accepter les routes `players`, `where`, `tool`, `play` et `stop`, avec le même contrat JSON que l'interface dans `web/index.html`. Il n'y a pas de serveur de démonstration partagé.

Les adresses privées et clés d'accès restent dans votre configuration locale, exclue de Git. Aucun ancien binaire contenant une configuration privée n'est distribué.

## Utilisation

Le programme sert son interface uniquement sur `http://127.0.0.1:47931/` et l'ouvre dans une fenêtre Edge ou Chrome. Si le studio tourne déjà, il tente d'ouvrir sa fenêtre. Après fermeture de l'interface, le programme s'arrête après environ 90 secondes sans activité.

Sans configuration de pont valide, l'interface reste accessible mais les commandes de jeu sont indisponibles. Les séquences sont stockées dans le profil du navigateur utilisé ; exportez-les avant de supprimer ce profil.

## Structure

- `main.go` : serveur local et relais vers votre pont.
- `config.go` : chargement de la configuration locale.
- `sys_windows.go` : ouverture de la fenêtre Windows.
- `web/` : interface, styles, polices et icônes.

## Vérifications

```powershell
go test ./...
go vet ./...
```

## Licence et ressources

Aucune licence de redistribution du code n'a encore été choisie. Les textures, icônes et polices conservent les droits et conditions de leurs auteurs. Projet indépendant, non affilié à Mojang ou Microsoft.
