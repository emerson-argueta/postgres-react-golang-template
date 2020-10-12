package main

import (
	"os"
	"os/signal"

	"emersonargueta/m/v1/authorization/jwt"
	"emersonargueta/m/v1/config"
	"emersonargueta/m/v1/data/postgres"
	"emersonargueta/m/v1/delivery/http"
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

	authorizationClient := jwt.NewClient(Config)

	httpServer := http.NewServer(Config)
	httpServer.Handler = &http.Handler{
		CommunitygoaltrackerHandler: http.NewCommunitygoaltrackerHandler(Config),
	}
	setUpHTTPServer(httpServer, databaseClient, authorizationClient)

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

func setUpHTTPServer(httpServer *http.Server, databaseClient *postgres.Client, authorizationClient *jwt.Client) {
	httpServer.Handler.Authorization.Jwt = authorizationClient.JwtService()

	httpServer.Handler.Communitygoaltracker.Achiever = databaseClient.AchieverService()
	httpServer.Handler.Communitygoaltracker.Goal = databaseClient.GoalService()

	identityClient := identity.NewClient(Config)
	identityClient.Services.Identity.User = databaseClient.UserService()
	identityClient.Services.Identity.Domain = databaseClient.DomainService()
	httpServer.Handler.Communitygoaltracker.Identity = &identityClient.Services.Identity

}
