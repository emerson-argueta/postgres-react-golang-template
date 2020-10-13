package main

import (
	"os"
	"os/signal"

	"emersonargueta/m/v1/authorization"
	"emersonargueta/m/v1/communitygoaltracker"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/data/postgres"
	"emersonargueta/m/v1/delivery/http"
	"emersonargueta/m/v1/delivery/middleware"
	"emersonargueta/m/v1/identity"
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

	httpServer := setUpHTTPServer(databaseClient, Config)

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

func setUpHTTPServer(databaseClient *postgres.Client, config *config.Config) (res *http.Server) {
	res = http.NewServer(Config)

	cgtHandler := http.NewCommunitygoaltrackerHandler()
	cgtHandler.Middleware = middleware.NewClient(config).JwtService()
	cgtHandler.Authorization = authorization.NewClient(config).JwtService()

	cgtHandler.Communitygoaltracker.Processes = communitygoaltracker.NewClient().Service()
	cgtHandler.Communitygoaltracker.Achiever = databaseClient.AchieverService()
	cgtHandler.Communitygoaltracker.Goal = databaseClient.GoalService()
	cgtHandler.Communitygoaltracker.Identity.Processes = identity.NewClient().Service()
	cgtHandler.Communitygoaltracker.Identity.User = databaseClient.UserService()
	cgtHandler.Communitygoaltracker.Identity.Domain = databaseClient.DomainService()

	cgtHandler.Initialize()

	res.Handler = &http.Handler{

		CommunitygoaltrackerHandler: cgtHandler,
	}

	return res

}
