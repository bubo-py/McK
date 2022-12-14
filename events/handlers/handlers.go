package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/events/service"
	"github.com/bubo-py/McK/types"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var badRequestReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrBadRequest.ErrorType,
	ErrorMessage: customErrors.ErrBadRequest.Error(),
}

var notFoundReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrNotFound.ErrorType,
	ErrorMessage: customErrors.ErrNotFound.Error(),
}

var unexpectedReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrUnexpected.ErrorType,
	ErrorMessage: customErrors.ErrUnexpected.Error(),
}

type Handler struct {
	bl  service.BusinessLogicInterface
	Mux *chi.Mux
}

func InitHandler(bl service.BusinessLogicInterface) Handler {
	var h Handler

	r := chi.NewRouter()

	r.Get("/", h.GetEventsHandler)
	r.Get("/{id}", h.GetEventHandler)
	r.Post("/", h.AddEventHandler)
	r.Put("/{id}", h.UpdateEventHandler)
	r.Delete("/{id}", h.DeleteEventHandler)

	h.Mux = r

	h.bl = bl
	return h
}

func (h *Handler) GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var f types.Filters
	query := r.URL.Query()

	_, present := query["day"]
	if present {
		day, err := strconv.Atoi(query.Get("day"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(badRequestReturn)
			if err != nil {
				log.Println(err)
			}
			return
		}
		f.Day = day
	}

	_, present = query["month"]
	if present {
		month, err := strconv.Atoi(query.Get("month"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(badRequestReturn)
			if err != nil {
				log.Println(err)
			}
			return
		}
		f.Month = month
	}

	_, present = query["year"]
	if present {
		year, err := strconv.Atoi(query.Get("year"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(badRequestReturn)
			if err != nil {
				log.Println(err)
			}
			return
		}
		f.Year = year
	}

	events, err := h.bl.GetEvents(r.Context(), f)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	event, err := h.bl.GetEvent(r.Context(), id)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) AddEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var e types.Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = h.bl.AddEvent(r.Context(), e)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = h.bl.DeleteEvent(r.Context(), id)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	var e types.Event
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = h.bl.UpdateEvent(r.Context(), e, id)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println(err)
	}
}

func errBasedReturn(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, customErrors.ErrBadRequest):
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
	case errors.Is(err, customErrors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(notFoundReturn)
		if err != nil {
			log.Println(err)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(unexpectedReturn)
		if err != nil {
			log.Println(err)
		}
	}
}
