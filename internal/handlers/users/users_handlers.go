package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pullRequestRepository "github.com/catarium/avito_test_task/internal/db/repositories/pullrequest"
	userRepository "github.com/catarium/avito_test_task/internal/db/repositories/user"
	"github.com/catarium/avito_test_task/internal/schemas"
	"github.com/catarium/avito_test_task/internal/services"
	"github.com/catarium/avito_test_task/internal/services/user"
	"github.com/catarium/avito_test_task/internal/utils/httputils"
)

type UserHandler struct {
	userService user.UserService
}

func (uh UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	body := schemas.UserIsActiveSchema{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 400)
		return
	}
	user, errResp, code := uh.userService.SetActive(body.UserId, body.IsActive)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(user)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(res))
}

func (uh UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userId := queryParams.Get("user_id")
	userPr, errResp, code := uh.userService.GetReviewedPullRequestsByUserId(userId)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(userPr)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(res))
}

func CreateUserRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	userHandler := UserHandler{
		userService: user.UserService{
			UserRepository:        userRepository.UserRepository{DB: db},
			PullRequestRepository: pullRequestRepository.PullRequestRepository{DB: db},
		},
	}
	mux.HandleFunc("POST /setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("GET /getReview", userHandler.GetReview)
	return mux
}
