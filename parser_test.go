package main

import (
	"testing"
)

type MetarTestScenario struct {
	RawValue          string
	ExpectedStation   string
	ExpectedDay       int32
	ExpectedVisiblity string
	ExpectedWindSpeed float32
}

func init() {
}

func TestParseFullMetar(t *testing.T) {
	testMetarBasic := MetarTestScenario{
		"KORD 210051Z 15007KT 10SM OVC060 05/01 A3010 RMK AO2 RAE02 SLP200 P0000 T00500011",
		"KORD",
		21,
		"10 miles",
		7,
	}

	checkMetarScenario(t, testMetarBasic)

	testMetarWithAuto := MetarTestScenario{
		"PANV 260236Z AUTO 04014G19KT 10SM OVC085 M11/M14 A2989 RMK AO1",
		"PANV",
		26,
		"10 miles",
		14,
	}

	checkMetarScenario(t, testMetarWithAuto)

	testMetarWithMultipleClouds := MetarTestScenario{
		"KPWK 300251Z 16009KT 10SM FEW150 BKN200 OVC250 00/M07 A3043 RMK AO2 SLP312 T00001067 58019",
		"KPWK",
		30,
		"10 miles",
		9,
	}

	checkMetarScenario(t, testMetarWithMultipleClouds)

	testMetarWithClear := MetarTestScenario{
		"KPWK 300252Z 15007KT 10SM CLR 00/M07 A3045 RMK AO2 SLP318 T00001067 58020",
		"KPWK",
		30,
		"10 miles",
		7,
	}

	checkMetarScenario(t, testMetarWithClear)
}

func checkMetarScenario(t *testing.T, testMetar MetarTestScenario) {
	metar, _ := ParseMetar(testMetar.RawValue)

	t.Logf("Evaluating %+v ", metar)
	if metar.Station != testMetar.ExpectedStation {
		t.Error("Station not correct")
	}
	if metar.Day != testMetar.ExpectedDay {
		t.Error("Day not correct")
	}
	if metar.Visibility != testMetar.ExpectedVisiblity {
		t.Error("Visiblity not correct")
	}
	if metar.WindSpeed != testMetar.ExpectedWindSpeed {
		t.Error("Wind speed not correct")
	}
	details := GetDetailMetar(metar)
	t.Logf("Details: %v", details)
	if details == "" {
		t.Error("Details failed to parse")
	}
}

func TestParseDayTime(t *testing.T) {
	const testDateTime = "210051Z"
	day, time := parseDayTime(testDateTime)
	t.Logf("Received %s, %v", day, time)
	if day != 21 {
		t.Error("day not correct")
	}
	if time.Hour() != 00 {
		t.Error("Time hour not correct")
	}
	if time.Minute() != 51 {
		t.Error("Time minute not correct")
	}
	t.Log("OK")
}

func TestParseWind(t *testing.T) {
	const testWind = "18055KT"
	direction, wind, degrees, gust := parseWind(testWind)
	t.Logf("Received %v, %v, %v, %v", direction, wind, degrees, gust)
	if direction != "S" {
		t.Error("Direction not correct")
	}
	if degrees != 180 {
		t.Error("Degrees not correct")
	}
	if wind != 55 {
		t.Error("Wind not correct")
	}
	if gust != 55 {
		t.Error("Gust not correct")
	}
	t.Log("OK")
}

func TestParseWindWithGust(t *testing.T) {
	const testWindWithGust = "34014G21KT"
	direction, wind, degrees, gust := parseWind(testWindWithGust)
	t.Logf("Received %v, %v, %v, %v", direction, wind, degrees, gust)
	if direction != "NNW" {
		t.Error("Direction not correct")
	}
	if degrees != 340 {
		t.Error("Degrees not correct")
	}
	if wind != 14 {
		t.Error("Wind not correct")
	}
	if gust != 21 {
		t.Error("Gust not correct")
	}
	t.Log("OK")
}

func TestParseVisibilityFraction(t *testing.T) {
	const testVisibility = "1/2SM"
	distance := parseVisibility(testVisibility)
	t.Logf("Received %v ", distance)
	if distance != "1/2 miles" {
		t.Error("Visiblity not correct")
	}
	t.Log("OK")
}

func TestParseCloudsMultiple(t *testing.T) {
	const testClouds = "FEW200 SCT250"
	clouds := parseClouds(testClouds)
	t.Logf("Received %v ", clouds)
	if len(clouds) != 2 {
		t.Error("Received wrong count of clouds")
	}
	t.Log("OK")
}

func TestParseCloudItem(t *testing.T) {
	const testCloud = "FEW200"
	cloud := parseCloudDescription(testCloud)
	t.Logf("Received %v ", cloud)
	if cloud != "FEW at 20000" {
		t.Error("Received wrong cloud value")
	}
	t.Log("OK")
}

func TestParseTempDew(t *testing.T) {
	const testTemp = "05/M01"
	temperature, dewPoint := parseTempDew(testTemp)
	t.Logf("Received %v, %v", temperature, dewPoint)
	if temperature != 5.0 {
		t.Error("Received wrong temperature")
	}

	if dewPoint != -1.0 {
		t.Error("Received wrong dew point")
	}
	t.Log("OK")
}

func TestParsePressure(t *testing.T) {
	const testPressure = "A3006"
	pressure := parsePressure(testPressure)
	t.Logf("Received %v", pressure)
	if pressure != 30.06 {
		t.Error("Received wrong pressure")
	}
}
