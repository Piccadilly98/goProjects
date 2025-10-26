package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/rate_limit"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
)

type AdminControlGlobalRl struct {
	stor     *storage.Storage
	globalRl *rate_limit.GlobalRateLimiter
	mtx      sync.Mutex
}

const (
	actionStop         = "stop"
	actionRestart      = "restart"
	actionStopAttacked = "attacked"
	actionPatchLimited = "changeLimits"
)

func NewAdminControlRL(stor *storage.Storage,
	global *rate_limit.GlobalRateLimiter) *AdminControlGlobalRl {

	return &AdminControlGlobalRl{
		stor:     stor,
		globalRl: global,
	}
}

func (a *AdminControlGlobalRl) ControlRlHandler(w http.ResponseWriter, r *http.Request) {
	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims, ok := r.Context().Value("jwtClaims").(*models.ClaimsJSON)
	if !ok {
		errors = "server error"
		w.WriteHeader(http.StatusInternalServerError)
		a.stor.NewLog(r, nil, httpCode, errors)
		w.Write([]byte(errors))
		return
	}
	defer func() {
		a.stor.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()
	param := chi.URLParam(r, "action")
	if param == actionStop {
		fmt.Println("stop")
		a.globalRl.StopRefillToken(false)
		response := `{"success": "true"}`
		attentions = append(attentions, "GLOBAL RL IS STOPED!!")
		w.Write([]byte(response))
		return
	} else if param == actionStopAttacked {
		fmt.Println("attacked")
		a.globalRl.StopRefillToken(true)
		response := `{"success": "true"}`
		attentions = append(attentions, "GLOBAL RL IS STOPED-ATTACKED!!")
		w.Write([]byte(response))
		return
	} else if param == actionRestart {
		fmt.Println("restart")
		if a.globalRl.GetAttackedStatus() {
			key := r.Header.Get("key")
			if key != os.Getenv("GLOBAL_RL_RESTART") {
				errors = "invalid restart key"
				httpCode = http.StatusBadRequest
				attentions = append(attentions, "INVALID RESTAT KEY!")
				w.WriteHeader(httpCode)
				responce := `{"success": "false","error":"ivalid restart key"}`
				w.Write([]byte(responce))
				return
			}
		}
		a.globalRl.Restart()
		response := `{"success": "true"}`
		attentions = append(attentions, "GLOBAL RL IS START!!")
		w.Write([]byte(response))
		return
	} else if param == actionPatchLimited {
		bodyJson := &models.PathGlobalModelsJSON{}
		b, err := io.ReadAll(r.Body)
		if err != nil {
			errors = err.Error()
			httpCode = http.StatusInternalServerError
			w.WriteHeader(httpCode)
			response := `{"success": "false","error":"server error"}`
			w.Write([]byte(response))
			return
		}

		err = json.Unmarshal(b, bodyJson)
		if err != nil {
			errors = err.Error()
			httpCode = http.StatusBadRequest
			w.WriteHeader(httpCode)
			response := `{"success": "false","error":"invalid body"}`
			w.Write([]byte(response))
			return
		}
		body, err := bodyJson.ToIntegerStruct()
		if err != nil {
			errors = err.Error()
			httpCode = http.StatusBadRequest
			w.WriteHeader(httpCode)
			response := `{"success": "false","error":"invalid body"}`
			w.Write([]byte(response))
			return
		}

		if !body.Validate() {
			httpCode = http.StatusBadRequest
			errors = "invalid body"
			w.WriteHeader(httpCode)
			response := `{"success": "false","error":"invalid body"}`
			w.Write([]byte(response))
			return
		}
		if a.globalRl.GetAttackedStatus() {
			key := r.Header.Get("key")
			if key != os.Getenv("GLOBAL_RL_RESTART") {
				errors = "invalid restart key"
				httpCode = http.StatusBadRequest
				attentions = append(attentions, "INVALID RESTAT KEY!")
				w.WriteHeader(httpCode)
				responce := `{"success": "false","error":"ivalid restart key"}`
				w.Write([]byte(responce))
				return
			}
		}
		a.globalRl.ChangeLimits(body.ReqInSecond, body.StartTokens)
		w.WriteHeader(httpCode)
		response := `{"success": "true"}`
		w.Write([]byte(response))
		attentions = append(attentions, fmt.Sprintf("CHANGE LIMITS IN GLOBAL RL: reqInSecond: %d  startTokens: %d", body.ReqInSecond, body.StartTokens))
		return

	} else {
		fmt.Println("hz")
		httpCode = http.StatusBadRequest
		errors = "Invalid action"
		w.WriteHeader(httpCode)
		response := `{"success": "false"}`
		w.Write([]byte(response))
		return
	}
}
