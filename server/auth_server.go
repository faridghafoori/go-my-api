package server

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"
	"log"
	"net/http"
)

type AuthServer struct{}

func (as AuthServer) RouterHandler() http.Handler {
	router := NewRouter()
	routes.AuthenticationRoutes(router)
	return router
}

func (as AuthServer) InitServer() {
	server := &http.Server{
		Addr:    configs.ENV_LAUNCH_AUTH_URI(),
		Handler: as.RouterHandler(),
	}

	g.Go(func() error {
		log.Println("Starting auth server at " + server.Addr)
		return server.ListenAndServe()
	})
}
