package handlers

import (
	"io"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bubo-py/McK/events/repositories/mocks"
	"github.com/bubo-py/McK/types"
	"github.com/golang/mock/gomock"
)

func TestGetEventsHandler(t *testing.T) {
	var f types.Filters
	var eventsRet []types.Event

	ti := time.Date(2020, 5, 15, 20, 30, 0, 0, time.Local)

	event := types.Event{
		ID:        300,
		Name:      "Daily meeting",
		StartTime: ti,
		EndTime:   ti,
	}

	eventsRet = append(eventsRet, event)

	r := httptest.NewRequest("GET", "/api/events", nil)
	w := httptest.NewRecorder()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
	mockBL.EXPECT().GetEvents(r.Context(), f).Return(eventsRet, nil).Times(1)

	handler := InitHandler(mockBL)

	handler.GetEventsHandler(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(data))

}
