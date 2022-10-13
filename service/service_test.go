package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bubo-py/McK/repositories/memoryStorage"
	"github.com/bubo-py/McK/types"
)

var ctx context.Context

func TestGetEvents(t *testing.T) {
	ti := time.Date(2015, 5, 15, 20, 30, 0, 0, time.Local)
	ti2 := time.Date(2017, 9, 15, 20, 30, 0, 0, time.Local)
	ti3 := time.Date(2019, 12, 16, 20, 30, 0, 0, time.Local)

	testCases := []struct {
		filters   types.Filters
		expOutput []types.Event
		expError  error
	}{
		{
			filters:   types.Filters{Day: 32, Month: 5, Year: 2020},
			expOutput: nil,
			expError:  errors.New("invalid day value"),
		},
		{
			filters:   types.Filters{Day: -31, Month: 5, Year: 2020},
			expOutput: nil,
			expError:  errors.New("invalid day value"),
		},
		{
			filters: types.Filters{Day: 15, Month: 10, Year: 2020},
			expOutput: []types.Event{{ID: 1, Name: "Daily meeting", StartTime: ti, EndTime: ti},
				{ID: 2, Name: "Weekly meeting", StartTime: ti2, EndTime: ti2}},
			expError: nil,
		},
		{
			filters:   types.Filters{Day: 10, Month: 15, Year: 2020},
			expOutput: nil,
			expError:  errors.New("invalid month value"),
		},
		{
			filters:   types.Filters{Day: 10, Month: -10, Year: 2020},
			expOutput: nil,
			expError:  errors.New("invalid month value"),
		},
		{
			filters:   types.Filters{Day: 30, Month: 12, Year: 2020},
			expOutput: []types.Event{{ID: 3, Name: "Yearly meeting", StartTime: ti3, EndTime: ti3}},
			expError:  nil,
		},
		{
			filters:  types.Filters{Day: 10, Month: 10, Year: -20},
			expError: errors.New("invalid year value"),
		},
		{
			filters:   types.Filters{Day: 4, Month: 10, Year: 2019},
			expOutput: []types.Event{{ID: 3, Name: "Yearly meeting", StartTime: ti3, EndTime: ti3}},
			expError:  nil,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i)
		t.Run(testName, func(t *testing.T) {
			db := memoryStorage.InitDatabase()
			bl := InitBusinessLogic(db)

			event := types.Event{
				ID:        300,
				Name:      "Daily meeting",
				StartTime: ti,
				EndTime:   ti,
			}
			event2 := types.Event{
				ID:        400,
				Name:      "Weekly meeting",
				StartTime: ti2,
				EndTime:   ti2,
			}
			event3 := types.Event{
				ID:        600,
				Name:      "Yearly meeting",
				StartTime: ti3,
				EndTime:   ti3,
			}

			err := bl.AddEvent(ctx, event)
			if err != nil {
				t.Error(err)
			}

			err = bl.AddEvent(ctx, event2)
			if err != nil {
				t.Error(err)
			}

			err = bl.AddEvent(ctx, event3)
			if err != nil {
				t.Error(err)
			}

			output, err := bl.GetEvents(ctx, tc.filters)
			if err != nil && tc.expError != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Failed to fetch events: got error: %v, expected: %v", err, tc.expError)
				}
			}

			if reflect.DeepEqual(output, tc.expOutput) == false {
				t.Errorf("Failed to fetch events: got output: %v, expected: %v", output, tc.expOutput)
			}

		})
	}
}

func TestAddEvent(t *testing.T) {
	db := memoryStorage.InitDatabase()
	bl := InitBusinessLogic(db)
	ti := time.Date(2020, 5, 15, 20, 30, 0, 0, time.Local)

	event := types.Event{
		ID:        300,
		Name:      "Daily meeting",
		StartTime: ti,
		EndTime:   ti,
	}

	event2 := types.Event{
		ID:        400,
		Name:      "",
		StartTime: ti,
		EndTime:   ti,
	}

	err := bl.AddEvent(ctx, event)
	if err != nil {
		t.Error(err)
	}

	addedEvent, err := bl.GetEvent(ctx, 1)
	if err != nil {
		t.Error(err)
	}
	if addedEvent.Name != event.Name {
		t.Errorf("Failed to add an event: expected name: %s, got: %s", event.Name, addedEvent.Name)
	}

	err = bl.AddEvent(ctx, event2)
	if err != nil {
		if err.Error() != "invalid post request" {
			t.Errorf("Failed to correctly validate an event, expected: invalid post request, got: %v",
				err.Error())
		}
	}
}

func TestValidatePostRequest(t *testing.T) {
	ti := time.Date(2020, 5, 15, 20, 30, 0, 0, time.Local)

	requestData := types.Event{
		ID:          300,
		Name:        "Daily meeting",
		StartTime:   ti,
		EndTime:     ti,
		Description: "A daily meeting for backend team",
		AlertTime:   ti,
	}

	invalidRequestData := types.Event{
		ID:          400,
		Name:        "No time meeting",
		Description: "A daily meeting for backend team",
	}

	err := validatePostRequest(requestData)
	if err != nil {
		t.Errorf("Should validate an event, got error instead: %v", err)
	}

	err = validatePostRequest(invalidRequestData)
	if err == nil {
		t.Errorf("Should throw an error, instead got: %v", err)
	}
}

func TestValidateLength(t *testing.T) {
	testCases := []struct {
		text     string
		expError error
	}{
		{text: "ok ok ok", expError: nil},
		{text: "😂", expError: nil},

		//len 288, rune 96
		{text: "嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複" +
			"嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複",
			expError: nil},

		// len 256, rune 64
		{text: "😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂" +
			"😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂😂",
			expError: nil},

		// len 789, rune 263
		{text: "嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略" +
			"》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉" +
			"定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》" +
			"影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略" +
			"》影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略》影印本重複影印本重複嘉定屠城紀略》影印本重複嘉定屠城紀略",
			expError: errors.New("length should be less than 255 characters")},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i)
		t.Run(testName, func(t *testing.T) {
			err := validateLength(tc.text)

			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Should return different error: got: %v, expected: %v", err, tc.expError)
				}
			}
		})
	}
}
