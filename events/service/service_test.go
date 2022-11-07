package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bubo-py/McK/contextHelpers"
	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/events/repositories/memoryStorage"
	"github.com/bubo-py/McK/events/repositories/mocks"
	"github.com/bubo-py/McK/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var timezoneFetchErr = fmt.Errorf("%w: failed to fetch timezone from context", customErrors.ErrUnexpected)

func TestGetEvents(t *testing.T) {
	ctx := context.Background()
	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Europe/Warsaw")

	ti := time.Date(2015, 5, 15, 20, 30, 0, 0, time.UTC)

	loc, _ := time.LoadLocation("Europe/Warsaw")
	tiWithTimezone := time.Date(2015, 5, 15, 22, 30, 0, 0, loc)

	ti2 := time.Date(2017, 9, 15, 15, 30, 0, 0, time.UTC)

	loc2, _ := time.LoadLocation("Europe/Warsaw")
	ti2WithTimezone := time.Date(2017, 9, 15, 17, 30, 0, 0, loc2)

	ti3 := time.Date(2019, 12, 16, 10, 30, 0, 0, time.UTC)

	loc3, _ := time.LoadLocation("Europe/Warsaw")
	ti3WithTimezone := time.Date(2019, 12, 16, 11, 30, 0, 0, loc3)

	testCases := []struct {
		testName   string
		filters    types.Filters
		noFilters  bool
		mockOutput []types.Event
		mockError  error
		expOutput  []types.Event
		expError   error
	}{
		{
			testName:  "GetEventsNoFilters",
			noFilters: true,
			mockOutput: []types.Event{{ID: 1, Name: "Daily meeting", StartTime: ti, EndTime: ti},
				{ID: 2, Name: "Weekly meeting", StartTime: ti2, EndTime: ti2}},
			expOutput: []types.Event{{ID: 1, Name: "Daily meeting", StartTime: tiWithTimezone, EndTime: tiWithTimezone},
				{ID: 2, Name: "Weekly meeting", StartTime: ti2WithTimezone, EndTime: ti2WithTimezone}},
		},
		{
			testName:  "GetEventsNoFiltersUnexpected",
			noFilters: true,
			mockError: customErrors.ErrUnexpected,
			expError:  customErrors.ErrUnexpected,
		},
		{
			testName: "GetEventsWithFiltersBadRequest_Day1",
			filters:  types.Filters{Day: 32, Month: 5, Year: 2020},
			expError: fmt.Errorf("%w: day", customErrors.ErrBadRequest),
		},
		{
			testName: "GetEventsWithFiltersBadRequest_Day2",
			filters:  types.Filters{Day: -31, Month: 5, Year: 2020},
			expError: fmt.Errorf("%w: day", customErrors.ErrBadRequest),
		},
		{
			testName: "GetEventsWithFilters1",
			filters:  types.Filters{Day: 15, Month: 10, Year: 2020},
			mockOutput: []types.Event{{ID: 1, Name: "Daily meeting", StartTime: ti, EndTime: ti},
				{ID: 2, Name: "Weekly meeting", StartTime: ti2, EndTime: ti2}},
			expOutput: []types.Event{{ID: 1, Name: "Daily meeting", StartTime: tiWithTimezone, EndTime: tiWithTimezone},
				{ID: 2, Name: "Weekly meeting", StartTime: ti2WithTimezone, EndTime: ti2WithTimezone}},
		},
		{
			testName: "GetEventsWithFiltersBadRequest_Month1",
			filters:  types.Filters{Day: 10, Month: 15, Year: 2020},
			expError: fmt.Errorf("%w: month", customErrors.ErrBadRequest),
		},
		{
			testName: "GetEventsWithFiltersBadRequest_Month2",
			filters:  types.Filters{Day: 10, Month: -10, Year: 2020},
			expError: fmt.Errorf("%w: month", customErrors.ErrBadRequest),
		},
		{
			testName:   "GetEventsWithFilters2",
			filters:    types.Filters{Day: 30, Month: 12, Year: 2020},
			mockOutput: []types.Event{{ID: 3, Name: "Yearly meeting", StartTime: ti3, EndTime: ti3}},
			expOutput:  []types.Event{{ID: 3, Name: "Yearly meeting", StartTime: ti3WithTimezone, EndTime: ti3WithTimezone}},
		},
		{
			testName: "GetEventsWithFiltersBadRequest_Year",
			filters:  types.Filters{Day: 10, Month: 10, Year: -20},
			expError: fmt.Errorf("%w: year", customErrors.ErrBadRequest),
		},
		{
			testName:   "GetEventsWithFilters3",
			filters:    types.Filters{Day: 4, Month: 10, Year: 2019},
			mockOutput: []types.Event{{ID: 3, Name: "Yearly meeting", StartTime: ti3, EndTime: ti3}},
			expOutput:  []types.Event{{ID: 3, Name: "Yearly meeting", StartTime: ti3WithTimezone, EndTime: ti3WithTimezone}},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDB := mocks.NewMockDatabaseRepository(mockCtrl)
			bl := InitBusinessLogic(mockDB)

			if tc.noFilters {
				mockDB.EXPECT().GetEvents(ctx).Return(tc.mockOutput, tc.mockError).Times(1)

				e, err := bl.GetEvents(ctx, tc.filters)
				require.Equal(t, tc.expOutput, e, "events should be equal")
				require.Equal(t, tc.expError, err, "errors should be equal")
			} else {
				if tc.expError == nil {
					mockDB.EXPECT().GetEventsFiltered(ctx, tc.filters).Return(tc.mockOutput, tc.mockError).Times(1)
				}

				e, err := bl.GetEvents(ctx, tc.filters)
				require.Equal(t, tc.expOutput, e, "events should be equal")
				require.Equal(t, tc.expError, err, "errors should be equal")
			}
		})
	}
}

