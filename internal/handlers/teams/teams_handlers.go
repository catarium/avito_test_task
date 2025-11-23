package teams

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/catarium/avito_test_task/internal/db/repositories/team"
	"github.com/catarium/avito_test_task/internal/db/repositories/user"
	"github.com/catarium/avito_test_task/internal/dto"
	"github.com/catarium/avito_test_task/internal/services"
	teamService "github.com/catarium/avito_test_task/internal/services/team"
	"github.com/catarium/avito_test_task/internal/utils/httputils"
)

type TeamHandler struct {
	teamService *teamService.TeamService
}

func (th TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	team := dto.Team{}
	err = json.Unmarshal(bodyBytes, &team)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 400)
		return
	}
	_, errResp, code := th.teamService.AddTeam(team.TeamName, team.Members)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(team)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(res))
}

func (th TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	teamName := queryParams.Get("team_name")
	team, errResp, code := th.teamService.GetTeam(teamName)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(team)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(res))
}

func CreateTeamRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	teamHandler := TeamHandler{teamService: &teamService.TeamService{
		TeamRepository: &team.TeamRepository{DB: db},
		UserRepository: &user.UserRepository{DB: db},
	}}
	mux.HandleFunc("POST /add", teamHandler.AddTeam)
	mux.HandleFunc("GET /get", teamHandler.GetTeam)
	return mux
}
