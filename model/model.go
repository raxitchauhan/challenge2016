package model

type Distributor struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	PermittedPlaces   []string `json:"permitted_places"`
	HasSubDistributor bool     `json:"has_sub_distributor"`
	Parent            string   `json:"parent"`
	Child             []string `json:"child"`
}