func TestGetEvent(t *testing.T) {
	ctx := context.Background()
	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Asia/Tokyo")

	tiUTC := time.Date(2015, 5, 15, 10, 30, 0, 0, time.UTC)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	tiJST := time.Date(2015, 5, 15, 19, 30, 0, 0, loc)

	testCases := []struct {
		testName    string
		callMock    bool
		id          int64
		eventFromDB types.Event
		eventToUser types.Event
		mockError   error
		expError    error
	}{
		{
			testName: "GetEventNoError",
			callMock: true,
			eventFromDB: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
				EndTime:   tiUTC,
				AlertTime: tiJST,
			},
			eventToUser: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiJST,
				EndTime:   tiJST,
				AlertTime: tiJST,
			},
		},
		{
			testName:  "GetEventNotFound",
			callMock:  true,
			mockError: customErrors.ErrNotFound,
			expError:  customErrors.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDB := mocks.NewMockDatabaseRepository(mockCtrl)
			bl := InitBusinessLogic(mockDB)

			if tc.callMock {
				mockDB.EXPECT().GetEvent(ctx, tc.id).Return(tc.eventFromDB, tc.mockError)
			}
			event, err := bl.GetEvent(ctx, tc.id)

			require.Equal(t, event, tc.eventToUser)
			require.Equal(t, err, tc.expError)
		})
	}
}

