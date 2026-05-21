package lib

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jenujari/go-swe-api/proto"
)

type PlanetShadbalaService interface {
	GetPlanetShadbala(ctx context.Context, timestamp string) (PlanetShadbalaView, error)
}

type SwePlanetShadbalaService struct {
	client SweGrpcClient
}

type PlanetShadbalaView struct {
	Timestamp string
	Planets   []PlanetShadbalaRecord
}

type PlanetShadbalaRecord struct {
	ID             string
	Name           string
	Symbol         string
	UdayBala       float64
	UchchaBala     float64
	VakraBala      float64
	NavamshaBala   float64
	KshetraBala    float64
	Total          float64
	Sign           string
	Nakshatra      string
	Speed          string
	State          string
	Retrograde     bool
	SpeedCategory  string
	Vedha          string
	LongitudeDMS   string
	LatitudeDMS    string
	SignLord       string
	SignLordship   string
	NavamsaSign    string
	Vargottama     string
	PowerSortOrder int
	NameSortOrder  int
}

func NewPlanetShadbalaService(client SweGrpcClient) PlanetShadbalaService {
	return &SwePlanetShadbalaService{client: client}
}

func (s *SwePlanetShadbalaService) GetPlanetShadbala(ctx context.Context, timestamp string) (PlanetShadbalaView, error) {
	if s == nil || s.client == nil {
		return PlanetShadbalaView{}, fmt.Errorf("planet shadbala service is not configured")
	}

	resp, err := s.client.GetBalas(ctx, timestamp)
	if err != nil {
		return PlanetShadbalaView{}, fmt.Errorf("fetch planet shadbala: %w", err)
	}

	return MapBalasResponse(timestamp, resp), nil
}

func MapBalasResponse(timestamp string, resp *proto.BalasResponse) PlanetShadbalaView {
	view := PlanetShadbalaView{Timestamp: timestamp, Planets: []PlanetShadbalaRecord{}}
	if resp == nil || len(resp.GetResults()) == 0 {
		return view
	}

	for key, bala := range resp.GetResults() {
		if bala == nil {
			continue
		}

		cords := bala.GetCords()
		name := titlePlanetName(key)
		if cords != nil && cords.GetName() != "" {
			name = titlePlanetName(cords.GetName())
		}

		retrograde := false
		sign := "-"
		nakshatra := "-"
		speed := "-"
		speedCategory := "-"
		vedha := "-"
		longitudeDMS := "-"
		latitudeDMS := "-"
		signLord := "-"
		signLordship := "-"
		navamsaSign := "-"
		vargottama := "No"
		if cords != nil {
			retrograde = cords.GetIsRetro()
			if cords.GetSign() != "" {
				sign = cords.GetSign()
			}
			if cords.GetNakshatra() != nil && cords.GetNakshatra().GetName() != "" {
				nakshatra = cords.GetNakshatra().GetName()
			}
			if cords.GetSpeedLong() != 0 {
				speed = fmt.Sprintf("%.4f° / day", cords.GetSpeedLong())
			}
			if cords.GetSpeedCategory() != "" {
				speedCategory = cords.GetSpeedCategory()
			}
			if cords.GetVedha() != "" {
				vedha = cords.GetVedha()
			}
			longitudeDMS = formatProtoDMS(cords.GetLongitudeDms())
			latitudeDMS = formatProtoDMS(cords.GetLatitudeDms())
			if cords.GetSignLord() != "" {
				signLord = cords.GetSignLord()
			}
			if cords.GetSignLordship() != "" {
				signLordship = cords.GetSignLordship()
			}
			if cords.GetNavamsaSign() != "" {
				navamsaSign = cords.GetNavamsaSign()
			}
			if cords.GetVargottama() {
				vargottama = "Yes"
			}
		}

		record := PlanetShadbalaRecord{
			ID:            planetID(name),
			Name:          name,
			Symbol:        planetSymbol(name),
			UdayBala:      clampPercent(bala.GetUdayBala()),
			UchchaBala:    clampPercent(bala.GetUchchaBala()),
			VakraBala:     clampPercent(bala.GetVakraBala()),
			NavamshaBala:  clampPercent(bala.GetNavamshaBala()),
			KshetraBala:   clampPercent(bala.GetKshetraBala()),
			Sign:          sign,
			Nakshatra:     nakshatra,
			Speed:         speed,
			Retrograde:    retrograde,
			SpeedCategory: speedCategory,
			Vedha:         vedha,
			LongitudeDMS:  longitudeDMS,
			LatitudeDMS:   latitudeDMS,
			SignLord:      signLord,
			SignLordship:  signLordship,
			NavamsaSign:   navamsaSign,
			Vargottama:    vargottama,
			State:         "Direct",
		}
		if record.Retrograde {
			record.State = "Retrograde"
		}
		record.Total = clampPercent((record.UdayBala + record.UchchaBala + record.VakraBala + record.NavamshaBala + record.KshetraBala) / 5)
		record.PowerSortOrder = 100 - int(record.Total)
		record.NameSortOrder = int([]rune(strings.ToLower(record.Name))[0])

		view.Planets = append(view.Planets, record)
	}

	sort.SliceStable(view.Planets, func(i, j int) bool {
		return planetSortOrder(view.Planets[i].Name) < planetSortOrder(view.Planets[j].Name)
	})

	return view
}

func formatProtoDMS(dms *proto.DMS) string {
	if dms == nil {
		return "-"
	}
	sign := ""
	if dms.GetIsNegative() {
		sign = "-"
	}
	return fmt.Sprintf("%s%d° %d′ %.2f″", sign, dms.GetD(), dms.GetM(), dms.GetS())
}

func clampPercent(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func titlePlanetName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Unknown"
	}
	words := strings.Fields(strings.ReplaceAll(name, "_", " "))
	for i, word := range words {
		word = strings.ToLower(word)
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

func planetID(name string) string {
	return strings.NewReplacer(" ", "-", "_", "-").Replace(strings.ToLower(strings.TrimSpace(name)))
}

func planetSymbol(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "sun", "surya":
		return "☉"
	case "moon", "chandra":
		return "☾"
	case "mars", "mangal":
		return "♂"
	case "mercury", "budha":
		return "☿"
	case "jupiter", "guru":
		return "♃"
	case "venus", "shukra":
		return "♀"
	case "saturn", "shani":
		return "♄"
	case "rahu":
		return "☊"
	case "ketu":
		return "☋"
	default:
		return "✦"
	}
}

func planetSortOrder(name string) int {
	order := map[string]int{"sun": 1, "moon": 2, "mars": 3, "mercury": 4, "jupiter": 5, "venus": 6, "saturn": 7, "rahu": 8, "ketu": 9}
	if v, ok := order[strings.ToLower(strings.TrimSpace(name))]; ok {
		return v
	}
	return 100
}
