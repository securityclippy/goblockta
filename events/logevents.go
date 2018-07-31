package events


type  LogEvent struct {
	UUID string `json:"uuid"`
	Published string `json:"published"`
	EventType string `json:"eventType"`
	Version string `json:"version"`
	Severity string `json:"severity"`
	LegacyEventType string `json:"legacyEventType"`
	DisplayMessage string `json:"displayMessage"`
	Actor Actor `json:"actor"`
	Client Client `json:"client"`
	Outcome Outcome `json:"outcome"`
	Target []Target `json:"target"`
	Transaction Transaction `json:"transaction"`
	DebugContext DebugContext `json:"debugContext"`
	AuthenticationContext Authentication `json:"authenticationContext"`
	SecurityContext SecurityContext `json:"securityContext"`
}

type Actor struct {
	ID string `json:"id"`
	Type string `json:"type"`
	AlternateID string `json:"alternateId"`
	DisplayName string `json:"displayName"`
	Detail interface{} `json:"detail"`
}

type Target struct {
	ID string `json:"id"`
	Type string `json:"type"`
	AlternateID string `json:"alternateId"`
	DisplayName string `json:"displayName"`
	Detail interface{} `json:"detail"`
}

type Client struct {
	ID interface{} `json:"id"`
	UserAgent UserAgent `json:"userAgent"`
	GeographicalContext GeographicalContext `json:"geographicalContext"`
	Zone string `json:"zone"`
	IPAddress string `json:"ipAddress"`
	Device string `json:"device"`
}

type UserAgent struct {
	Browser string `json:"Browser"`
	OS string `json:"OS"`
	RawUserAgent string `json:"RawUserAgent"`
}

type Request struct {
	IPChain []IPAddress `json:"ipChain"`
}


type GeographicalContext struct {
	Geolocation Geolocation `json:"geolocation"`
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type Geolocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`

}

type Outcome struct {
	Result string `json:"result"`
	Reason string `json:"reason"`
}

type Transaction struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Detail interface{} `json:"detail"`
}

type DebugContext struct {
	DebugData interface{} `json:"debugData"`
}

type Authentication struct {
	AuthenticationProvider string `json:"authenticationProvider"`
	CredentialProvider string `json:"credentialProvider"`
	CredentialType string `json:"credentialType"`
	Issuer Issuer `json:"issuer"`
	ExternalSessionID string `json:"externalSessionId"`
	Interface string `json:"interface"`
}

type Issuer struct {
	ID string `json:"id"`
	Type string `json:"type"`
}

type SecurityContext struct {
	ASNumber int `json:"asNumber"`
	ASOrg string `json:"asOrg"`
	ISP string `json:"isp"`
	Domain string `json:"domain"`
	ISProxy bool `json:"isProxy"`
}



type IPAddress struct {
	IP string `json:"ip"`
	GeographicalContext GeographicalContext `json:"geographicalContext"`
	Version string `json:"version"`
	Source string `json:"source"`
}