package model

import uuid "github.com/satori/go.uuid"

type FlightRequest struct {
	DestinationFrom string `json:"destination_from" binding:"required"`
	DestinationTo   string `json:"destination_to" binding:"required"`
}

type FlightDTO struct {
	FlightID        uuid.UUID `db:"flight_id"`
	DestinationFrom string    `db:"destination_from"`
	DestinationTo   string    `db:"destination_to"`
}

type FlightResponse struct {
	ID              uuid.UUID `json:"id"`
	DestinationFrom string    `json:"destination_from"`
	DestinationTo   string    `json:"destination_to"`
}

func (dto *FlightDTO) ToFlightResponse() FlightResponse {
	if dto == nil {
		return FlightResponse{}
	}
	return FlightResponse{
		ID:              dto.FlightID,
		DestinationFrom: dto.DestinationFrom,
		DestinationTo:   dto.DestinationTo,
	}
}

func ToFlightDTO(req FlightRequest) FlightDTO {
	return FlightDTO{
		DestinationFrom: req.DestinationFrom,
		DestinationTo:   req.DestinationTo,
	}
}

func ToFlightDTOWithID(req FlightRequest, id uuid.UUID) FlightDTO {
	return FlightDTO{
		FlightID:        id,
		DestinationFrom: req.DestinationFrom,
		DestinationTo:   req.DestinationTo,
	}
}
