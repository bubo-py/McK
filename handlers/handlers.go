package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bubo-py/McK/service"
	"github.com/bubo-py/McK/types"
	"github.com/go-chi/chi"
)

var bl = service.BusinessLogic{}

type Handler struct{}

func (h Handler) GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var f types.Filters
	query := r.URL.Query()

	f.Day, _ = strconv.Atoi(query.Get("day"))
	f.Month, _ = strconv.Atoi(query.Get("month"))
	f.Year, _ = strconv.Atoi(query.Get("year"))

	err := json.NewEncoder(w).Encode(bl.GetEvents(f))
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	event, err := bl.GetEvent(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) AddEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e types.Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	err = bl.AddEvent(e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	err = bl.DeleteEvent(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Println(err)
	}

	var e types.Event
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		log.Println(err)
	}

	err = bl.UpdateEvent(e, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(err.Error())
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
