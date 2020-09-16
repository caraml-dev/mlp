package models

type Application struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Href        string `json:"href"`
	IconName    string `json:"icon" gorm:"column:icon"`
	UseProjects bool   `json:"use_projects"`
	IsInBeta    bool   `json:"is_in_beta"`
	IsDisabled  bool   `json:"is_disabled"`
}
