package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"jenujari/go-sbc-webapp/config"
	"jenujari/go-sbc-webapp/html"
	"jenujari/go-sbc-webapp/lib"
	"jenujari/go-sbc-webapp/sqls"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jenujari/go-swe-api/proto"
)

type astrologyGenerateResultData struct {
	PanchangUpserted        int64
	PlanetPositionsUpserted int64
	Errors                  []string
}

func astrologyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	services, ok := ctx.Value("services").(map[string]any)
	if !ok || services == nil {
		http.Error(w, "Services not available", http.StatusInternalServerError)
		return
	}

	db, ok := services["db"].(*lib.DBService)
	if !ok || db == nil {
		http.Error(w, "Database connection is not available. Check db.url or DATABASE_URL.", http.StatusInternalServerError)
		return
	}

	sweClient, ok := services["sweClient"].(lib.SweGrpcClient)
	if !ok || sweClient == nil {
		http.Error(w, "Swe gRPC Client is not available.", http.StatusInternalServerError)
		return
	}

	fromDate := r.FormValue("from_date")
	toDate := r.FormValue("to_date")
	specificTime := r.FormValue("time")

	if fromDate == "" || toDate == "" || specificTime == "" {
		http.Error(w, "All fields (From Date, To Date, and Time) are required", http.StatusBadRequest)
		return
	}

	startDay, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		http.Error(w, "Invalid From Date format", http.StatusBadRequest)
		return
	}

	endDay, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		http.Error(w, "Invalid To Date format", http.StatusBadRequest)
		return
	}

	if endDay.Before(startDay) {
		http.Error(w, "From Date must be before or equal to To Date", http.StatusBadRequest)
		return
	}

	// 10 years max: 3653 days
	if endDay.Sub(startDay) > 3653*24*time.Hour {
		http.Error(w, "Date range cannot exceed 10 years", http.StatusBadRequest)
		return
	}

	// Validate time format (HH:MM)
	if _, err := time.Parse("15:04", specificTime); err != nil {
		http.Error(w, "Invalid Time format. Must be HH:MM", http.StatusBadRequest)
		return
	}

	var panchangUpserted int64
	var planetPositionsUpserted int64
	var errorsList []string

	for day := startDay; !day.After(endDay); day = day.AddDate(0, 0, 1) {
		// Target timestamp in UTC: "2026-05-25T09:30:00Z"
		timestampStr := day.Format("2006-01-02") + "T" + specificTime + ":00Z"
		parsedTime, _ := time.Parse(time.RFC3339, timestampStr)

		dbTime := pgtype.Timestamptz{Time: parsedTime, Valid: true}

		// 1. Fetch and Upsert Panchang
		tithyResp, err := sweClient.Tithy(ctx, timestampStr)
		if err != nil {
			errorsList = appendLimited(errorsList, day.Format("2006-01-02")+": tithy gRPC fetch failed: "+err.Error())
		} else {
			pParams := sqls.UpsertPanchangParams{
				Time:      dbTime,
				Tithy:     int16(tithyResp.GetTithy()),
				Nakshatra: normalizeNakshatra(tithyResp.GetNakshatra()),
				Weekday:   normalizeWeekDay(tithyResp.GetWeekday()),
			}
			if err := db.Queries.UpsertPanchang(ctx, pParams); err != nil {
				errorsList = appendLimited(errorsList, day.Format("2006-01-02")+": panchang DB upsert failed: "+err.Error())
			} else {
				panchangUpserted++
			}
		}

		// 2. Fetch and Upsert Planet Positions
		balasResp, err := sweClient.GetBalas(ctx, timestampStr)
		if err != nil {
			errorsList = appendLimited(errorsList, day.Format("2006-01-02")+": planet positions gRPC fetch failed: "+err.Error())
		} else {
			for key, bala := range balasResp.GetResults() {
				if bala == nil || bala.GetCords() == nil {
					continue
				}
				cords := bala.GetCords()

				isRetro := cords.GetIsRetro()
				vargottama := cords.GetVargottama()

				signStr := cords.GetSign()
				var sign *sqls.SignType
				if signStr != "" {
					sVal := normalizeSign(signStr)
					sign = &sVal
				}

				nakNameStr := ""
				var nakPada *int16
				if cords.GetNakshatra() != nil {
					nakNameStr = cords.GetNakshatra().GetName()
					padaVal := int16(cords.GetNakshatra().GetPada())
					nakPada = &padaVal
				}

				var nakName *sqls.NakshatraType
				if nakNameStr != "" {
					nVal := normalizeNakshatra(nakNameStr)
					nakName = &nVal
				}

				signLordStr := cords.GetSignLord()
				var signLord *sqls.PlanetType
				if signLordStr != "" {
					slVal := normalizePlanet(signLordStr)
					signLord = &slVal
				}

				signLordshipStr := cords.GetSignLordship()
				var signLordship *sqls.RelType
				if signLordshipStr != "" {
					slsVal := normalizeRel(signLordshipStr)
					signLordship = &slsVal
				}

				navamsaSignStr := cords.GetNavamsaSign()
				var navamsaSign *sqls.SignType
				if navamsaSignStr != "" {
					nsVal := normalizeSign(navamsaSignStr)
					navamsaSign = &nsVal
				}

				udayVal := bala.GetUdayBala()
				uchchaVal := bala.GetUchchaBala()
				vakraVal := bala.GetVakraBala()
				kshetraVal := bala.GetKshetraBala()
				navamshaVal := bala.GetNavamshaBala()

				ppParams := sqls.UpsertPlanetPositionParams{
					ObservationTime: dbTime,
					PlanetName:      normalizePlanet(key),
					Longitude:       cords.GetLongitude(),
					Latitude:        cords.GetLatitude(),
					Distance:        cords.GetDistance(),
					SpeedLong:       cords.GetSpeedLong(),
					SpeedLat:        cords.GetSpeedLat(),
					SpeedDist:       cords.GetSpeedDist(),
					SpeedCategory:   normalizeSpeed(cords.GetSpeedCategory()),
					Vedha:           normalizeVedha(cords.GetVedha()),
					Sign:            sign,
					NakshatraName:   nakName,
					NakshatraPada:   nakPada,
					IsRetro:         &isRetro,
					SignLord:        signLord,
					SignLordship:    signLordship,
					NavamsaSign:     navamsaSign,
					Vargottama:      &vargottama,
					UdayBala:        &udayVal,
					UchchaBala:      &uchchaVal,
					VakraBala:       &vakraVal,
					KshetraBala:     &kshetraVal,
					NavamshaBala:    &navamshaVal,
					LongitudeDms:    marshalDMS(cords.GetLongitudeDms()),
					LatitudeDms:     marshalDMS(cords.GetLatitudeDms()),
					SpeedLongDms:    marshalDMS(cords.GetSpeedLongDms()),
				}

				if err := db.Queries.UpsertPlanetPosition(ctx, ppParams); err != nil {
					errorsList = appendLimited(errorsList, day.Format("2006-01-02")+" ("+key+"): planet DB upsert failed: "+err.Error())
				} else {
					planetPositionsUpserted++
				}
			}
		}
	}

	data := astrologyGenerateResultData{
		PanchangUpserted:        panchangUpserted,
		PlanetPositionsUpserted: planetPositionsUpserted,
		Errors:                  errorsList,
	}

	renderAstrologyGenerateResult(w, data)
}

