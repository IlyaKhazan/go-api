package model

type FlightRequest struct {
	DestinationFrom string `json:"destination_from" binding:"required"`
	DestinationTo   string `json:"destination_to" binding:"required"`
}

type FlightDTO struct {
	FlightID        int    `db:"flight_id"`
	DestinationFrom string `db:"destination_from"`
	DestinationTo   string `db:"destination_to"`
}

type FlightResponse struct {
	ID              int    `json:"id"`
	DestinationFrom string `json:"destination_from"`
	DestinationTo   string `json:"destination_to"`
}
