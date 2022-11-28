package models

type Application struct {
	Name          string             `json:"name" validate:"required"`
	Description   string             `json:"description"`
	Configuration *ApplicationConfig `json:"config" validate:"dive"`
}

type ApplicationConfig struct {
	API        string               `json:"api"`
	Homepage   string               `json:"homepage"`
	IconName   string               `json:"icon"`
	Navigation []NavigationMenuItem `json:"navigation"`
}

type NavigationMenuItem struct {
	Label       string `json:"label"`
	Destination string `json:"destination"`
}
