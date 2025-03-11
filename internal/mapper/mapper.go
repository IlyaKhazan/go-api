package mapper

import "go-api/internal/model"

func ToFlightDTO(req model.FlightRequest) model.FlightDTO {
	return model.FlightDTO{
		DestinationFrom: req.DestinationFrom,
		DestinationTo:   req.DestinationTo,
	}
}

func ToFlightDTOWithID(req model.FlightRequest, id int) model.FlightDTO {
	return model.FlightDTO{
		FlightID:        id,
		DestinationFrom: req.DestinationFrom,
		DestinationTo:   req.DestinationTo,
	}
}

func ToFlightResponse(dto model.FlightDTO) model.FlightResponse {
	return model.FlightResponse{
		ID:              dto.FlightID,
		DestinationFrom: dto.DestinationFrom,
		DestinationTo:   dto.DestinationTo,
	}
}
