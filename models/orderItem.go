package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Quantity  int                `json:"quantity,omitempty" bson:"quantity,omitempty" validate:"required"`
	UnitPrice float64            `json:"unitPrice,omitempty" bson:"unitPrice,omitempty" validate:"required"`
	FoodID    primitive.ObjectID `json:"foodID,omitempty" bson:"foodID,omitempty" validate:"required"`
	OrderID   primitive.ObjectID `json:"orderID,omitempty" bson:"orderID,omitempty" validate:"required"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
}
