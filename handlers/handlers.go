package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bubo-py/McK/events"
	"github.com/go-chi/chi"
)

var db = events.InitDatabase()

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(db.Storage)
	if err != nil {
		log.Println(err)
	}
}

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	ok, index := db.CheckEvent(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode("Event with specified ID not found")
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(db.Storage[index])
}

func AddEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e events.Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	db.AppendEvent(e)

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println(err)
	}
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	ok := db.DeleteEvent(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode("Event with specified ID not found")
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	var e events.Event
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	ok := db.UpdateEvent(e, id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode("Event with specified ID not found")
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode("Event updated")
	if err != nil {
		log.Println(err)
	}
}
