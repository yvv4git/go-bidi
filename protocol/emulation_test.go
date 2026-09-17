package protocol

import (
	"encoding/json"
	"testing"
)

func TestEmulationSetGeolocationOverrideParams(t *testing.T) {
	params := EmulationSetGeolocationOverrideParams{
		Context: "c1",
		Coordinates: &EmulationCoordinates{
			Latitude:  1.5,
			Longitude: 2.5,
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(data); got != `{"context":"c1","coordinates":{"latitude":1.5,"longitude":2.5}}` {
		t.Errorf("got %s", got)
	}
}

func TestEmulationSetGeolocationOverrideCleared(t *testing.T) {
	params := EmulationSetGeolocationOverrideParams{Context: "c1"}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(data); got != `{"context":"c1"}` {
		t.Errorf("got %s", got)
	}
}

func TestEmulationSetTimezoneOverrideParams(t *testing.T) {
	params := EmulationSetTimezoneOverrideParams{
		Context:  "c1",
		Timezone: "Europe/Berlin",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if got := string(data); got != `{"context":"c1","timezone":"Europe/Berlin"}` {
		t.Errorf("got %s", got)
	}
}
