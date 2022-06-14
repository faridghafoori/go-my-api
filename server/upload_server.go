package server

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/middleware"
	"gin-mongo-api/routes"
	"log"
	"net/http"
)

type UploadServer struct{}

func (us UploadServer) RouterHandler() http.Handler {
	router := NewRouter()
	router.Use(middleware.TokenAuthMiddleware())
	routes.IORoutes(router)
	return router
}

func (us UploadServer) InitServer() {
	server := &http.Server{
		Addr:    configs.ENV_LAUNCH_UPLOAD_URI(),
		Handler: us.RouterHandler(),
	}

	g.Go(func() error {
		log.Println("Starting upload server at " + server.Addr)
		return server.ListenAndServe()
	})
}
