package main

import (
	"database/sql"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func writeMeta(db *sql.DB, key string, value string) bool {
	logger := log.WithFields(log.Fields{
		"key":   key,
		"value": value,
	})

	logger.Trace("Writing meta")
	const errMsg = "Failed to write meta"

	tx, err := db.Begin()
	if err != nil {
		logger.WithFields(log.Fields{"err": err, "step": "begin transaction"}).Fatal(errMsg)
		return false
	}

	stmt, err := tx.Prepare("insert into meta(key, value) values (?, ?)")
	if err != nil {
		logger.WithFields(log.Fields{"err": err, "step": "prepare transaction"}).Fatal(errMsg)
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, value)
	if err != nil {
		logger.WithFields(log.Fields{"err": err, "step": "execute"}).Fatal(errMsg)
		return false
	}

	tx.Commit()
	return true
}

func createTable(db *sql.DB, table string, schema string) bool {
	log.WithFields(log.Fields{
		"table": table,
	}).Trace("Creating table")
	sqlStmt := fmt.Sprintf(`
	create table "%s" (%s)
	`, table, schema)

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.WithFields(log.Fields{
			"table":   table,
			"err":     err,
			"sqlStmt": sqlStmt,
		}).Fatal("Failed to create table")
		return false
	}

	return true
}

func createDB(file string) bool {
	log.WithFields(log.Fields{
		"file": file,
	}).Debug("Creating database")

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.WithFields(log.Fields{
			"file": file,
			"err":  err,
			"step": "open",
		}).Info("Failed to create database")
	}
	defer db.Close()

	createTable(db, "meta", "key text not null primary key, value text not null")
	createTable(db, "tags", "name text not null primary key")
	createTable(db, "sounds", "file text not null primary key")
	createTable(db, "applied_tags", "sound text not null, tag text not null, primary key (sound, tag), foreign key(sound) references sounds(id), foreign key(tag) references tags(name)")

	writeMeta(db, "DatabaseSchemaVersion", strconv.Itoa(DatabaseSchemaVersion))

	return true
}
