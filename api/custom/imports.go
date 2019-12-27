package custom

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-media-repo/api"
	"github.com/turt2live/matrix-media-repo/common/config"
	"github.com/turt2live/matrix-media-repo/controllers/data_controller"
)

type ImportStarted struct {
	ImportID string `json:"import_id"`
	TaskID   int    `json:"task_id"`
}

func StartImport(r *http.Request, log *logrus.Entry, user api.UserInfo) interface{} {
	if !config.Get().Archiving.Enabled {
		return api.BadRequest("archiving is not enabled")
	}

	defer r.Body.Close()
	task, importId, err := data_controller.StartImport(r.Body, log)
	if err != nil {
		log.Error(err)
		return api.InternalServerError("fatal error starting import")
	}

	return &api.DoNotCacheResponse{Payload: &ImportStarted{
		TaskID:   task.ID,
		ImportID: importId,
	}}
}

func AppendToImport(r *http.Request, log *logrus.Entry, user api.UserInfo) interface{} {
	if !config.Get().Archiving.Enabled {
		return api.BadRequest("archiving is not enabled")
	}

	params := mux.Vars(r)

	importId := params["importId"]

	defer r.Body.Close()
	err := data_controller.AppendToImport(importId, r.Body)
	if err != nil {
		log.Error(err)
		return api.InternalServerError("fatal error appending to import")
	}

	return &api.DoNotCacheResponse{Payload: &api.EmptyResponse{}}
}

func StopImport(r *http.Request, log *logrus.Entry, user api.UserInfo) interface{} {
	if !config.Get().Archiving.Enabled {
		return api.BadRequest("archiving is not enabled")
	}

	params := mux.Vars(r)

	importId := params["importId"]

	err := data_controller.StopImport(importId)
	if err != nil {
		log.Error(err)
		return api.InternalServerError("fatal error stopping import")
	}

	return &api.DoNotCacheResponse{Payload: &api.EmptyResponse{}}
}