func TestAddEvent(t *testing.T) {
	ctx := context.Background()
	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Asia/Tokyo")

	tiUTC := time.Date(2015, 5, 15, 10, 30, 0, 0, time.UTC)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	tiJST := time.Date(2015, 5, 15, 19, 30, 0, 0, loc)

	testCases := []struct {
		testName               string
		badRequestPresent      bool
		eventToAdd             types.Event
		eventConvertedTimezone types.Event
		mockError              error
		expError               error
	}{
		{
			testName: "AddEventNoError",
			eventToAdd: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiJST,
				EndTime:   tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
				EndTime:   tiUTC,
			},
		},
		{
			testName:          "AddEventBadRequestNoName",
			badRequestPresent: true,
			eventToAdd: types.Event{
				ID:        3,
				Name:      "",
				StartTime: tiJST,
				EndTime:   tiJST,
			},
			expError: fmt.Errorf("%w: invalid post request", customErrors.ErrBadRequest),
		},
		{
			testName:          "AddEventBadRequestNoStartTime",
			badRequestPresent: true,
			eventToAdd: types.Event{
				ID:      3,
				Name:    "hello",
				EndTime: tiJST,
			},
			expError: fmt.Errorf("%w: invalid post request", customErrors.ErrBadRequest),
		},
		{
			testName:          "AddEventBadRequestNoEndTime",
			badRequestPresent: true,
			eventToAdd: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
			},
			expError: fmt.Errorf("%w: invalid post request", customErrors.ErrBadRequest),
		},
		{
			testName: "AddEventUnexpected",
			eventToAdd: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiJST,
				EndTime:   tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
				EndTime:   tiUTC,
			},
			mockError: customErrors.ErrUnexpected,
			expError:  customErrors.ErrUnexpected,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDB := mocks.NewMockDatabaseRepository(mockCtrl)
			bl := InitBusinessLogic(mockDB)

			if tc.badRequestPresent {
				err := bl.AddEvent(ctx, tc.eventToAdd)

				require.Equal(t, tc.expError, err, "errors should be equal")
			} else {
				mockDB.EXPECT().AddEvent(ctx, tc.eventConvertedTimezone).Return(tc.mockError)

				err := bl.AddEvent(ctx, tc.eventToAdd)
				require.Equal(t, tc.expError, err, "errors should be equal")
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	ctx := context.Background()
	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Asia/Tokyo")

	tiUTC := time.Date(2015, 5, 15, 10, 30, 0, 0, time.UTC)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	tiJST := time.Date(2015, 5, 15, 19, 30, 0, 0, loc)

	testCases := []struct {
		testName               string
		id                     int64
		badRequestPresent      bool
		eventToUpdate          types.Event
		eventConvertedTimezone types.Event
		mockError              error
		expError               error
	}{
		{
			testName: "UpdateEventNoError",
			id:       3,
			eventToUpdate: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiJST,
				EndTime:   tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
				EndTime:   tiUTC,
			},
		},
		{
			testName: "UpdateEventNoName",
			id:       3,
			eventToUpdate: types.Event{
				ID:        3,
				StartTime: tiJST,
				EndTime:   tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:        3,
				StartTime: tiUTC,
				EndTime:   tiUTC,
			},
		},
		{
			testName: "UpdateEventNoStartTime",
			id:       3,
			eventToUpdate: types.Event{
				ID:      3,
				Name:    "hello",
				EndTime: tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:      3,
				Name:    "hello",
				EndTime: tiUTC,
			},
		},
		{
			testName: "UpdateEventNoEndTime",
			id:       3,
			eventToUpdate: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
			},
		},
		{
			testName: "UpdateEventUnexpected",
			id:       3,
			eventToUpdate: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiJST,
				EndTime:   tiJST,
			},
			eventConvertedTimezone: types.Event{
				ID:        3,
				Name:      "hello",
				StartTime: tiUTC,
				EndTime:   tiUTC,
			},
			mockError: customErrors.ErrUnexpected,
			expError:  customErrors.ErrUnexpected,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockDB := mocks.NewMockDatabaseRepository(mockCtrl)
			bl := InitBusinessLogic(mockDB)

			if tc.badRequestPresent {
				err := bl.UpdateEvent(ctx, tc.eventToUpdate, tc.id)

				require.Equal(t, tc.expError, err, "errors should be equal")
			} else {
				mockDB.EXPECT().UpdateEvent(ctx, tc.eventConvertedTimezone, tc.id).Return(tc.mockError)

				err := bl.UpdateEvent(ctx, tc.eventToUpdate, tc.id)
				require.Equal(t, tc.expError, err, "errors should be equal")
			}
		})
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

	badReq := fmt.Errorf("%w: invalid post request", customErrors.ErrBadRequest)
	err := validatePostRequest(invalidRequestData)
	require.Equal(t, badReq, err)

	err = validatePostRequest(requestData)
	require.Nil(t, err)
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
			expError: fmt.Errorf("%w: length should be less than 255 characters", customErrors.ErrBadRequest)},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i)
		t.Run(testName, func(t *testing.T) {
			err := validateLength(tc.text)

			if tc.expError != nil {
				require.Equal(t, tc.expError, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestEventToUserTime(t *testing.T) {
	ctx := context.Background()

	db := memoryStorage.InitDatabase()
	bl := InitBusinessLogic(db)

	// UTC time, in Europe/Warsaw should be 12
	ti := time.Date(2020, 5, 15, 10, 0, 0, 0, time.UTC)

	event := types.Event{
		ID:        300,
		Name:      "Daily meeting",
		StartTime: ti,
		EndTime:   ti,
	}

	_, err := bl.eventToUserTime(ctx, event.StartTime)
	require.NotNilf(t, err, "Should return an error: %v", timezoneFetchErr)
	require.Equal(t, err, timezoneFetchErr)

	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Europe/Warsaw")

	event.StartTime, err = bl.eventToUserTime(ctx, event.StartTime)
	require.Nil(t, err)

	require.Equal(t, event.StartTime.Hour(), 12, "Failed to convert timezone")

}

func TestEventToUTC(t *testing.T) {
	ctx := context.Background()

	db := memoryStorage.InitDatabase()
	bl := InitBusinessLogic(db)

	loc, _ := time.LoadLocation("Europe/Warsaw")

	// Should convert to 10 UTC
	ti := time.Date(2020, 5, 15, 12, 0, 0, 0, loc)
	event := types.Event{
		ID:        300,
		Name:      "Daily meeting",
		StartTime: ti,
		EndTime:   ti,
	}

	_, err := bl.eventToUTC(ctx, event)
	require.NotNilf(t, err, "Should return an error: %v", timezoneFetchErr)
	require.Equal(t, err, timezoneFetchErr)

	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Europe/Warsaw")

	event, err = bl.eventToUTC(ctx, event)
	require.Nil(t, err)

	require.Equal(t, event.StartTime.Hour(), 10, "Failed to convert timezone")
}
