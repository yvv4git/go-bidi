package browser

import (
	"context"

	"github.com/yvv4git/go-bidi/protocol"
)

// SetPermission grants, denies or prompts for a permission on the given
// origin. State is one of the permission state constants.
func (s *Session) SetPermission(
	ctx context.Context,
	params protocol.PermissionsSetPermissionParams,
) error {
	return s.caller.Call(ctx, protocol.PermissionsSetPermission, params, nil)
}
