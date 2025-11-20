package dto

type Team struct {
	TeamName string chan `json:"team_name"`
	Members  []TeamMember `json:"members"`
}
