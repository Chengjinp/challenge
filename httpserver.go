// Author: Harsh Nayyar
// Copyright 2017 - PocketHealth

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

// todo: need read from environment
var listeningPort = "8080"
var sqliteDatabasename = "sqlite-challenge.db"

// SubscribeConfirmPage type of subscribe info
type SubscribeConfirmPage struct {
	ID              string `json:"Id"`
	Name            string `json:"Name" validate:"required,name"`
	Email           string `json:"Email" validate:"required,email"`
	Tel             string `json:"Tel" validate:"required,tel"`
	FavouriteColour string `json:"FavouriteColour"`
}

//handleSubscribe sanitize the user input
func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmpl/subscribe.html")
		t.Execute(w, "")
	} else {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

//SanitizeSubscribe sanitize the user input
func SanitizeSubscribe(subscribeInfo SubscribeConfirmPage) bool {
	if len(subscribeInfo.Name) < 1 && len(subscribeInfo.Name) > 255 {
		return false
	}
	if len(subscribeInfo.Email) < 3 && len(subscribeInfo.Email) > 255 {
		return false
	}

	if len(subscribeInfo.Tel) < 10 && len(subscribeInfo.Tel) > 15 {
		return false
	}

	if len(subscribeInfo.FavouriteColour) > 30 {
		return false
	}

	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if emailRegex.MatchString(subscribeInfo.Email) != true {
		return false
	}

	return true
}

// handleSubscribeConfirm sanitize subscribe and save it to sqlite
func handleSubscribeConfirm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleSubscribeConfirm")
	if r.Method == "POST" {
		subscribeConfirmPage := SubscribeConfirmPage{
			Name:            r.FormValue("name"),
			Email:           r.FormValue("email"),
			Tel:             r.FormValue("tel"),
			FavouriteColour: r.FormValue("favouriteColour"),
		}

		if SanitizeSubscribe(subscribeConfirmPage) != true {
			http.Redirect(w, r, "/subscribe", http.StatusSeeOther)
			return
		}

		t, _ := template.ParseFiles("tmpl/subscribeconfirm.html")

		// INSERT RECORDS
		insertSubscribe(sqliteDatabasename, subscribeConfirmPage)

		t.Execute(w, subscribeConfirmPage)
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// handleDefault home page
func handleDefault(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/default.html")
	t.Execute(w, "")
}

// handleSubscribeList return list of subscribe
func handleSubscribeList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := displaySubscribe(sqliteDatabasename)
	json.NewEncoder(w).Encode(data)
}

func main() {
	log.Println("Server is starting ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleDefault)
	mux.HandleFunc("/subscribe", handleSubscribe)
	mux.HandleFunc("/subscribeconfirm", handleSubscribeConfirm)
	mux.HandleFunc("/v1/internal/subscriptions/list", handleSubscribeList)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// create database is not exists
	log.Println("Create DB if not exists ...")
	createDB(sqliteDatabasename)

	// specify listening port
	log.Println("Lisening port : " + listeningPort)
	server := &http.Server{Addr: ":" + listeningPort, Handler: mux}
	// start the server, listening to listeningPort
	server.ListenAndServe()
}
