package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/constants"
	"github.com/iota-agency/iota-sdk/pkg/services"
)

func Tabs() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				u, err := composables.UseUser(r.Context())
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				app, err := composables.UseApp(r.Context())
				if err != nil {
					panic(err)
				}
				tabService := app.Service(services.TabService{}).(*services.TabService)
				tabs, err := tabService.GetUserTabs(r.Context(), u.ID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				ctx := context.WithValue(r.Context(), constants.TabsKey, tabs)
				logger, err := composables.UseLogger(r.Context())
				if err == nil {
					logger.WithField("duration", time.Since(start)).Info("middleware.Tabs")
				}
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
