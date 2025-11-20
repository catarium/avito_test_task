package dto

type TeamMember struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	IsActive bool   `json:"is_active"`
}
