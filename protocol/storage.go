package protocol

// StorageCookie is a cookie stored in the browser. SameSite is one of the
// cookie same site constants.
type StorageCookie struct {
	Name     string            `json:"name"`
	Value    NetworkBytesValue `json:"value"`
	Domain   string            `json:"domain"`
	Path     string            `json:"path"`
	Size     int64             `json:"size"`
	Expires  int64             `json:"expires"`
	HTTPOnly bool              `json:"httpOnly"`
	Secure   bool              `json:"secure"`
	SameSite string            `json:"sameSite,omitempty"`
}

// StorageCookieFilter narrows the cookies selected by storage.getCookies and
// storage.deleteCookies. Nil fields match any value.
type StorageCookieFilter struct {
	Name     *string `json:"name,omitempty"`
	Value    *string `json:"value,omitempty"`
	Domain   *string `json:"domain,omitempty"`
	Path     *string `json:"path,omitempty"`
	HTTPOnly *bool   `json:"httpOnly,omitempty"`
	Secure   *bool   `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
}

// StorageGetCookiesParams are the parameters of storage.getCookies.
type StorageGetCookiesParams struct {
	Filter *StorageCookieFilter `json:"filter,omitempty"`
}

// StorageGetCookiesResult is the result of storage.getCookies.
type StorageGetCookiesResult struct {
	Cookies []StorageCookie `json:"cookies"`
}

// StorageSetCookieParams are the parameters of storage.setCookie.
type StorageSetCookieParams struct {
	Cookie StorageCookie `json:"cookie"`
}

// StorageDeleteCookiesParams are the parameters of storage.deleteCookies.
type StorageDeleteCookiesParams struct {
	Filter *StorageCookieFilter `json:"filter,omitempty"`
}
