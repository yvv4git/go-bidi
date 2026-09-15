package protocol

// NetworkBytesValue is a byte string used in headers, cookies and request
// bodies. Type is one of the bytes value kind constants.
type NetworkBytesValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NetworkHeader is a single HTTP header. Value is encoded as a bytes value.
type NetworkHeader struct {
	Name  string            `json:"name"`
	Value NetworkBytesValue `json:"value"`
}

// NetworkCookie describes a cookie sent with a request or received in a
// response.
type NetworkCookie struct {
	Name     string            `json:"name"`
	Value    NetworkBytesValue `json:"value"`
	Domain   string            `json:"domain,omitempty"`
	Path     string            `json:"path,omitempty"`
	HTTPOnly bool              `json:"httpOnly,omitempty"`
	Secure   bool              `json:"secure,omitempty"`
	SameSite string            `json:"sameSite,omitempty"`
	Size     int64             `json:"size,omitempty"`
}

// NetworkTimingInfo reports the timing of the network phases of a request in
// milliseconds. Nullable fields are nil when the browser did not report them.
type NetworkTimingInfo struct {
	OriginTime    int64  `json:"originTime"`
	RequestTime   int64  `json:"requestTime"`
	RedirectStart *int64 `json:"redirectStart"`
	RedirectEnd   *int64 `json:"redirectEnd"`
	FetchStart    *int64 `json:"fetchStart"`
	DNSStart      *int64 `json:"dnsStart"`
	DNSEnd        *int64 `json:"dnsEnd"`
	ConnectStart  *int64 `json:"connectStart"`
	ConnectEnd    *int64 `json:"connectEnd"`
	TLSStart      *int64 `json:"tlsStart"`
	RequestStart  *int64 `json:"requestStart"`
	ResponseStart *int64 `json:"responseStart"`
	ResponseEnd   *int64 `json:"responseEnd"`
}

// NetworkInitiator describes what started a request.
type NetworkInitiator struct {
	Type    string `json:"type"`
	Request string `json:"request,omitempty"`
}

// NetworkRequestData describes a request sent by the browser.
type NetworkRequestData struct {
	Request     string            `json:"request"`
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     []NetworkHeader   `json:"headers"`
	Cookies     []NetworkCookie   `json:"cookies"`
	HeadersSize int64             `json:"headersSize"`
	BodySize    *int64            `json:"bodySize"`
	Destination string            `json:"destination"`
	Initiator   NetworkInitiator  `json:"initiator"`
	Timing      NetworkTimingInfo `json:"timing"`
}

// NetworkResponseData describes the response received for a request.
type NetworkResponseData struct {
	URL         string          `json:"url"`
	Protocol    string          `json:"protocol"`
	Status      int             `json:"status"`
	StatusText  string          `json:"statusText"`
	FromCache   bool            `json:"fromCache"`
	Headers     []NetworkHeader `json:"headers"`
	MimeType    string          `json:"mimeType"`
	KnownLength *int64          `json:"knownLength"`
	BodySize    *int64          `json:"bodySize"`
	Cookies     []NetworkCookie `json:"cookies"`
}

// NetworkBeforeRequestSentParams is the payload of network.beforeRequestSent.
type NetworkBeforeRequestSentParams struct {
	Context       string             `json:"context"`
	IsBlocked     bool               `json:"isBlocked"`
	Navigation    *string            `json:"navigation"`
	RedirectCount int64              `json:"redirectCount"`
	Request       NetworkRequestData `json:"request"`
	Timestamp     int64              `json:"timestamp"`
}

// NetworkResponseStartedParams is the payload of network.responseStarted.
type NetworkResponseStartedParams struct {
	Context       string              `json:"context"`
	IsBlocked     bool                `json:"isBlocked"`
	Navigation    *string             `json:"navigation"`
	RedirectCount int64               `json:"redirectCount"`
	Request       NetworkRequestData  `json:"request"`
	Response      NetworkResponseData `json:"response"`
	Timestamp     int64               `json:"timestamp"`
}

// NetworkResponseCompletedParams is the payload of network.responseCompleted.
type NetworkResponseCompletedParams struct {
	Context       string              `json:"context"`
	IsBlocked     bool                `json:"isBlocked"`
	Navigation    *string             `json:"navigation"`
	RedirectCount int64               `json:"redirectCount"`
	Request       NetworkRequestData  `json:"request"`
	Response      NetworkResponseData `json:"response"`
	Timestamp     int64               `json:"timestamp"`
}

