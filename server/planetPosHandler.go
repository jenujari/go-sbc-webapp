package server

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"net/http"
	"sort"
	"time"

	"github.com/jenujari/go-swe-api/proto"
	plLib "github.com/jenujari/planets-lib"
)

func planetPosHandler(w http.ResponseWriter, r *http.Request) {
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	webData := app.PageData()
	webData["currentTime"] = time.Now().Format("2006-01-02T15:04")

	html.RenderPage(w, webData, "planet_pos.html")
}

func positionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	webData := app.PageData()
	sweClient := app.SweClient

	datetime := r.FormValue("datetime")
	parsedDate, err := time.Parse("2006-01-02T15:04", datetime)

	webData["displayDate"] = datetime
	if err == nil {
		webData["displayDate"] = parsedDate.Format("January 02, 2006 15:04")
	}

	posResp, err := sweClient.GetPos(ctx, parsedDate.Format(time.RFC3339), "")
	if err != nil {
		config.GetLogger().Println("get pos failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	webData["planets"] = getSortedArray(posResp.GetResults())

	html.RenderPartial(w, "position_table.html", webData)
}

type PlanetTableRecord struct {
	Name         string
	Long         float64
	Lat          float64
	Longitude    string
	Latitude     string
	IsRetrograde string
	Sign         string
	Nakshatra    string
}

func getSortedArray(m map[string]*proto.PlanetCord) []PlanetTableRecord {
	result := []PlanetTableRecord{}

	for k, v := range m {

		longDMS := plLib.DMS{
			D:          int(v.LongitudeDms.D),
			M:          int(v.LongitudeDms.M),
			S:          float32(v.LongitudeDms.S),
			IsNegative: v.IsRetro,
		}

		latDMS := plLib.DMS{
			D:          int(v.LatitudeDms.D),
			M:          int(v.LatitudeDms.M),
			S:          float32(v.LatitudeDms.S),
			IsNegative: v.IsRetro,
		}

		isRetro := "No"
		if v.IsRetro {
			isRetro = "Yes"
		}

		result = append(result, PlanetTableRecord{
			Name:         k,
			Long:         v.Longitude,
			Lat:          v.Latitude,
			Longitude:    longDMS.String(),
			Latitude:     latDMS.String(),
			IsRetrograde: isRetro,
			Sign:         v.Sign,
			Nakshatra:    v.Nakshatra.Name,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Long > result[j].Long
	})

	return result
}
