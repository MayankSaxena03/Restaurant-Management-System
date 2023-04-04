package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TableID   primitive.ObjectID `json:"tableID,omitempty" bson:"tableID,omitempty" validate:"required"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
}
