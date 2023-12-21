package models

import (
	cm "geotask/module/courier/models"
	om "geotask/module/order/models"
)

type CourierStatus struct {
	Courier cm.Courier `json:"courier"`
	Orders  []om.Order `json:"orders"`
}
