package structs

import (
	"github.com/bearts/go-fiber/app/models"
)

// get allOrders
// in this, take base struct as models.Order
// replace the Location in that to models.Location

type GetAllUnassignedOrders struct {
	models.Order
	Location models.Place `json:"location"`
}