func renderAstrologyGenerateResult(w http.ResponseWriter, data astrologyGenerateResultData) {
	w.WriteHeader(http.StatusOK)
	tpl, err := html.GetTpl().Clone()
	if err != nil {
		config.GetLogger().Println("template clone failed:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	tpl, err = tpl.ParseFS(html.GetViewsFs(), "astrology_generate_result.html")
	if err != nil {
		config.GetLogger().Println("template not found:", err)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	if err := tpl.ExecuteTemplate(w, "astrology_generate_result.html", data); err != nil {
		config.GetLogger().Println("template execution failed:", err)
		http.Error(w, "template execution failed", http.StatusInternalServerError)
	}
}

func marshalDMS(dms *proto.DMS) []byte {
	if dms == nil {
		return nil
	}
	data := map[string]any{
		"d":           dms.GetD(),
		"m":           dms.GetM(),
		"s":           dms.GetS(),
		"is_negative": dms.GetIsNegative(),
	}
	b, _ := json.Marshal(data)
	return b
}

func normalizeNakshatra(s string) sqls.NakshatraType {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	switch s {
	case "ashwini":
		return sqls.NakshatraTypeAshwini
	case "bharani":
		return sqls.NakshatraTypeBharani
	case "krittika":
		return sqls.NakshatraTypeKrittika
	case "rohini":
		return sqls.NakshatraTypeRohini
	case "mrigashirsha", "mrigashira", "mrigasira":
		return sqls.NakshatraTypeMrigashirsha
	case "ardra":
		return sqls.NakshatraTypeArdra
	case "punarvasu":
		return sqls.NakshatraTypePunarvasu
	case "pushya":
		return sqls.NakshatraTypePushya
	case "ashlesha":
		return sqls.NakshatraTypeAshlesha
	case "magha":
		return sqls.NakshatraTypeMagha
	case "purva phalguni", "purvaphalguni", "pubba":
		return sqls.NakshatraTypePurvaPhalguni
	case "uttara phalguni", "uttaraphalguni", "uttara":
		return sqls.NakshatraTypeUttaraPhalguni
	case "hasta":
		return sqls.NakshatraTypeHasta
	case "chitra", "chithra":
		return sqls.NakshatraTypeChitra
	case "swati":
		return sqls.NakshatraTypeSwati
	case "vishakha":
		return sqls.NakshatraTypeVishakha
	case "anuradha":
		return sqls.NakshatraTypeAnuradha
	case "jyestha", "jyeshtha":
		return sqls.NakshatraTypeJyestha
	case "moola", "mula":
		return sqls.NakshatraTypeMoola
	case "purva ashadha", "purvaashadha", "poorvashadha":
		return sqls.NakshatraTypePurvaAshadha
	case "uttara ashadha", "uttaraashadha", "uttarashadha":
		return sqls.NakshatraTypeUttaraAshadha
	case "abhijit":
		return sqls.NakshatraTypeAbhijit
	case "shravana", "shravan":
		return sqls.NakshatraTypeShravana
	case "dhanishtha", "dhanishta":
		return sqls.NakshatraTypeDhanishtha
	case "shatabhisha", "shatataraka", "satabhisha":
		return sqls.NakshatraTypeShatabhisha
	case "purva bhadrapada", "purvabhadrapada", "poorvabhadrapada":
		return sqls.NakshatraTypePurvaBhadrapada
	case "uttara bhadrapada", "uttarabhadrapada":
		return sqls.NakshatraTypeUttaraBhadrapada
	case "revati":
		return sqls.NakshatraTypeRevati
	}
	return sqls.NakshatraType(s)
}

func normalizePlanet(s string) sqls.PlanetType {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "sun", "surya":
		return sqls.PlanetTypeSun
	case "moon", "chandra":
		return sqls.PlanetTypeMoon
	case "mercury", "budha":
		return sqls.PlanetTypeMercury
	case "venus", "shukra":
		return sqls.PlanetTypeVenus
	case "mars", "mangal":
		return sqls.PlanetTypeMars
	case "jupiter", "guru":
		return sqls.PlanetTypeJupiter
	case "saturn", "shani":
		return sqls.PlanetTypeSaturn
	case "uranus":
		return sqls.PlanetTypeUranus
	case "neptune":
		return sqls.PlanetTypeNeptune
	case "pluto":
		return sqls.PlanetTypePluto
	case "rahu":
		return sqls.PlanetTypeRahu
	case "ketu":
		return sqls.PlanetTypeKetu
	}
	return sqls.PlanetType(s)
}

func normalizeSign(s string) sqls.SignType {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "aries":
		return sqls.SignTypeAries
	case "taurus":
		return sqls.SignTypeTaurus
	case "gemini":
		return sqls.SignTypeGemini
	case "cancer":
		return sqls.SignTypeCancer
	case "leo":
		return sqls.SignTypeLeo
	case "virgo":
		return sqls.SignTypeVirgo
	case "libra":
		return sqls.SignTypeLibra
	case "scorpio":
		return sqls.SignTypeScorpio
	case "sagittarius":
		return sqls.SignTypeSagittarius
	case "capricorn":
		return sqls.SignTypeCapricorn
	case "aquarius":
		return sqls.SignTypeAquarius
	case "pisces":
		return sqls.SignTypePisces
	}
	return sqls.SignType(s)
}

func normalizeWeekDay(s string) sqls.WeekDayType {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "sunday":
		return sqls.WeekDayTypeSunday
	case "monday":
		return sqls.WeekDayTypeMonday
	case "tuesday":
		return sqls.WeekDayTypeTuesday
	case "wednesday":
		return sqls.WeekDayTypeWednesday
	case "thursday":
		return sqls.WeekDayTypeThursday
	case "friday":
		return sqls.WeekDayTypeFriday
	case "saturday":
		return sqls.WeekDayTypeSaturday
	}
	return sqls.WeekDayType(s)
}

