package main

import (
	"os"
	"os/signal"

	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/data/postgres"
	"emersonargueta/m/v1/delivery/http"
	"emersonargueta/m/v1/delivery/middleware"
)

// Config variables for services
var Config *config.Config

func init() {
	Config = config.NewConfig()
}

func main() {
	databaseClient := postgres.NewClient(Config)
	if err := databaseClient.Open(); err != nil {
		panic(err)
	}

	httpServer := http.NewServer(Config)
	httpServer.Handler = &http.Handler{
		CommunitygoaltrackerHandler: http.NewCommunitygoaltrackerHandler(Config),
	}
	setUpHTTPServer(httpServer, databaseClient, Config)

	if err := httpServer.Open(); err != nil {
		panic(err)
	}
	defer httpServer.Close()

	// Block until an OS interrupt is received.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sig := <-ch
	println("Got signal:", sig)

}

func setUpHTTPServer(httpServer *http.Server, databaseClient *postgres.Client, config *config.Config) {
	httpServer.Handler.Authorization = authorization.NewClient(config).JwtService()
	httpServer.Handler.Middleware = middleware.NewClient(config).JwtService()

	httpServer.Handler.Communitygoaltracker.Achiever = databaseClient.AchieverService()
	httpServer.Handler.Communitygoaltracker.Goal = databaseClient.GoalService()

	httpServer.Handler.Communitygoaltracker.Identity.User = databaseClient.UserService()
	httpServer.Handler.Communitygoaltracker.Identity.Domain = databaseClient.DomainService()

}
