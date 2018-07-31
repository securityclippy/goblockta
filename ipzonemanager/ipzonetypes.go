package ipzonemanager

type IPZone struct {
	Type string `json:"type"`
	ID string `json:"id"`
	Name string `json:"name"`
	Gateways []IPAddress `json:"gateways"`
	Proxies []IPAddress `json:"proxies"`
}

type IPAddress struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

