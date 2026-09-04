package server

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"net/http"
	"time"
)

func tithyHandler(w http.ResponseWriter, r *http.Request) {
	app, ok := requestApp(r)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	webData := app.WebData
	webData["currentDate"] = time.Now().Format("2006-01-02")

	html.RenderPage(w, webData, "tithy.html")
}

func tithyTableHandler(w http.ResponseWriter, r *http.Request) {
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
	webData := app.WebData
	sweClient := app.SweClient

	selectedDate := r.FormValue("selected_date")
	if selectedDate == "" {
		http.Error(w, "selected_date is required", http.StatusBadRequest)
		return
	}

	baseDate, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		http.Error(w, "invalid selected_date", http.StatusBadRequest)
		return
	}

	records := make([]TithyTableRecord, 0, 21)
	for dayOffset := -10; dayOffset <= 10; dayOffset++ {
		day := baseDate.AddDate(0, 0, dayOffset)
		timestamp := time.Date(day.Year(), day.Month(), day.Day(), 3, 30, 0, 0, time.UTC)

		resp, err := sweClient.Tithy(ctx, timestamp.Format(time.RFC3339))
		if err != nil {
			config.GetLogger().Println("tithy fetch failed", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		records = append(records, TithyTableRecord{
			DateLabel:         day.Format("January 02, 2006"),
			DateValue:         day.Format("2006-01-02"),
			Timestamp:         timestamp.Format(time.RFC3339),
			Weekday:           resp.GetWeekday(),
			Nakshatra:         resp.GetNakshatra(),
			TithyValue:        resp.GetTithy(),
			DisplayPaksha:     getPaksha(resp.GetTithy()),
			DisplayTithyValue: getDisplayTithyValue(resp.GetTithy()),
		})
	}

	webData["displayDate"] = baseDate.Format("January 02, 2006")
	webData["selectedDateValue"] = baseDate.Format("2006-01-02")
	webData["rangeStart"] = records[0].DateLabel
	webData["rangeEnd"] = records[len(records)-1].DateLabel
	webData["tithies"] = records

	html.RenderPartial(w, "tithy_table.html", webData)
}

type TithyTableRecord struct {
	DateLabel         string
	DateValue         string
	Timestamp         string
	Weekday           string
	Nakshatra         string
	TithyValue        int32
	DisplayPaksha     string
	DisplayTithyValue int32
}

func getPaksha(tithy int32) string {
	if tithy >= 1 && tithy <= 15 {
		return "Shukla Paksha"
	}

	return "Krishna Paksha"
}

func getDisplayTithyValue(tithy int32) int32 {
	if tithy <= 15 {
		return tithy
	}

	return tithy - 15
}
