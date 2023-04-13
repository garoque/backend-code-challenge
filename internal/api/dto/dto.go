package dto

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Err  error       `json:"error,omitempty"`
}

type CreateUser struct {
	Name string `json:"name" validate:"required"`
}

type CreateTransaction struct {
	SourceUserId      string  `json:"sourceUserId" validate:"required"`
	DestinationUserId string  `json:"destinationUserId" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
}

type IncreaseBalanceUser struct {
	UserId string  `json:"userId" validate:"required"`
	Value  float64 `json:"value" validate:"required"`
}
