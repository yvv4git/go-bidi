package protocol

// EmulationCoordinates are the geolocation coordinates used by the emulation
// module. Accuracy is optional and reset when nil.
type EmulationCoordinates struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Accuracy  *float64 `json:"accuracy,omitempty"`
}

// EmulationSetGeolocationOverrideParams are the parameters of
// emulation.setGeolocationOverride. The override is cleared when Coordinates
// is nil.
type EmulationSetGeolocationOverrideParams struct {
	Context     string                `json:"context"`
	Coordinates *EmulationCoordinates `json:"coordinates,omitempty"`
}

// EmulationSetTimezoneOverrideParams are the parameters of
// emulation.setTimezoneOverride.
type EmulationSetTimezoneOverrideParams struct {
	Context  string `json:"context"`
	Timezone string `json:"timezone"`
}
