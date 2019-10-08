package handlers

import (
	"context"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/harrytucker/auction-api/models"
	"github.com/harrytucker/auction-api/restapi/operations/statistics"
)

// CreateGetStatsHandler instantiates a GetStats struct to handle requests
func CreateGetStatsHandler(ctx context.Context) *GetStats {
	return &GetStats{ctx: ctx}
}

// GetStats stores context for this handler
type GetStats struct {
	ctx context.Context
}

// Handle contains business logic for getting statistics about an item
func (gs *GetStats) Handle(params statistics.GetBidStatsParams) middleware.Responder {
	log.WithFields(log.Fields{
		"Item Number": params.ItemNumber,
	}).Info("Item statistics request received")

	db := gs.ctx.Value(CtxKey("database")).(*gorm.DB)

	bids := models.ItemSummary{}
	db.Where("item_number = ?", params.ItemNumber).Order("bid_amount desc").Find(&bids)
	noOfBids := len(bids)


	return statistics.NewGetBidStatsOK().WithPayload(&models.Statistics{NoOfBids: int64(noOfBids)})
}
