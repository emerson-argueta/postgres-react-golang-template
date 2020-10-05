package main

import (
	"os"
	"os/signal"

	"emersonargueta/m/v1/authorization/jwt"
	"emersonargueta/m/v1/data/postgres"
	"emersonargueta/m/v1/delivery/http"
	"emersonargueta/m/v1/paymentgateway/stripe"
)

func main() {
	databaseClient := postgres.NewClient()
	if err := databaseClient.Open(); err != nil {
		panic(err)
	}

	paymentgatewayClient := stripe.NewClient()
	setUpPaymentgateway(paymentgatewayClient, databaseClient)
	paymentgatewayClient.Initialize()

	authorizationClient := jwt.NewClient()
	setUpAuthorization(authorizationClient, databaseClient)
	authorizationClient.Initialize()

	httpServer := http.NewServer()
	httpServer.Handler = &http.Handler{
		AdministratorHandler: http.NewAdministratorHandler(),
	}
	setUpHTTPServer(httpServer, databaseClient, authorizationClient, paymentgatewayClient)

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

func setUpPaymentgateway(paymentgatewayClient *stripe.Client, databaseClient *postgres.Client) {
	paymentgatewayClient.Services.Administrator.Usecase.Services.Administrator = databaseClient.AdministratorService()

	paymentgatewayClient.Services.Administrator.Usecase.Services.User = databaseClient.UserService()
}

func setUpHTTPServer(httpServer *http.Server, databaseClient *postgres.Client, authorizationClient *jwt.Client, paymentgatewayClient *stripe.Client) {
	httpServer.Handler.AdministratorHandler.Authorization.Administrator = authorizationClient.AdministratorService()
	httpServer.Handler.AdministratorHandler.PaymentGateway.Administrator = paymentgatewayClient.AdministratorService()

	httpServer.Handler.AdministratorHandler.Usecase.Services.Administrator = databaseClient.AdministratorService()

	httpServer.Handler.AdministratorHandler.Usecase.Services.User = databaseClient.UserService()

	httpServer.Handler.AdministratorHandler.Usecase.Services.Church = databaseClient.ChurchService()

	httpServer.Handler.AdministratorHandler.Usecase.Services.Donator = databaseClient.DonatorService()

	httpServer.Handler.AdministratorHandler.Usecase.Services.Transaction = databaseClient.TransactionService()
}

func setUpAuthorization(authorizationClient *jwt.Client, databaseClient *postgres.Client) {
	authorizationClient.Services.Administrator.Usecase.Services.Administrator = databaseClient.AdministratorService()
}
