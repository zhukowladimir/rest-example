package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	path_to_db     = "../store/store.db"
	path_to_schema = "../store/schema.sql"
)

func connectToBD() *sql.DB {
	db, err := sql.Open("sqlite3", path_to_db)
	if err != nil {
		panic(err)
	}

	fSchema, err := ioutil.ReadFile(path_to_schema)
	if err != nil {
		panic(err)
	}

	db.Exec(string(fSchema))

	return db
}

func dbAddPlayer(db *sql.DB, player *Player) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	cmds := []string{
		"INSERT INTO players (id, username, avatar, sex, email) VALUES (?, ?, ?, ?, ?)",
		"INSERT INTO stats (loss_count, win_count, duration) VALUES (0, 0, 0)",
		"INSERT INTO pid_sid (pid, sid) VALUES (?, ?)",
	}
	var stmts []*sql.Stmt

	for _, cmd := range cmds {
		stmt, err := db.Prepare(cmd)
		if err != nil {
			return err
		}
		stmts = append(stmts, stmt)
	}

	_, err = tx.Stmt(stmts[0]).Exec(
		player.ID,
		player.Username,
		player.Avatar,
		player.Sex,
		player.Email,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := tx.Stmt(stmts[1]).Exec()
	if err != nil {
		tx.Rollback()
		return err
	}
	sid, _ := res.LastInsertId()

	_, err = tx.Stmt(stmts[2]).Exec(
		player.ID,
		sid,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func dbGetOnePlayer(db *sql.DB, id string) (*Player, error) {
	sqlRequest := "SELECT * FROM players WHERE id = $1"
	row := db.QueryRow(sqlRequest, id)

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
	sqlRequest := "SELECT * FROM players"
	rows, err := db.Query(sqlRequest)
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

func dbRemakeOnePlayer(db *sql.DB, p *Player) (*Player, error) {
	sqlRequest := "UPDATE players SET username = $1, avatar = $2, sex = $3, email = $4 WHERE id = $5"

	old, err := dbGetOnePlayer(db, p.ID)
	if err != nil {
		return nil, err
	}

	if p.Avatar != "" {
		old.Avatar = p.Avatar
	}
	if p.Username != "" {
		old.Username = p.Username
	}
	if p.Sex != "" {
		old.Sex = p.Sex
	}
	if p.Avatar != "" {
		old.Email = p.Email
	}

	_, err = db.Exec(sqlRequest,
		old.Username,
		old.Avatar,
		old.Sex,
		old.Email,
		old.ID,
	)
	if err != nil {
		return nil, err
	}

	return old, nil
}

func dbDeleteOnePlayer(db *sql.DB, id string) error {
	sqlRequests := []string{
		"DELETE FROM stats WHERE id = (SELECT id FROM stats INNER JOIN pid_sid ON pid_sid.sid = stats.id WHERE pid_sid.pid = $1)",
		"DELETE FROM pid_sid WHERE pid = $1",
		"DELETE FROM players WHERE id = $1",
	}

	for _, sqlRequest := range sqlRequests {
		_, err := db.Exec(sqlRequest, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func dbUpdatePlayerStats(db *sql.DB, gameStats *GameStats, playerID string) error {
	var sqlRequest string
	if gameStats.IsWin {
		sqlRequest = "UPDATE stats SET win_count = win_count + 1, duration = duration + $1 WHERE id = ( SELECT s.id FROM players AS p INNER JOIN pid_sid AS ps ON p.id = ps.pid INNER JOIN stats AS s ON ps.sid = s.id WHERE p.id = $2 )"
	} else {
		sqlRequest = "UPDATE stats SET loss_count = loss_count + 1, duration = duration + $1 WHERE id = ( SELECT s.id FROM players AS p INNER JOIN pid_sid AS ps ON p.id = ps.pid INNER JOIN stats AS s ON ps.sid = s.id WHERE p.id = $2 )"
	}

	_, err := db.Exec(sqlRequest, gameStats.Duration, playerID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func dbGetPlayerStats(db *sql.DB, playerID string) (*QueueMsg, error) {
	sqlRequest := "SELECT p.id, p.username, p.avatar, p.sex, p.email, s.loss_count, s.win_count, s.duration FROM players AS p INNER JOIN pid_sid AS ps ON ps.pid = p.id INNER JOIN stats AS s ON ps.sid = s.id WHERE ps.pid == $1"

	msg := QueueMsg{
		Filename: strconv.Itoa(time.Now().Nanosecond()) + ".pdf",
	}
	row := db.QueryRow(sqlRequest, playerID)
	err := row.Scan(
		&msg.ID, &msg.Username, &msg.Avatar, &msg.Sex, &msg.Email,
		&msg.LossCount, &msg.WinCount, &msg.Duration,
	)
	if err != nil {
		return nil, err
	}

	log.Print(msg)

	return &msg, nil
}
