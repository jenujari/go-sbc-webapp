package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"jenujari/go-sbc-webapp/sqls"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jenujari/go-swe-api/proto"
)

type AstrologyGenerateRequest struct {
	From  time.Time
	To    time.Time
	Clock string
}

type AstrologyGenerateResult struct {
	PanchangUpserted        int64
	PlanetPositionsUpserted int64
	Errors                  []string
}

type AstrologyStore interface {
	UpsertPanchang(ctx context.Context, arg sqls.UpsertPanchangParams) error
	UpsertPlanetPosition(ctx context.Context, arg sqls.UpsertPlanetPositionParams) error
}

type AstrologyGenerator struct {
	Swe   SweGrpcClient
	Store AstrologyStore
}

func NewAstrologyGenerator(swe SweGrpcClient, db *DBService) *AstrologyGenerator {
	if swe == nil || db == nil || db.Queries == nil {
		return nil
	}
	return &AstrologyGenerator{Swe: swe, Store: db.Queries}
}

func ParseAstrologyRange(fromDate, toDate, clock string) (AstrologyGenerateRequest, error) {
	if fromDate == "" || toDate == "" || clock == "" {
		return AstrologyGenerateRequest{}, fmt.Errorf("All fields (From Date, To Date, and Time) are required")
	}

	startDay, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return AstrologyGenerateRequest{}, fmt.Errorf("Invalid From Date format")
	}

	endDay, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return AstrologyGenerateRequest{}, fmt.Errorf("Invalid To Date format")
	}

	if endDay.Before(startDay) {
		return AstrologyGenerateRequest{}, fmt.Errorf("From Date must be before or equal to To Date")
	}

	if endDay.Sub(startDay) > 3653*24*time.Hour {
		return AstrologyGenerateRequest{}, fmt.Errorf("Date range cannot exceed 10 years")
	}

	if _, err := time.Parse("15:04", clock); err != nil {
		return AstrologyGenerateRequest{}, fmt.Errorf("Invalid Time format. Must be HH:MM")
	}

	return AstrologyGenerateRequest{From: startDay, To: endDay, Clock: clock}, nil
}

func (s *AstrologyGenerator) Generate(ctx context.Context, req AstrologyGenerateRequest) AstrologyGenerateResult {
	var result AstrologyGenerateResult
	if s == nil || s.Swe == nil || s.Store == nil {
		result.Errors = []string{"astrology generator is not configured"}
		return result
	}

	for day := req.From; !day.After(req.To); day = day.AddDate(0, 0, 1) {
		timestampStr := day.Format("2006-01-02") + "T" + req.Clock + ":00Z"
		parsedTime, _ := time.Parse(time.RFC3339, timestampStr)
		dbTime := pgtype.Timestamptz{Time: parsedTime, Valid: true}
		dayLabel := day.Format("2006-01-02")

		tithyResp, err := s.Swe.Tithy(ctx, timestampStr)
		if err != nil {
			result.Errors = appendLimited(result.Errors, dayLabel+": tithy gRPC fetch failed: "+err.Error())
		} else {
			pParams := sqls.UpsertPanchangParams{
				Time:      dbTime,
				Tithy:     int16(tithyResp.GetTithy()),
				Nakshatra: normalizeNakshatra(tithyResp.GetNakshatra()),
				Weekday:   normalizeWeekDay(tithyResp.GetWeekday()),
			}
			if err := s.Store.UpsertPanchang(ctx, pParams); err != nil {
				result.Errors = appendLimited(result.Errors, dayLabel+": panchang DB upsert failed: "+err.Error())
			} else {
				result.PanchangUpserted++
			}
		}

		balasResp, err := s.Swe.GetBalas(ctx, timestampStr)
		if err != nil {
			result.Errors = appendLimited(result.Errors, dayLabel+": planet positions gRPC fetch failed: "+err.Error())
			continue
		}

		for key, bala := range balasResp.GetResults() {
			if bala == nil || bala.GetCords() == nil {
				continue
			}
			if err := s.Store.UpsertPlanetPosition(ctx, planetPositionParams(dbTime, key, bala)); err != nil {
				result.Errors = appendLimited(result.Errors, dayLabel+" ("+key+"): planet DB upsert failed: "+err.Error())
				continue
			}
			result.PlanetPositionsUpserted++
		}
	}

	return result
}

func planetPositionParams(dbTime pgtype.Timestamptz, key string, bala *proto.PlanetBalas) sqls.UpsertPlanetPositionParams {
	cords := bala.GetCords()
	isRetro := cords.GetIsRetro()
	vargottama := cords.GetVargottama()

	var sign *sqls.SignType
	if signStr := cords.GetSign(); signStr != "" {
		sVal := normalizeSign(signStr)
		sign = &sVal
	}

	var nakPada *int16
	nakNameStr := ""
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

	var signLord *sqls.PlanetType
	if signLordStr := cords.GetSignLord(); signLordStr != "" {
		slVal := normalizePlanet(signLordStr)
		signLord = &slVal
	}

	var signLordship *sqls.RelType
	if signLordshipStr := cords.GetSignLordship(); signLordshipStr != "" {
		slsVal := normalizeRel(signLordshipStr)
		signLordship = &slsVal
	}

	var navamsaSign *sqls.SignType
	if navamsaSignStr := cords.GetNavamsaSign(); navamsaSignStr != "" {
		nsVal := normalizeSign(navamsaSignStr)
		navamsaSign = &nsVal
	}

	udayVal := bala.GetUdayBala()
	uchchaVal := bala.GetUchchaBala()
	vakraVal := bala.GetVakraBala()
	kshetraVal := bala.GetKshetraBala()
	navamshaVal := bala.GetNavamshaBala()

	return sqls.UpsertPlanetPositionParams{
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
