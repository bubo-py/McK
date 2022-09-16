package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(db)
	if err != nil {
		log.Println(err)
	}
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	AppendEvent(e)

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println(err)
	}
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	for i, event := range db {
		if event.ID == id {
			copy(db[i:], db[i+1:])
			db[len(db)-1] = Event{}
			db = db[:len(db)-1]
		}
	}

	err = json.NewEncoder(w).Encode("Event deleted")
	if err != nil {
		log.Println(err)
	}
}
