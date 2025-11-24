package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Structures de données
type Joueur struct {
	ID      int    `json:"id"`
	Pseudo  string `json:"pseudo"`
	Couleur string `json:"couleur"`
}

type Session struct {
	Grille        [][]string
	Joueur1       Joueur
	Joueur2       Joueur
	JoueurActuel  int
	Tour          int
	Gagnant       *Joueur
	Egalite       bool
}

type Partie struct {
	Date     string  `json:"date"`
	Joueur1  Joueur  `json:"joueur1"`
	Joueur2  Joueur  `json:"joueur2"`
	Gagnant  *Joueur `json:"gagnant"`
	Egalite  bool    `json:"egalite"`
	Tour     int     `json:"tour"`
}

var session Session

func main() {
	// Routes statiques
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes de l'application
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/game/init", handleInit)
	http.HandleFunc("/game/play", handlePlay)
	http.HandleFunc("/game/end", handleEnd)
	http.HandleFunc("/game/scoreboard", handleScoreboard)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
// Page d'accueil
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)
}
// Page d'initialisation
func handleInit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Récupérer les données du formulaire
		joueur1Pseudo := r.FormValue("joueur1")
		joueur2Pseudo := r.FormValue("joueur2")
		couleur1 := r.FormValue("couleur1")

		// Validation basique
		if joueur1Pseudo == "" || joueur2Pseudo == "" {
			data := map[string]string{"Error": "Veuillez renseigner les deux pseudos"}
			tmpl := template.Must(template.ParseFiles("templates/init.html"))
			tmpl.Execute(w, data)
			return
		}

		if joueur1Pseudo == joueur2Pseudo {
			data := map[string]string{"Error": "Les pseudos doivent être différents"}
			tmpl := template.Must(template.ParseFiles("templates/init.html"))
			tmpl.Execute(w, data)
			return
		}

		// Déterminer la couleur du joueur 2
		couleur2 := "jaune"
		if couleur1 == "jaune" {
			couleur2 = "rouge"
		}
		// Initialiser la session
		session = Session{
			Grille:       creerGrilleVide(),
			Joueur1:      Joueur{ID: 1, Pseudo: joueur1Pseudo, Couleur: couleur1},
			Joueur2:      Joueur{ID: 2, Pseudo: joueur2Pseudo, Couleur: couleur2},
			JoueurActuel: 1,
			Tour:         1,
			Gagnant:      nil,
			Egalite:      false,
		}

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
		return
	}
	// Afficher le formulaire
	tmpl := template.Must(template.ParseFiles("templates/init.html"))
	tmpl.Execute(w, nil)
}
// Page de jeu
func handlePlay(w http.ResponseWriter, r *http.Request) {
	// Vérifier qu'une session existe
	if session.Joueur1.Pseudo == "" {
		http.Redirect(w, r, "/game/init", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		colonneStr := r.FormValue("colonne")
		colonne, err := strconv.Atoi(colonneStr)

		if err != nil || colonne < 1 || colonne > 7 {
			data := map[string]interface{}{
				"Grille":       session.Grille,
				"Joueur1":      session.Joueur1,
				"Joueur2":      session.Joueur2,
				"JoueurActuel": session.JoueurActuel,
				"Tour":         session.Tour,
				"Erreur":       "Veuillez entrer un numéro de colonne entre 1 et 7",
			}
			tmpl := template.Must(template.ParseFiles("templates/play.html"))
			tmpl.Execute(w, data)
			return
		}
		// Convertir en index (0-6)
		colonne--

		// Placer le jeton
		couleur := session.Joueur1.Couleur
		if session.JoueurActuel == 2 {
			couleur = session.Joueur2.Couleur
		}

		ligne := placerJeton(session.Grille, colonne, couleur)
		if ligne == -1 {
			data := map[string]interface{}{
				"Grille":       session.Grille,
				"Joueur1":      session.Joueur1,
				"Joueur2":      session.Joueur2,
				"JoueurActuel": session.JoueurActuel,
				"Tour":         session.Tour,
				"Erreur":       "Cette colonne est pleine, choisissez-en une autre",
			}
			tmpl := template.Must(template.ParseFiles("templates/play.html"))
			tmpl.Execute(w, data)
			return
		}
		// Vérifier victoire
		if verifierVictoire(session.Grille, ligne, colonne, couleur) {
			if session.JoueurActuel == 1 {
				session.Gagnant = &session.Joueur1
			} else {
				session.Gagnant = &session.Joueur2
			}
			sauvegarderPartie()
			http.Redirect(w, r, "/game/end", http.StatusSeeOther)
			return
		}

		// Vérifier égalité
		if grilleComplete(session.Grille) {
			session.Egalite = true
			sauvegarderPartie()
			http.Redirect(w, r, "/game/end", http.StatusSeeOther)
			return
		}
		// Changer de joueur
		if session.JoueurActuel == 1 {
			session.JoueurActuel = 2
		} else {
			session.JoueurActuel = 1
		}
		session.Tour++
	}
	// Afficher la grille
	data := map[string]interface{}{
		"Grille":       session.Grille,
		"Joueur1":      session.Joueur1,
		"Joueur2":      session.Joueur2,
		"JoueurActuel": session.JoueurActuel,
		"Tour":         session.Tour,
	}
	tmpl := template.Must(template.ParseFiles("templates/play.html"))
	tmpl.Execute(w, data)
}
