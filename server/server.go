package server

import (
	"log"

	"golang.org/x/sync/errgroup"
)

// type ServerFunctialities interface {
// 	AuthRouterHandler() http.Handler
// 	InitServer(http.Server)
// }

// type ServerStruct struct {
// 	server http.Server
// }

// func (s ServerStruct) AuthRouterHandler() http.Handler {
// 	router := NewRouter()
// 	routes.AuthenticationRoutes(router)
// 	return router
// }

// func (s ServerStruct) InitServer(serverConfig http.Server) {
// 	server := &serverConfig

// 	g.Go(func() error {
// 		log.Println("Starting auth server at " + server.Addr)
// 		return server.ListenAndServe()
// 	})
// }

// func ServerInitializer(sf ServerFunctialities) {
// 	// sf.InitServer()
// 	log.Panic(sf)
// }

var g errgroup.Group

func Init() {
	AuthServer.InitServer(AuthServer{})
	PublicServer.InitServer(PublicServer{})
	AdminServer.InitServer(AdminServer{})
	UploadServer.InitServer(UploadServer{})
	// authServer := &ServerStruct{}
	// authServer.server.Addr = configs.ENV_LAUNCH_AUTH_URI()
	// authServer.server.Handler = authServer.AuthRouterHandler()
	// ServerInitializer(authServer)

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

// func Init() {
// 	InitAuthServer()
// 	InitAdminServer()
// 	InitUploadServer()
// 	InitPublicServer()
// }

// func InitAdminServer() {
// 	router := NewRouter()
// router.Use(middleware.TokenAuthMiddleware())
// routes.UserRoutes(router)
// routes.RoleRoutes(router)
// 	router.Run(configs.ENV_LAUNCH_ADMIN_URI())
// }

// func InitUploadServer() {
// 	router := NewRouter()
// 	router.Use(middleware.TokenAuthMiddleware())
// 	routes.IORoutes(router)
// 	router.Run(configs.ENV_LAUNCH_UPLOAD_URI())
// }

// func InitPublicServer() {
// 	router := NewRouter()
// 	routes.EpisodeRoutes(router)
// 	router.Run(configs.ENV_LAUNCH_PUBLIC_URI())
// }
