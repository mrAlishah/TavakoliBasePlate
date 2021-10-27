package main

import (
	"GolangTraining/internal/http/rest"
	"GolangTraining/internal/metrics"

	"fmt"
	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
)

func SetupRouter(handler *rest.Handler, cfg *MainConfig, p *metrics.Prometheus) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, rest.NotFound)
	})

	r.GET("/health", handler.Health)

	v1 := r.Group("/v1")
	if cfg.Server.AuthEnabled {
		v1.Use(gin.BasicAuth(gin.Accounts{
			cfg.Server.User: cfg.Server.Pass,
		}))
	}

	{
		v1.GET("/browse", handler.Health)
		v1.GET("/sigin", handler.SignIn)
		//v1.GET("/msg", handler.GetMessages)
	}

	p.MetricsPath = fmt.Sprintf("/%s", "metrics")
	if cfg.Server.AuthEnabled {
		p.UseWithAuth(r, gin.Accounts{
			cfg.Server.User: cfg.Server.Pass,
		})
	} else {
		p.Use(r)
	}

	var AllowedRoutes = make(map[string]bool, 0)
	routes := r.Routes()
	for _, i := range routes {
		AllowedRoutes[i.Path] = true
	}

	paramStripMap := make(map[string]bool, 0)

	for _, sp := range []string{"q", "username"} {
		paramStripMap[sp] = true
	}

	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if _, ok := paramStripMap[p.Key]; ok {
				// Found param
				url = strings.Replace(url, p.Value, fmt.Sprintf(":%s", p.Key), 1)
			}
		}
		if _, ok := AllowedRoutes[url]; ok {
			return url
		} else {
			return ""
		}
	}
	return r
}
