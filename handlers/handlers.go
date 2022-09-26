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

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(bl.GetEvents())
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

func AddEventHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
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

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
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

func GetEventsByDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	day := chi.URLParam(r, "day")

	events, err := bl.GetEventsByDay(day)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}

func GetEventsByMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	month := chi.URLParam(r, "month")

	events, err := bl.GetEventsByMonth(month)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}

func GetEventsByYear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	year := chi.URLParam(r, "year")

	events, err := bl.GetEventsByYear(year)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}
