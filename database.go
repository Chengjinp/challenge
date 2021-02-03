package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

// createDB create database if not exists
func createDB(dababaseName string) {

	if _, err := os.Stat(dababaseName); os.IsNotExist(err) {
		// SQLite DB is not exists.
		log.Println("Creating " + dababaseName)
		file, err := os.Create(dababaseName) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println(dababaseName + " created")
	}
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+dababaseName) // Open the created SQLite File
	defer sqliteDatabase.Close()                                // Defer Closing the database
	createTable(sqliteDatabase)                                 // Create Database Tables
}

//createTable make sure table exists all the time.
func createTable(db *sql.DB) {
	createSubcribeTableSQL := `CREATE TABLE IF NOT EXISTS subscribe (
			"idSubscribe" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
			"name" NVARCHAR(255),
			"email" NVARCHAR(255),
			"telephone" NVARCHAR(15),
			"favouriteColour" NVARCHAR(30)		
		  );` // SQL Statement for Create Table

	log.Println("Create subscribe table if not exists...")
	statement, err := db.Prepare(createSubcribeTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("subscribe table created")
}

//insertSubscribe passing db reference connection from main to our method with other parameters
func insertSubscribe(databaseName string, subscribeConfirmPage SubscribeConfirmPage) {
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+databaseName) // Open the created SQLite File
	defer sqliteDatabase.Close()                                // Defer Closing the database

	log.Println("Inserting subscribe record ...")
	insertSubscribeSQL := `INSERT INTO subscribe(name, email, telephone, favouriteColour) VALUES (?, ?, ?, ?)`
	statement, err := sqliteDatabase.Prepare(insertSubscribeSQL) // Prepare statement.

	// This is good to avoid SQL injections - # 7
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(subscribeConfirmPage.Name, subscribeConfirmPage.Email, subscribeConfirmPage.Tel, subscribeConfirmPage.FavouriteColour)
	if err != nil {
		log.Fatalln(err.Error())
	}

}

// displaySubscribe get all subscribe records
func displaySubscribe(databaseName string) []SubscribeConfirmPage {
	sqliteDatabase, _ := sql.Open("sqlite3", "./"+databaseName) // Open the SQLite File
	defer sqliteDatabase.Close()
	row, err := sqliteDatabase.Query("SELECT idSubscribe, name, email, telephone, favouriteColour FROM subscribe ORDER BY idSubscribe")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	subscribes := []SubscribeConfirmPage{}

	for row.Next() { // Iterate and fetch the records from result cursor
		subscribe := SubscribeConfirmPage{}
		row.Scan(&subscribe.ID, &subscribe.Name, &subscribe.Email, &subscribe.Tel, &subscribe.FavouriteColour)
		subscribes = append(subscribes, subscribe)
		log.Println("subscribe: ", subscribe.ID, " ", subscribe.Name, " ", subscribe.Email, " ", subscribe.Email, " ", subscribe.FavouriteColour)
	}
	return subscribes
}
