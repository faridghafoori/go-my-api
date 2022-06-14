package server

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"
	"log"
	"net/http"
)

type PublicServer struct{}

func (ps PublicServer) RouterHandler() http.Handler {
	router := NewRouter()
	routes.EpisodeRoutes(router)
	return router
}

func (ps PublicServer) InitServer() {
	server := &http.Server{
		Addr:    configs.ENV_LAUNCH_PUBLIC_URI(),
		Handler: ps.RouterHandler(),
	}

	g.Go(func() error {
		log.Println("Starting public server at " + server.Addr)
		return server.ListenAndServe()
	})
}
