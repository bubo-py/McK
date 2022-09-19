package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bubo-py/McK/events"
	"github.com/go-chi/chi"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(events.Db)
	if err != nil {
		log.Println(err)
	}
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	for i, event := range events.Db {
		if event.ID == id {
			err = json.NewEncoder(w).Encode(events.Db[i])
		}
	}
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e events.Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	events.ID += 1
	e.ID = events.ID

	events.AppendEvent(e)

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

	for i, event := range events.Db {
		if event.ID == id {
			copy(events.Db[i:], events.Db[i+1:])
			events.Db[len(events.Db)-1] = events.Event{}
			events.Db = events.Db[:len(events.Db)-1]
		}
	}

	err = json.NewEncoder(w).Encode("Event deleted")
	if err != nil {
		log.Println(err)
	}
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	var e events.Event
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	for i, event := range events.Db {
		if event.ID == id {
			events.Db[i].Name = e.Name
			events.Db[i].StartTime = e.StartTime
			events.Db[i].EndTime = e.EndTime
			events.Db[i].Description = e.Description
			events.Db[i].AlertTime = e.AlertTime
		}
	}

	err = json.NewEncoder(w).Encode("Event updated")
	if err != nil {
		log.Println(err)
	}
}
