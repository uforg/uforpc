package openapi

type Spec struct {
	OpenAPI    string     `json:"openapi"`
	Info       Info       `json:"info"`
	Servers    []Server   `json:"servers,omitempty"`
	Tags       []Tag      `json:"tags,omitempty"`
	Paths      Paths      `json:"paths,omitempty"`
	Components Components `json:"components,omitzero"`
}

type Info struct {
	Title       string `json:"title,omitzero"`
	Version     string `json:"version"`
	Description string `json:"description,omitzero"`
}

type Server struct {
	URL string `json:"url"`
}

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitzero"`
}

type Paths map[string]any

type Components struct {
	Schemas       map[string]any `json:"schemas,omitempty"`
	RequestBodies map[string]any `json:"requestBodies,omitempty"`
	Responses     map[string]any `json:"responses,omitempty"`
}
