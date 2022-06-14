package server

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/middleware"
	"gin-mongo-api/routes"
	"log"
	"net/http"
)

type AdminServer struct{}

func (as AdminServer) RouterHandler() http.Handler {
	router := NewRouter()
	router.Use(middleware.TokenAuthMiddleware())
	routes.UserRoutes(router)
	routes.RoleRoutes(router)
	routes.EpisodePrivateRoutes(router)
	return router
}

func (as AdminServer) InitServer() {
	server := &http.Server{
		Addr:    configs.ENV_LAUNCH_ADMIN_URI(),
		Handler: as.RouterHandler(),
	}

	g.Go(func() error {
		log.Println("Starting admin server at " + server.Addr)
		return server.ListenAndServe()
	})
}
