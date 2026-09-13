package lib

import "testing"

func TestPageDataDoesNotMutateSharedWebData(t *testing.T) {
	app := &App{WebData: WebData{"appname": "webapp"}}

	page := app.PageData()
	page["currentTime"] = "mutated"

	if _, ok := app.WebData["currentTime"]; ok {
		t.Fatal("shared WebData was mutated by request-scoped PageData")
	}
	if app.WebData["appname"] != "webapp" {
		t.Fatalf("expected appname preserved, got %v", app.WebData["appname"])
	}
}
