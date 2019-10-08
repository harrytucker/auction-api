// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/harrytucker/auction-api/database"
	"github.com/harrytucker/auction-api/models"
	"github.com/harrytucker/auction-api/restapi/operations"
	"github.com/harrytucker/auction-api/restapi/operations/bidding"
	"github.com/harrytucker/auction-api/restapi/operations/statistics"

	"github.com/harrytucker/auction-api/handlers"
)

//go:generate swagger generate server --target ../../auction-api --name Auction --spec ../swagger.yml

func configureFlags(api *operations.AuctionAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AuctionAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// logging configuration
	api.Logger = log.Infof

	logFile, err := os.Create("logs/log_" + time.Now().Format("2006-01-02_15:04:05"))
	if err != nil {
		log.Fatal("Failed to create logfile")
	}
	mw := io.Writer(logFile)
	log.SetOutput(mw)

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	db, err := database.GetDB()
	if err != nil {
		log.Error(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, handlers.CtxKey("database"), db)

	api.BiddingMakeBidHandler = handlers.CreateMakeBidHandler(ctx)
	api.StatisticsGetBidStatsHandler = handlers.CreateGetStatsHandler(ctx)

	// DISABLED WHILE AUCTION IS ACTIVE TO PREVENT MEDDLING BY CHEEKY ENGINEERS
	//
	// api.BiddingGetAllBidsHandler = handlers.CreateGetBidsHandler(ctx)
	// api.BiddingGetAllBidsForItemHandler = handlers.CreateGetBidsForItemHandler(ctx)
	disabledError := models.ErrorResponse{ErrorMessage: "SELF DESTRUCT ACTIVATED: Please stop messing with the auction API"}

	api.BiddingGetAllBidsHandler = bidding.GetAllBidsHandlerFunc(func(params bidding.GetAllBidsParams) middleware.Responder {
		log.WithFields(log.Fields{
			"path":      params.HTTPRequest.URL.Path,
			"requester": params.HTTPRequest.RemoteAddr,
		}).Warn("Request attempting to see extra bid info")
		return bidding.NewGetAllBidsUnauthorized().WithPayload(&disabledError)
	})

	api.BiddingGetAllBidsForItemHandler = bidding.GetAllBidsForItemHandlerFunc(func(params bidding.GetAllBidsForItemParams) middleware.Responder {
		log.WithFields(log.Fields{
			"path":      params.HTTPRequest.URL.Path,
			"requester": params.HTTPRequest.RemoteAddr,
		}).Warn("Request attempting to see extra bid info")
		return bidding.NewGetAllBidsUnauthorized().WithPayload(&disabledError)
	})
	if api.BiddingMakeBidHandler == nil {
		api.BiddingMakeBidHandler = bidding.MakeBidHandlerFunc(func(params bidding.MakeBidParams) middleware.Responder {
			return middleware.NotImplemented("operation bidding.MakeBid has not yet been implemented")
		})
	}
	if api.StatisticsGetBidStatsHandler == nil {
		api.StatisticsGetBidStatsHandler = statistics.GetBidStatsHandlerFunc(func(params statistics.GetBidStatsParams) middleware.Responder {
			return middleware.NotImplemented("operation statistics.GetBidStats has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {
		db.Close()
		logFile.Close()
	}

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
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}
