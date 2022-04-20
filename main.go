package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type player struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
	Avatar   string `json:"Avatar"`
	Sex      string `json:"Sex"`
	Email    string `json:"Email"`
}

var players = []player{
	{
		ID:       "773ab1d5-c36c-49dd-b45f-5a26d104601b",
		Username: "kek",
		Avatar:   "",
		Sex:      "male",
		Email:    "kek@gmail.com",
	},
	{
		ID:       "e33ab9c0-f7b3-475d-a3ca-ee9b4231b1ad",
		Username: "lol",
		Avatar:   "",
		Sex:      "female",
		Email:    "lol@mail.ru",
	},
}

func registerPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Body)
	var newPlayer player
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Some shit happens while reading request body")
	}

	json.Unmarshal(reqBody, &newPlayer)
	newPlayer.ID = uuid.New().String()

	players = append(players, newPlayer)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newPlayer)
}

func getAllPlayers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(players)
}

func getOnePlayer(w http.ResponseWriter, r *http.Request) {
	playerID := mux.Vars(r)["id"]

	for _, player := range players {
		if player.ID == playerID {
			json.NewEncoder(w).Encode(player)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/players", registerPlayer).Methods("POST")
	router.HandleFunc("/players", getAllPlayers).Methods("GET")
	router.HandleFunc("/players/{id}", getOnePlayer).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
