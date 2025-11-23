package team

import (
	"database/sql"
	"net/http"

	"github.com/catarium/avito_test_task/internal/db/repositories/team"
	"github.com/catarium/avito_test_task/internal/db/repositories/user"
	"github.com/catarium/avito_test_task/internal/dto"
	"github.com/catarium/avito_test_task/internal/services"
)

type TeamService struct {
	TeamRepository *team.TeamRepository
	UserRepository *user.UserRepository
}

func (ts TeamService) AddTeam(teamName string, members []dto.TeamMember) (*dto.Team, *dto.ErrorDto, int) {
	res := dto.Team{TeamName: teamName}
	exists, err := ts.TeamRepository.Exists(teamName)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if exists {
		return nil, ErrTeamExists(teamName), http.StatusBadRequest
	}
	_, err = ts.TeamRepository.CreateTeam(teamName)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	for _, m := range members {
		exists, err := ts.UserRepository.Exists(m.UserId)
		if err != nil {
			return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
		}
		if exists {
			_, err := ts.UserRepository.Update(m.UserId, m.UserName, teamName, m.IsActive)
			if err != nil {
				return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
			}
		} else {
			_, err := ts.UserRepository.Create(m.UserId, m.UserName, teamName, m.IsActive)
			if err != nil {
				return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
			}
		}
		res.Members = append(res.Members, m)
	}
	return &res, nil, http.StatusCreated
}

func (ts TeamService) GetTeam(teamName string) (*dto.Team, *dto.ErrorDto, int) {
	res := dto.Team{TeamName: teamName}
	exists, err := ts.TeamRepository.Exists(teamName)
	if err != nil {
		return nil, services.ErrUnknown(err.Error()), http.StatusInternalServerError
	}
	if !exists {
		return nil, &services.ErrNotFound, http.StatusNotFound
	}
	members, err := ts.UserRepository.GetByTeamName(teamName)
	if err == sql.ErrNoRows {
		return &res, nil, http.StatusOK
	}
	var member dto.TeamMember
	for _, m := range members {
		member.UserId = m.UserId
		member.UserName = m.Username
		member.IsActive = m.IsActive
		res.Members = append(res.Members, member)
	}
	return &res, nil, http.StatusOK
}
