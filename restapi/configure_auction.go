// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/harrytucker/auction-api/database"
	"github.com/harrytucker/auction-api/restapi/operations"
	"github.com/harrytucker/auction-api/restapi/operations/bidding"

	"github.com/harrytucker/auction-api/handlers"
)

//go:generate swagger generate server --target ../../auction-api --name Auction --spec ../swagger.yml

func configureFlags(api *operations.AuctionAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AuctionAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	db, err := database.GetDB()
	if err != nil {
		log.Error(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, handlers.CtxKey("database"), db)

	api.BiddingMakeBidHandler = handlers.CreateMakeBidHandler(ctx)
	api.BiddingGetAllBidsHandler = handlers.CreateGetBidsHandler(ctx)

	if api.BiddingGetAllBidsHandler == nil {
		api.BiddingGetAllBidsHandler = bidding.GetAllBidsHandlerFunc(func(params bidding.GetAllBidsParams) middleware.Responder {
			return middleware.NotImplemented("operation bidding.GetAllBids has not yet been implemented")
		})
	}
	if api.BiddingGetAllBidsForItemHandler == nil {
		api.BiddingGetAllBidsForItemHandler = bidding.GetAllBidsForItemHandlerFunc(func(params bidding.GetAllBidsForItemParams) middleware.Responder {
			return middleware.NotImplemented("operation bidding.GetAllBidsForItem has not yet been implemented")
		})
	}
	if api.BiddingMakeBidHandler == nil {
		api.BiddingMakeBidHandler = bidding.MakeBidHandlerFunc(func(params bidding.MakeBidParams) middleware.Responder {
			return middleware.NotImplemented("operation bidding.MakeBid has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
