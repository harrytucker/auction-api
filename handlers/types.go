package handlers

import (
	"github.com/harrytucker/auction-api/models"
	"github.com/jinzhu/gorm"
)

// CtxKey is the key type for contexts in the handlers package
type CtxKey string

// db representation of a bid
type bid struct {
	gorm.Model
	models.Bid
}
