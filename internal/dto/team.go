package dto

type TeamDto struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type TeamCreateDto struct {
	Team TeamDto `json:"team"`
}
