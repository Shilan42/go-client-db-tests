package main

import (
	"database/sql"
)

type Client struct {
	ID       int
	FIO      string
	Login    string
	Birthday string
	Email    string
}

func main() {

}

func selectClient(db *sql.DB, id int) (Client, error) {
	cl := Client{}

	row := db.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&cl.ID, &cl.FIO, &cl.Login, &cl.Birthday, &cl.Email)
	if err != nil {
		return cl, err
	}

	return cl, nil
}

func insertClient(db *sql.DB, client Client) (int, error) {
	res, err := db.Exec("INSERT INTO clients (fio, login, birthday, email) VALUES (:fio, :login, :birthday, :email)",
		sql.Named("fio", client.FIO),
		sql.Named("login", client.Login),
		sql.Named("birthday", client.Birthday),
		sql.Named("email", client.Email))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func deleteClient(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM clients WHERE id = :id", sql.Named("id", id))

	return err
}
