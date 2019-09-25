package handlers

import (
	"context"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/harrytucker/auction-api/models"
	"github.com/harrytucker/auction-api/restapi/operations/bidding"
)

// CreateMakeBidHandler instantiates a MakeBid struct to handle requests
func CreateMakeBidHandler(ctx context.Context) *MakeBid {
	return &MakeBid{ctx: ctx}
}

// MakeBid stores context for this handler
type MakeBid struct {
	ctx context.Context
}

// Handle contains business logic for making a bid on an item
func (mb *MakeBid) Handle(params bidding.MakeBidParams) middleware.Responder {
	log.WithFields(log.Fields{
		"Item Number": params.ItemNumber,
		"Bid Amount":  *params.Bid.BidAmount,
		"Bidder":      *params.Bid.BidderName,
	}).Info("Bid request received")

	db := mb.ctx.Value(CtxKey("database")).(*gorm.DB)
	db.AutoMigrate(&bid{})

	bid := bid{Bid: models.Bid{
		ItemNumber:  &params.ItemNumber,
		BidAmount:   params.Bid.BidAmount,
		BidderName:  params.Bid.BidderName,
		BidderEmail: params.Bid.BidderEmail,
	}}
	db.Create(&bid)

	log.WithField("Results", bid).Info("Successful bid!")
	return bidding.NewMakeBidOK().WithPayload(params.Bid)
}
