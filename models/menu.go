package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Category  string             `json:"category,omitempty" bson:"category,omitempty" validate:"required"`
	StartDate time.Time          `json:"startDate,omitempty" bson:"startDate,omitempty"`
	EndDate   time.Time          `json:"endDate,omitempty" bson:"endDate,omitempty"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
}
