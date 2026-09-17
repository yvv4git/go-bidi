package browser

import (
	"context"

	"github.com/yvv4git/go-bidi/protocol"
)

// SetGeolocationOverride overrides the geolocation of a browsing context.
// The override is cleared when Coordinates is nil.
func (s *Session) SetGeolocationOverride(
	ctx context.Context,
	params protocol.EmulationSetGeolocationOverrideParams,
) error {
	return s.caller.Call(ctx, protocol.EmulationSetGeolocationOverride, params, nil)
}

// SetTimezoneOverride overrides the timezone of a browsing context.
func (s *Session) SetTimezoneOverride(
	ctx context.Context,
	params protocol.EmulationSetTimezoneOverrideParams,
) error {
	return s.caller.Call(ctx, protocol.EmulationSetTimezoneOverride, params, nil)
}
