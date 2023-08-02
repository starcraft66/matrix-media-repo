package rcontext

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/turt2live/matrix-media-repo/common/config"
)

func Initial() RequestContext {
	return RequestContext{
		Context: context.Background(),
		Log:     logrus.WithFields(logrus.Fields{"internal_flag": 1}),
		Config: config.DomainRepoConfig{
			MinimumRepoConfig: config.Get().MinimumRepoConfig,
			Downloads:         config.Get().Downloads.DownloadsConfig,
			Thumbnails:        config.Get().Thumbnails.ThumbnailsConfig,
			UrlPreviews:       config.Get().UrlPreviews.UrlPreviewsConfig,
		},
		Request: nil,
	}.populate()
}

func InitialNoConfig() RequestContext {
	return RequestContext{
		Context: context.Background(),
		Log:     logrus.WithFields(logrus.Fields{"internal_flag": 2}),
		Config:  config.DomainRepoConfig{},
		Request: nil,
	}.populate()
}

type RequestContext struct {
	context.Context

	// These are also stored on the context object itself
	Log     *logrus.Entry           // mmr.logger
	Config  config.DomainRepoConfig // mmr.serverConfig
	Request *http.Request           // mmr.request
}

func (c RequestContext) populate() RequestContext {
	c.Context = context.WithValue(c.Context, "mmr.logger", c.Log)
	c.Context = context.WithValue(c.Context, "mmr.serverConfig", c.Config)
	c.Context = context.WithValue(c.Context, "mmr.request", c.Request)
	return c
}

func (c RequestContext) ReplaceLogger(log *logrus.Entry) RequestContext {
	ctx := context.WithValue(c.Context, "mmr.logger", log)
	return RequestContext{
		Context: ctx,
		Log:     log,
		Config:  c.Config,
		Request: c.Request,
	}
}

func (c RequestContext) LogWithFields(fields logrus.Fields) RequestContext {
	return c.ReplaceLogger(c.Log.WithFields(fields))
}
