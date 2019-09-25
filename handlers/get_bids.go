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

// Handle contains business logic for getting all bids in the database
func (gab *GetAllBids) Handle(params bidding.GetAllBidsParams) middleware.Responder {
	log.Info("Get bids request received")
	db := gab.ctx.Value(CtxKey("database")).(*gorm.DB)

	var allBids []*models.Bid
	db.Find(&allBids)

	var items []models.ItemSummary
	items = append(items, allBids)

	log.WithField("Results", items).Info("Successful get!")
	return bidding.NewGetAllBidsOK().WithPayload(items)
}
