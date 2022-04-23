package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type Player struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
	Avatar   string `json:"Avatar"`
	Sex      string `json:"Sex"`
	Email    string `json:"Email"`
}

type GameStats struct {
	IsWin    bool `json:"IsWin,string"`
	Duration int  `json:"Duration,string"`
}

type QueueMsg struct {
	ID        string `json:"ID"`
	Username  string `json:"Username"`
	Avatar    string `json:"Avatar"`
	Sex       string `json:"Sex"`
	Email     string `json:"Email"`
	WinCount  int    `json:"WinCount"`
	LossCount int    `json:"LossCount"`
	Duration  int    `json:"Duration"`
	Filename  string `json:"Filename"`
}

var (
	db     *sql.DB
	ch     *amqp.Channel
	wQueue amqp.Queue
	rQueue amqp.Queue
	corrId string
)

const (
	path_to_store = "../store"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func addPlayer(w http.ResponseWriter, r *http.Request) {
	var newPlayer Player
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Some shit happens while reading request body")
	}

	json.Unmarshal(reqBody, &newPlayer)
	newPlayer.ID = uuid.New().String()

	err = dbAddPlayer(db, &newPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newPlayer)
}

func getAllPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := dbGetAllPlayers(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(players)
}

func getOnePlayer(w http.ResponseWriter, r *http.Request) {
	playerID := mux.Vars(r)["id"]

	player, err := dbGetOnePlayer(db, playerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(player)
}

func updatePlayerStats(w http.ResponseWriter, r *http.Request) {
	playerID := mux.Vars(r)["id"]

	var stats GameStats
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Some shit happens while reading request body")
	}

	err = json.Unmarshal(reqBody, &stats)
	if err != nil {
		log.Fatal(err)
	}

	err = dbUpdatePlayerStats(db, &stats, playerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func getPlayerStats(w http.ResponseWriter, r *http.Request) {
	playerID := mux.Vars(r)["id"]

	msg, err := dbGetPlayerStats(db, playerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	err = ch.Publish(
		"",
		wQueue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			CorrelationId: corrId,
			MessageId:     playerID + "/" + msg.Filename,
			ReplyTo:       rQueue.Name,
			Body:          []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(r.Host + "/players/" + playerID + "/stats/" + msg.Filename))
}

func getPdf(w http.ResponseWriter, r *http.Request) {
	playerID := mux.Vars(r)["id"]
	filename := mux.Vars(r)["filename"]

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	fw, err := writer.CreateFormFile("pdf", filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	file, err := os.Open(path_to_store + "/pdfs/" + playerID + "/" + filename)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", writer.FormDataContentType())
	w.Write(body.Bytes())
	w.WriteHeader(http.StatusOK)
}

func main() {
	db = connectToBD()
	defer db.Close()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	corrId = uuid.New().String()

	wQueue, err = ch.QueueDeclare(
		"work_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a wQueue")

	rQueue, err = ch.QueueDeclare(
		"reply_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a wQueue")

	msgs, err := ch.Consume(
		rQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to reguster a consumer")

	go func() {
		for d := range msgs {
			if d.CorrelationId == corrId {
				err = ioutil.WriteFile(path_to_store+"/pdfs/"+d.MessageId, d.Body, 0666)
				if err != nil {
					log.Println("Failed to write pdf: ", err.Error())
				}

			}
		}
	}()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/players", addPlayer).Methods("POST")
	router.HandleFunc("/players", getAllPlayers).Methods("GET")
	router.HandleFunc("/players/{id}", getOnePlayer).Methods("GET")

	router.HandleFunc("/players/{id}/stats", updatePlayerStats).Methods("PUT")
	router.HandleFunc("/players/{id}/stats", getPlayerStats).Methods("GET")
	router.HandleFunc("/players/{id}/stats/{filename}", getPdf).Methods("GET")

	// router.HandleFunc("/auth/sign-up", signUpPlayer).Methods("POST")
	// router.HandleFunc("/auth/login", loginPlayer).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
