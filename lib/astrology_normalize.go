package lib

import (
	"strings"

	"jenujari/go-sbc-webapp/sqls"
)

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
