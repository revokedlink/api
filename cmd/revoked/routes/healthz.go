package routes

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

func HealthzRoute(app core.App) {
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {

		e.Router.GET("/healthz", func(e *core.RequestEvent) error {
			return e.String(http.StatusOK, "ok")
		})

		return e.Next()
	})
}
