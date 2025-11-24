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