package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"jenujari/go-sbc-webapp/lib"
	"jenujari/go-sbc-webapp/sqls"

	"github.com/jenujari/go-swe-api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockQuerier is a mock of the sqls.Querier interface
type mockQuerier struct {
	mock.Mock
}

func (m *mockQuerier) CreateOLHC(ctx context.Context, arg sqls.CreateOLHCParams) (sqls.TblOhlc, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqls.TblOhlc), args.Error(1)
}

func (m *mockQuerier) DeleteOLHC(ctx context.Context, arg sqls.DeleteOLHCParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockQuerier) GetOLHC(ctx context.Context, arg sqls.GetOLHCParams) (sqls.TblOhlc, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqls.TblOhlc), args.Error(1)
}

func (m *mockQuerier) ListOLHCs(ctx context.Context) ([]sqls.TblOhlc, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sqls.TblOhlc), args.Error(1)
}

func (m *mockQuerier) ListTickers(ctx context.Context) ([]sqls.TblTicker, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sqls.TblTicker), args.Error(1)
}

func (m *mockQuerier) UpdateOLHC(ctx context.Context, arg sqls.UpdateOLHCParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockQuerier) UpsertOLHC(ctx context.Context, arg sqls.UpsertOLHCParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockQuerier) UpsertPanchang(ctx context.Context, arg sqls.UpsertPanchangParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockQuerier) UpsertPlanetPosition(ctx context.Context, arg sqls.UpsertPlanetPositionParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// mockSweGrpcClient is a mock of the lib.SweGrpcClient interface
type mockSweGrpcClient struct {
	mock.Mock
}

func (m *mockSweGrpcClient) Ping(ctx context.Context) (*proto.PingResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*proto.PingResponse), args.Error(1)
}

func (m *mockSweGrpcClient) GetPos(ctx context.Context, datetime string, planet string) (*proto.PosResponse, error) {
	args := m.Called(ctx, datetime, planet)
	return args.Get(0).(*proto.PosResponse), args.Error(1)
}

func (m *mockSweGrpcClient) Tithy(ctx context.Context, timestamp string) (*proto.TithyResponse, error) {
	args := m.Called(ctx, timestamp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*proto.TithyResponse), args.Error(1)
}

func (m *mockSweGrpcClient) FindConjunction(ctx context.Context, start, end, planet1, planet2 string, orb int32, step float64) (*proto.ConjunctionResponse, error) {
	args := m.Called(ctx, start, end, planet1, planet2, orb, step)
	return args.Get(0).(*proto.ConjunctionResponse), args.Error(1)
}

func (m *mockSweGrpcClient) GetBalas(ctx context.Context, ts string) (*proto.BalasResponse, error) {
	args := m.Called(ctx, ts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*proto.BalasResponse), args.Error(1)
}

func newTestContext(db *lib.DBService, swe lib.SweGrpcClient) context.Context {
	return context.WithValue(context.Background(), "services", map[string]any{
		"webData":   lib.WebData{"appname": "webapp"},
		"db":        db,
		"sweClient": swe,
	})
}

func TestAstrologyGenerateHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ohlc-upload/generate-astrology", nil)
	rr := httptest.NewRecorder()

	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestAstrologyGenerateHandler_MissingParameters(t *testing.T) {
	form := url.Values{}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(newTestContext(&lib.DBService{Queries: new(mockQuerier)}, new(mockSweGrpcClient)))
	rr := httptest.NewRecorder()

	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "All fields (From Date, To Date, and Time) are required")
}

func TestAstrologyGenerateHandler_InvalidDates(t *testing.T) {
	form := url.Values{
		"from_date": {"invalid-date"},
		"to_date":   {"2026-05-28"},
		"time":      {"09:30"},
	}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(newTestContext(&lib.DBService{Queries: new(mockQuerier)}, new(mockSweGrpcClient)))
	rr := httptest.NewRecorder()

	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid From Date format")
}

func TestAstrologyGenerateHandler_ToDateBeforeFromDate(t *testing.T) {
	form := url.Values{
		"from_date": {"2026-05-28"},
		"to_date":   {"2026-05-25"},
		"time":      {"09:30"},
	}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(newTestContext(&lib.DBService{Queries: new(mockQuerier)}, new(mockSweGrpcClient)))
	rr := httptest.NewRecorder()

	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "From Date must be before or equal to To Date")
}

func TestAstrologyGenerateHandler_RangeTooLarge(t *testing.T) {
	// 10 years + 2 days
	form := url.Values{
		"from_date": {"2016-05-25"},
		"to_date":   {"2026-05-27"},
		"time":      {"09:30"},
	}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(newTestContext(&lib.DBService{Queries: new(mockQuerier)}, new(mockSweGrpcClient)))
	rr := httptest.NewRecorder()

	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Date range cannot exceed 10 years")
}

func TestAstrologyGenerateHandler_DatabaseNotAvailable(t *testing.T) {
	form := url.Values{
		"from_date": {"2026-05-25"},
		"to_date":   {"2026-05-25"},
		"time":      {"09:30"},
	}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// db service missing in context
	req = req.WithContext(context.WithValue(req.Context(), "services", map[string]any{
		"webData":   lib.WebData{"appname": "webapp"},
		"sweClient": new(mockSweGrpcClient),
	}))

	rr := httptest.NewRecorder()
	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Database connection is not available")
}

