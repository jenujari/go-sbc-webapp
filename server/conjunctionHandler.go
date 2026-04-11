package server

import (
	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
	"net/http"
	"time"

	plLib "github.com/jenujari/planets-lib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func conjunctionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	services := ctx.Value("services").(map[string]any)

	webData := services["webData"].(lib.WebData)
	webData["planets"] = plLib.PLANET_NAMES
	webData["defaultPlanet1"] = plLib.SUN
	webData["defaultPlanet2"] = plLib.MOON

	now := time.Now().UTC().Truncate(time.Minute)
	webData["defaultStart"] = now.Format("2006-01-02T15:04")
	webData["defaultEnd"] = now.Add(24 * time.Hour).Format("2006-01-02T15:04")

	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.ParseFS(html.GetViewsFs(), "layout.html", "conjunction.html")
	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "layout.html", webData)
	if err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func conjunctionSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	services := ctx.Value("services").(map[string]any)
	webData := services["webData"].(lib.WebData)
	sweClient := services["sweClient"].(lib.SweGrpcClient)
	delete(webData, "conjunctionError")
	delete(webData, "conjunctionNotFound")
	delete(webData, "conjunctionResult")

	searchRecord := ConjunctionSearchRecord{
		Planet1:      r.FormValue("planet1"),
		Planet2:      r.FormValue("planet2"),
		Start:        r.FormValue("start"),
		End:          r.FormValue("end"),
		DisplayStart: r.FormValue("start"),
		DisplayEnd:   r.FormValue("end"),
	}

	webData["conjunctionSearch"] = searchRecord

	if searchRecord.Planet1 == "" || searchRecord.Planet2 == "" || searchRecord.Start == "" || searchRecord.End == "" {
		webData["conjunctionError"] = "planet1, planet2, start, and end are required."
		renderConjunctionResult(w, webData)
		return
	}

	if searchRecord.Planet1 == searchRecord.Planet2 {
		webData["conjunctionError"] = "Please choose two different planets."
		renderConjunctionResult(w, webData)
		return
	}

	startTime, err := time.Parse("2006-01-02T15:04", searchRecord.Start)
	if err != nil {
		webData["conjunctionError"] = "Invalid start date."
		renderConjunctionResult(w, webData)
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04", searchRecord.End)
	if err != nil {
		webData["conjunctionError"] = "Invalid end date."
		renderConjunctionResult(w, webData)
		return
	}

	if !endTime.After(startTime) {
		webData["conjunctionError"] = "End must be after start."
		renderConjunctionResult(w, webData)
		return
	}

	searchRecord.DisplayStart = startTime.Format("January 02, 2006 15:04 UTC")
	searchRecord.DisplayEnd = endTime.Format("January 02, 2006 15:04 UTC")
	webData["conjunctionSearch"] = searchRecord

	resp, err := sweClient.FindConjunction(
		ctx,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		searchRecord.Planet1,
		searchRecord.Planet2,
		1,
		1.0/24.0,
	)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			webData["conjunctionNotFound"] = true
			renderConjunctionResult(w, webData)
			return
		}

		config.GetLogger().Println("find conjunction failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	result, err := getConjunctionResultRecord(resp)
	if err != nil {
		config.GetLogger().Println("invalid conjunction response", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	webData["conjunctionResult"] = result
	renderConjunctionResult(w, webData)
}

type ConjunctionSearchRecord struct {
	Planet1      string
	Planet2      string
	Start        string
	End          string
	DisplayStart string
	DisplayEnd   string
}

type ConjunctionResultRecord struct {
	Start        string
	End          string
	DisplayStart string
	DisplayEnd   string
}

func renderConjunctionResult(w http.ResponseWriter, webData lib.WebData) {
	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tpl, err = tpl.ParseFS(html.GetViewsFs(), "conjunction_result.html")
	if err != nil {
		config.GetLogger().Println("template not found", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "conjunction_result.html", webData)
	if err != nil {
		config.GetLogger().Println("template execution failed", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func getConjunctionResultRecord(resp interface {
	GetStart() string
	GetEnd() string
}) (ConjunctionResultRecord, error) {
	startTime, err := time.Parse(time.RFC3339, resp.GetStart())
	if err != nil {
		return ConjunctionResultRecord{}, err
	}

	endTime, err := time.Parse(time.RFC3339, resp.GetEnd())
	if err != nil {
		return ConjunctionResultRecord{}, err
	}

	return ConjunctionResultRecord{
		Start:        resp.GetStart(),
		End:          resp.GetEnd(),
		DisplayStart: startTime.Format("January 02, 2006 15:04:05 UTC"),
		DisplayEnd:   endTime.Format("January 02, 2006 15:04:05 UTC"),
	}, nil
}
