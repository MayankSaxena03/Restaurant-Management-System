package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderID        primitive.ObjectID `json:"orderID,omitempty" bson:"orderID,omitempty" validate:"required"`
	PaymentMethod  string             `json:"paymentMethod,omitempty" bson:"paymentMethod,omitempty" validate:"required"`
	PaymentStatus  string             `json:"paymentStatus,omitempty" bson:"paymentStatus,omitempty" validate:"required"`
	PaymentDueDate time.Time          `json:"paymentDueDate,omitempty" bson:"paymentDueDate,omitempty"`
	CreatedOn      time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn      time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
}
