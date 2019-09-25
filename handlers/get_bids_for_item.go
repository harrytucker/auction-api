package handlers

import (
	"context"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/harrytucker/auction-api/models"
	"github.com/harrytucker/auction-api/restapi/operations/bidding"
)

// CreateGetBidsForItemHandler instantiates a GetAllBidsForItem struct to handle requests
func CreateGetBidsForItemHandler(ctx context.Context) *GetAllBidsForItem {
	return &GetAllBidsForItem{ctx: ctx}
}

// GetAllBidsForItem stores context for this handler
type GetAllBidsForItem struct {
	ctx context.Context
}

// Handle contains business logic for making a bid on an item
func (gab *GetAllBidsForItem) Handle(params bidding.GetAllBidsForItemParams) middleware.Responder {
	log.Info("Request received to get all bids")
	db := gab.ctx.Value(CtxKey("database")).(*gorm.DB)

	bids := models.ItemSummary{}
	db.Where("item_number = ?", params.ItemNumber).Find(&bids)

	return bidding.NewGetAllBidsForItemOK().WithPayload(bids)
}
