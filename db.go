package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func connectToBD() *sql.DB {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}

	fSchema, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}
	log.Println(string(fSchema))

	db.Exec(string(fSchema))

	return db
}

func dbAddPlayer(db *sql.DB, player *Player) error {
	insertSQL := "INSERT INTO players (id, username, avatar, sex, email) VALUES ($1, $2, $3, $4, $5)"
	res, err := db.Exec(insertSQL,
		player.ID,
		player.Username,
		player.Avatar,
		player.Sex,
		player.Email,
	)
	if err != nil {
		return err
	}

	log.Println(res.LastInsertId())
	log.Println(res.RowsAffected())

	return nil
}

func dbGetOnePlayer(db *sql.DB, id string) (*Player, error) {
	selectSQL := "SELECT * FROM players WHERE id = $1"
	row := db.QueryRow(selectSQL, id)

	player := Player{}
	err := row.Scan(
		&player.ID,
		&player.Username,
		&player.Avatar,
		&player.Sex,
		&player.Email,
	)

	return &player, err
}

func dbGetAllPlayers(db *sql.DB) (*[]Player, error) {
	selectSQL := "SELECT * FROM players"
	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []Player{}
	for rows.Next() {
		p := Player{}
		err = rows.Scan(
			&p.ID,
			&p.Username,
			&p.Avatar,
			&p.Sex,
			&p.Email,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(p)
		players = append(players, p)
	}

	return &players, nil
}
