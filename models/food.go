package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=50"`
	Price     float64            `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	FoodImage string             `json:"foodImage,omitempty" bson:"foodImage,omitempty" validate:"required"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UpdatedOn time.Time          `json:"updatedOn,omitempty" bson:"updatedOn,omitempty"`
	MenuID    primitive.ObjectID `json:"menuID,omitempty" bson:"menuID,omitempty"`
}
