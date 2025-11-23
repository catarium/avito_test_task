package pullrequests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pullRequestRepository "github.com/catarium/avito_test_task/internal/db/repositories/pullrequest"
	"github.com/catarium/avito_test_task/internal/db/repositories/user"
	"github.com/catarium/avito_test_task/internal/schemas"
	"github.com/catarium/avito_test_task/internal/services"
	"github.com/catarium/avito_test_task/internal/services/pullrequest"
	"github.com/catarium/avito_test_task/internal/utils/httputils"
)

type PullRequestHandler struct {
	pullRequestService *pullrequest.PullRequestService
}

func (ph PullRequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	body := schemas.PullRequestSchema{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 400)
		return
	}
	pullRequest, errResp, code := ph.pullRequestService.Create(body.PullRequestId, body.PullRequestName, body.AuthorId)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(pullRequest)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(res))
}

func (ph PullRequestHandler) Merge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	body := schemas.PullRequestMergeSchema{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 400)
		return
	}
	pullRequest, errResp, code := ph.pullRequestService.Merge(body.PullRequestId)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(pullRequest)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(res))
}

func (ph PullRequestHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		errResp := services.ErrUnknown(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	body := schemas.PullRequestReassignSchema{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 400)
		return
	}
	pullRequest, errResp, code := ph.pullRequestService.Reassign(body.PullRequestId, body.OldReviewerId)
	if errResp != nil {
		httputils.SendJSONError(w, errResp, code)
		return
	}
	res, err := json.Marshal(pullRequest)
	if err != nil {
		errResp := services.ErrInvalidJson(err.Error())
		httputils.SendJSONError(w, errResp, 500)
		return
	}
	w.WriteHeader(code)
	fmt.Fprintf(w, string(res))
}

func CreatePullRequestRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	pullRequestHandler := PullRequestHandler{
		pullRequestService: &pullrequest.PullRequestService{
			PullRequestRepository: &pullRequestRepository.PullRequestRepository{DB: db},
			UserRepository:        &user.UserRepository{DB: db},
		},
	}
	mux.HandleFunc("POST /create", pullRequestHandler.Create)
	mux.HandleFunc("POST /merge", pullRequestHandler.Merge)
	mux.HandleFunc("POST /reassign", pullRequestHandler.Reassign)
	return mux
}
