package models

type Application struct {
	Name                           string                          `json:"name" validate:"required"`
	Description                    string                          `json:"description"`
	Homepage                       string                          `json:"homepage"`
	Configuration                  *ApplicationConfig              `json:"config" validate:"dive"`
	IsProjectAgnostic              bool                            `json:"is_project_agnostic"`
	StreamlitPlaceholderPageConfig *StreamlitPlaceholderPageConfig `json:"streamlit_placeholder_page_config"`
}

type ApplicationConfig struct {
	API        string               `json:"api"`
	IconName   string               `json:"icon"`
	Navigation []NavigationMenuItem `json:"navigation"`
}

type NavigationMenuItem struct {
	Label       string `json:"label"`
	Destination string `json:"destination"`
}

type StreamlitPlaceholderPageConfig struct {
	StreamlitURL string `json:"streamlit_url"`
}
