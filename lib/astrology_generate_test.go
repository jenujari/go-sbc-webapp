package lib

import (
	"context"
	"testing"

	"jenujari/go-sbc-webapp/sqls"

	"github.com/jenujari/go-swe-api/proto"
	"github.com/stretchr/testify/mock"
)

type stubAstrologyStore struct {
	panchang int
	planets  int
}

func (s *stubAstrologyStore) UpsertPanchang(ctx context.Context, arg sqls.UpsertPanchangParams) error {
	s.panchang++
	return nil
}

func (s *stubAstrologyStore) UpsertPlanetPosition(ctx context.Context, arg sqls.UpsertPlanetPositionParams) error {
	s.planets++
	return nil
}

func TestParseAstrologyRangeValidation(t *testing.T) {
	cases := []struct {
		from, to, clock, want string
	}{
		{"", "2026-05-28", "09:30", "All fields (From Date, To Date, and Time) are required"},
		{"bad", "2026-05-28", "09:30", "Invalid From Date format"},
		{"2026-05-28", "2026-05-25", "09:30", "From Date must be before or equal to To Date"},
		{"2016-05-25", "2026-05-27", "09:30", "Date range cannot exceed 10 years"},
		{"2026-05-25", "2026-05-25", "99:99", "Invalid Time format. Must be HH:MM"},
	}
	for _, tc := range cases {
		_, err := ParseAstrologyRange(tc.from, tc.to, tc.clock)
		if err == nil || err.Error() != tc.want {
			t.Fatalf("from=%s to=%s clock=%s: got %v want %s", tc.from, tc.to, tc.clock, err, tc.want)
		}
	}
}

func TestNormalizeAliases(t *testing.T) {
	if normalizePlanet("surya") != sqls.PlanetTypeSun {
		t.Fatal("surya should map to sun")
	}
	if normalizeNakshatra("pubba") != sqls.NakshatraTypePurvaPhalguni {
		t.Fatal("pubba should map to purva phalguni")
	}
	if normalizeSpeed("ativakra") != sqls.SpeedTypeAtiVakra {
		t.Fatal("ativakra should map to ati-vakra")
	}
}

func TestAstrologyGeneratorGenerate(t *testing.T) {
	client := NewMockSweGrpcClient(t)
	store := &stubAstrologyStore{}
	ts := "2026-05-25T09:30:00Z"
	client.EXPECT().Tithy(mock.Anything, ts).Return(&proto.TithyResponse{Tithy: 5, Nakshatra: "Magha", Weekday: "Monday"}, nil)
	client.EXPECT().GetBalas(mock.Anything, ts).Return(&proto.BalasResponse{Results: map[string]*proto.PlanetBalas{
		"sun": {Cords: &proto.PlanetCord{Name: "sun", Sign: "Leo", Longitude: 120.5}},
	}}, nil)

	req, err := ParseAstrologyRange("2026-05-25", "2026-05-25", "09:30")
	if err != nil {
		t.Fatal(err)
	}
	result := (&AstrologyGenerator{Swe: client, Store: store}).Generate(context.Background(), req)
	if result.PanchangUpserted != 1 || result.PlanetPositionsUpserted != 1 || len(result.Errors) != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if store.panchang != 1 || store.planets != 1 {
		t.Fatalf("unexpected store calls: %+v", store)
	}
}
