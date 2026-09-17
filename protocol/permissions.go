package protocol

// PermissionsDescriptor identifies a permission by name and optional
// properties. UserVisibleOnly, Sysex and AllowWithoutPrompt are only
// meaningful for relevant permissions.
type PermissionsDescriptor struct {
	Name               string `json:"name"`
	UserVisibleOnly    *bool  `json:"userVisibleOnly,omitempty"`
	Sysex              *bool  `json:"sysex,omitempty"`
	AllowWithoutPrompt *bool  `json:"allowWithoutPrompt,omitempty"`
}

// PermissionsSetPermissionParams are the parameters of
// permissions.setPermission. State is one of the permission state constants.
type PermissionsSetPermissionParams struct {
	Descriptor PermissionsDescriptor `json:"descriptor"`
	State      string                `json:"state"`
	Origin     string                `json:"origin"`
}
