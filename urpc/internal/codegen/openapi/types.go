package openapi

type Spec struct {
	OpenAPI string   `json:"openapi"`
	Info    Info     `json:"info"`
	Servers []Server `json:"servers,omitempty"`
}

type Info struct {
	Title       string `json:"title,omitzero"`
	Description string `json:"description,omitzero"`
	Version     string `json:"version"`
}

type Server struct {
	URL string `json:"url"`
}
