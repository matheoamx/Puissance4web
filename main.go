package main

import (
	"encoding/json"
	"fmt"
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
