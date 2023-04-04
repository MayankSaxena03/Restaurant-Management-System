package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Number    int                `json:"number,omitempty" bson:"number,omitempty" validate:"required"`
	Guests    int                `json:"guests,omitempty" bson:"guests,omitempty" validate:"required"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
}
