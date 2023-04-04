package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text      string             `json:"text,omitempty" bson:"text,omitempty"`
	OrderID   primitive.ObjectID `json:"orderID,omitempty" bson:"orderID,omitempty"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
}
