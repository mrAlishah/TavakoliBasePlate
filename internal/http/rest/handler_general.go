package rest

import (
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"time"
)

func (h *Handler) Health(c *gin.Context) {
	uptime := sigar.Uptime{}
	uptime.Get()
	avg := sigar.LoadAverage{}
	avg.Get()
	hcDTO := HealthCheckResponse{
		Status:       "ok",
		Container:    h.VersionInfo.ContainerName,
		GitCommit:    h.VersionInfo.GitCommit,
		BuildTime:    h.VersionInfo.BuildTime,
		Version:      runtime.Version(),
		Uptime:       uptime.Format(),
		BinaryUptime: time.Since(h.VersionInfo.BinaryStartTime).String(),
		LAOne:        avg.One,
		LAFive:       avg.Five,
		LAFifteen:    avg.Fifteen,
	}
	c.JSON(http.StatusOK, GetSuccessResponse(hcDTO.Status))
}

type HealthCheckResponse struct {
	Status       string  `json:"status"`
	Container    string  `json:"container"`
	GitCommit    string  `json:"git_commit"`
	BuildTime    string  `json:"build_time"`
	Version      string  `json:"go_version"`
	Uptime       string  `json:"uptime"`
	BinaryUptime string  `json:"binary_uptime"`
	LAOne        float64 `json:"load_average_one"`
	LAFive       float64 `json:"load_average_five"`
	LAFifteen    float64 `json:"load_average_fifteen"`
}
