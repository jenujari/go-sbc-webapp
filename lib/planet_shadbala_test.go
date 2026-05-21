package lib

import (
	"context"
	"testing"

	"github.com/jenujari/go-swe-api/proto"
	"github.com/stretchr/testify/mock"
)

func TestSwePlanetShadbalaServiceGetPlanetShadbalaHappyPath(t *testing.T) {
	client := NewMockSweGrpcClient(t)
	client.EXPECT().GetBalas(mock.Anything, "2026-05-21T00:00:00Z").Return(&proto.BalasResponse{Results: map[string]*proto.PlanetBalas{
		"sun": {
			Cords:        &proto.PlanetCord{Name: "sun", Sign: "Leo", Nakshatra: &proto.NakshatraPada{Name: "Magha"}, SpeedLong: 0.98},
			UdayBala:     90,
			UchchaBala:   80,
			VakraBala:    0,
			NavamshaBala: 70,
			KshetraBala:  60,
		},
	}}, nil)

	view, err := NewPlanetShadbalaService(client).GetPlanetShadbala(context.Background(), "2026-05-21T00:00:00Z")
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
	if len(view.Planets) != 1 {
		t.Fatalf("expected one planet, got %d", len(view.Planets))
	}
	planet := view.Planets[0]
	if planet.Name != "Sun" || planet.Symbol != "☉" || planet.Sign != "Leo" || planet.Nakshatra != "Magha" {
		t.Fatalf("unexpected mapped planet: %+v", planet)
	}
	if planet.Total != 60 {
		t.Fatalf("expected total 60, got %v", planet.Total)
	}
}

func TestSwePlanetShadbalaServiceGetPlanetShadbalaEmptyData(t *testing.T) {
	client := NewMockSweGrpcClient(t)
	client.EXPECT().GetBalas(mock.Anything, "2026-05-21T00:00:00Z").Return(&proto.BalasResponse{Results: map[string]*proto.PlanetBalas{}}, nil)

	view, err := NewPlanetShadbalaService(client).GetPlanetShadbala(context.Background(), "2026-05-21T00:00:00Z")
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
	if len(view.Planets) != 0 {
		t.Fatalf("expected empty planets, got %d", len(view.Planets))
	}
}

func TestMapBalasResponseHandlesNilFields(t *testing.T) {
	view := MapBalasResponse("2026-05-21T00:00:00Z", &proto.BalasResponse{Results: map[string]*proto.PlanetBalas{
		"mars": nil,
		"moon": {UdayBala: 10, UchchaBala: 20, VakraBala: 30, NavamshaBala: 40, KshetraBala: 50},
	}})

	if len(view.Planets) != 1 {
		t.Fatalf("expected one non-nil planet, got %d", len(view.Planets))
	}
	if view.Planets[0].Sign != "-" || view.Planets[0].Nakshatra != "-" || view.Planets[0].Speed != "-" {
		t.Fatalf("expected nil-safe defaults, got %+v", view.Planets[0])
	}
}
