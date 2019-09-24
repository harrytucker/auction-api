package handlers

import (
	"context"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/harrytucker/auction-api/models"
	"github.com/harrytucker/auction-api/restapi/operations/bidding"
)

// CreateGetBidsHandler instantiates a GetAllBids struct to handle requests
func CreateGetBidsHandler(ctx context.Context) *GetAllBids {
	return &GetAllBids{ctx: ctx}
}

// GetAllBids stores context for this handler
type GetAllBids struct {
	ctx context.Context
}

// Handle contains business logic for making a bid on an item
func (gab *GetAllBids) Handle(params bidding.GetAllBidsParams) middleware.Responder {
	log.Info("Request received to get all bids")
	db := gab.ctx.Value(CtxKey("database")).(*gorm.DB)

	var allBids []*models.Bid
	db.Find(&allBids)

	var results []models.ItemSummary
	results = append(results, allBids)

	log.WithField("Results", allBids).Info("Checking first bid in DB")
	return bidding.NewGetAllBidsOK().WithPayload(results)
}