func normalizeSpeed(s string) sqls.SpeedType {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, "_", "-")
	switch s {
	case "kutil":
		return sqls.SpeedTypeKutil
	case "ati-vakra", "ativakra":
		return sqls.SpeedTypeAtiVakra
	case "vakra":
		return sqls.SpeedTypeVakra
	case "ati-mand", "atimand":
		return sqls.SpeedTypeAtiMand
	case "mand":
		return sqls.SpeedTypeMand
	case "madhyam":
		return sqls.SpeedTypeMadhyam
	case "sama":
		return sqls.SpeedTypeSama
	case "sheeghra":
		return sqls.SpeedTypeSheeghra
	case "ati-sheeghra", "atisheeghra":
		return sqls.SpeedTypeAtiSheeghra
	case "n/a", "na":
		return sqls.SpeedTypeNA
	}
	return sqls.SpeedTypeNA
}

func normalizeVedha(s string) sqls.VedhaType {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "left":
		return sqls.VedhaTypeLeft
	case "right":
		return sqls.VedhaTypeRight
	case "front":
		return sqls.VedhaTypeFront
	case "no":
		return sqls.VedhaTypeNo
	case "n/a", "na":
		return sqls.VedhaTypeNA
	}
	return sqls.VedhaTypeNA
}

func normalizeRel(s string) sqls.RelType {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "friend":
		return sqls.RelTypeFriend
	case "neutral":
		return sqls.RelTypeNeutral
	case "enemy":
		return sqls.RelTypeEnemy
	case "self":
		return sqls.RelTypeSelf
	}
	return sqls.RelTypeNeutral
}