func TestAstrologyGenerateHandler_SuccessPath(t *testing.T) {
	form := url.Values{
		"from_date": {"2026-05-25"},
		"to_date":   {"2026-05-26"}, // 2 days range: May 25, May 26
		"time":      {"09:30"},
	}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockQ := new(mockQuerier)
	mockSwe := new(mockSweGrpcClient)

	dbSvc := &lib.DBService{
		Queries: mockQ,
	}

	req = req.WithContext(newTestContext(dbSvc, mockSwe))

	// Setup expectations for sweClient.Tithy
	// Times must be in UTC! "2026-05-25T09:30:00Z" and "2026-05-26T09:30:00Z"
	t1 := "2026-05-25T09:30:00Z"
	t2 := "2026-05-26T09:30:00Z"

	mockSwe.On("Tithy", mock.Anything, t1).Return(&proto.TithyResponse{
		Tithy:     5,
		Nakshatra: "Magha",
		Weekday:   "Monday",
	}, nil)

	mockSwe.On("Tithy", mock.Anything, t2).Return(&proto.TithyResponse{
		Tithy:     6,
		Nakshatra: "Purva Phalguni",
		Weekday:   "Tuesday",
	}, nil)

	// Setup expectations for sweClient.GetBalas
	mockSwe.On("GetBalas", mock.Anything, t1).Return(&proto.BalasResponse{
		Results: map[string]*proto.PlanetBalas{
			"sun": {
				Cords: &proto.PlanetCord{
					Name:      "sun",
					Sign:      "Leo",
					Nakshatra: &proto.NakshatraPada{Name: "Magha", Pada: 1},
					Longitude: 120.5,
					Latitude:  0.0,
					Distance:  1.0,
					SpeedLong: 0.98,
					IsRetro:   false,
				},
				UdayBala: 90,
			},
		},
	}, nil)

	mockSwe.On("GetBalas", mock.Anything, t2).Return(&proto.BalasResponse{
		Results: map[string]*proto.PlanetBalas{
			"sun": {
				Cords: &proto.PlanetCord{
					Name:      "sun",
					Sign:      "Leo",
					Nakshatra: &proto.NakshatraPada{Name: "Magha", Pada: 2},
					Longitude: 121.5,
					Latitude:  0.0,
					Distance:  1.0,
					SpeedLong: 0.98,
					IsRetro:   false,
				},
				UdayBala: 91,
			},
		},
	}, nil)

	// Setup expectations for DB UpsertPanchang
	mockQ.On("UpsertPanchang", mock.Anything, mock.MatchedBy(func(p sqls.UpsertPanchangParams) bool {
		return p.Tithy == 5 && p.Nakshatra == sqls.NakshatraTypeMagha && p.Weekday == sqls.WeekDayTypeMonday
	})).Return(nil)

	mockQ.On("UpsertPanchang", mock.Anything, mock.MatchedBy(func(p sqls.UpsertPanchangParams) bool {
		return p.Tithy == 6 && p.Nakshatra == sqls.NakshatraTypePurvaPhalguni && p.Weekday == sqls.WeekDayTypeTuesday
	})).Return(nil)

	// Setup expectations for DB UpsertPlanetPosition
	mockQ.On("UpsertPlanetPosition", mock.Anything, mock.MatchedBy(func(p sqls.UpsertPlanetPositionParams) bool {
		return p.PlanetName == sqls.PlanetTypeSun && p.Longitude == 120.5 && p.UdayBala != nil && *p.UdayBala == 90
	})).Return(nil)

	mockQ.On("UpsertPlanetPosition", mock.Anything, mock.MatchedBy(func(p sqls.UpsertPlanetPositionParams) bool {
		return p.PlanetName == sqls.PlanetTypeSun && p.Longitude == 121.5 && p.UdayBala != nil && *p.UdayBala == 91
	})).Return(nil)

	rr := httptest.NewRecorder()
	astrologyGenerateHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.Contains(t, body, "Success")
	assert.Contains(t, body, "Panchang Records")
	assert.Contains(t, body, "Planet Positions")

	mockQ.AssertExpectations(t)
	mockSwe.AssertExpectations(t)
}

func TestAstrologyGenerateHandler_SweClientError(t *testing.T) {
	form := url.Values{
		"from_date": {"2026-05-25"},
		"to_date":   {"2026-05-25"},
		"time":      {"09:30"},
	}
	req := httptest.NewRequest(http.MethodPost, "/ohlc-upload/generate-astrology", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mockQ := new(mockQuerier)
	mockSwe := new(mockSweGrpcClient)

	dbSvc := &lib.DBService{
		Queries: mockQ,
	}

	req = req.WithContext(newTestContext(dbSvc, mockSwe))

	t1 := "2026-05-25T09:30:00Z"

	mockSwe.On("Tithy", mock.Anything, t1).Return(nil, errors.New("tithy gRPC error"))
	mockSwe.On("GetBalas", mock.Anything, t1).Return(nil, errors.New("balas gRPC error"))

	rr := httptest.NewRecorder()
	astrologyGenerateHandler(rr, req)

	// Even with an error, the handler should complete but report the error inside the results
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "tithy gRPC error")
	assert.Contains(t, rr.Body.String(), "balas gRPC error")
}