// NetworkFetchErrorParams is the payload of network.fetchError.
type NetworkFetchErrorParams struct {
	Context       string             `json:"context"`
	IsBlocked     bool               `json:"isBlocked"`
	Navigation    *string            `json:"navigation"`
	RedirectCount int64              `json:"redirectCount"`
	Request       NetworkRequestData `json:"request"`
	Timestamp     int64              `json:"timestamp"`
	ErrorText     string             `json:"errorText"`
}

// NetworkAuthRequiredParams is the payload of network.authRequired.
type NetworkAuthRequiredParams struct {
	Context       string              `json:"context"`
	IsBlocked     bool                `json:"isBlocked"`
	Navigation    *string             `json:"navigation"`
	RedirectCount int64               `json:"redirectCount"`
	Request       NetworkRequestData  `json:"request"`
	Response      NetworkResponseData `json:"response"`
	Timestamp     int64               `json:"timestamp"`
}

// NetworkContinueRequestParams are the parameters of network.continueRequest.
type NetworkContinueRequestParams struct {
	Request string             `json:"request"`
	Body    *NetworkBytesValue `json:"body,omitempty"`
	Cookies []NetworkCookie    `json:"cookies,omitempty"`
	Headers []NetworkHeader    `json:"headers,omitempty"`
	Method  string             `json:"method,omitempty"`
	URL     string             `json:"url,omitempty"`
}

// NetworkContinueResponseParams are the parameters of
// network.continueResponse.
type NetworkContinueResponseParams struct {
	Request      string                  `json:"request"`
	Cookies      []NetworkCookie         `json:"cookies,omitempty"`
	Credentials  *NetworkAuthCredentials `json:"credentials,omitempty"`
	Headers      []NetworkHeader         `json:"headers,omitempty"`
	ReasonPhrase string                  `json:"reasonPhrase,omitempty"`
	StatusCode   *int                    `json:"statusCode,omitempty"`
}

// NetworkContinueWithAuthParams are the parameters of network.continueWithAuth.
// Action is one of the auth action constants.
type NetworkContinueWithAuthParams struct {
	Request     string                  `json:"request"`
	Action      string                  `json:"action"`
	Credentials *NetworkAuthCredentials `json:"credentials,omitempty"`
}

// NetworkAuthCredentials carries credentials for network.continueWithAuth.
// Type is one of the auth credentials type constants.
type NetworkAuthCredentials struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NetworkFailRequestParams are the parameters of network.failRequest.
type NetworkFailRequestParams struct {
	Request string `json:"request"`
}

// NetworkProvideResponseParams are the parameters of network.provideResponse.
type NetworkProvideResponseParams struct {
	Request      string             `json:"request"`
	Body         *NetworkBytesValue `json:"body,omitempty"`
	Cookies      []NetworkCookie    `json:"cookies,omitempty"`
	Headers      []NetworkHeader    `json:"headers,omitempty"`
	ReasonPhrase string             `json:"reasonPhrase,omitempty"`
	StatusCode   *int               `json:"statusCode,omitempty"`
}

// NetworkURLPattern matches request URLs. Type is one of the url pattern
// constants; Hostname, Port, Path and Search are matched with glob syntax.
type NetworkURLPattern struct {
	Type     string `json:"type"`
	Protocol string `json:"protocol,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Port     string `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
	Search   string `json:"search,omitempty"`
}

// NetworkAddInterceptParams are the parameters of network.addIntercept.
// Phases holds one or more of the intercept phase constants.
type NetworkAddInterceptParams struct {
	Phases      []string            `json:"phases"`
	Contexts    []string            `json:"contexts,omitempty"`
	URLPatterns []NetworkURLPattern `json:"urlPatterns,omitempty"`
}

// NetworkAddInterceptResult is the result of network.addIntercept.
type NetworkAddInterceptResult struct {
	Intercept string `json:"intercept"`
}

// NetworkRemoveInterceptParams are the parameters of network.removeIntercept.
type NetworkRemoveInterceptParams struct {
	Intercept string `json:"intercept"`
}

// NetworkSetCacheBehaviorParams are the parameters of
// network.setCacheBehavior. CacheBehavior is one of the cache behavior
// constants.
type NetworkSetCacheBehaviorParams struct {
	CacheBehavior string   `json:"cacheBehavior"`
	Contexts      []string `json:"contexts,omitempty"`
}
