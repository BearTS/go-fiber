package interfaces

type Body_VerifyOtp struct {
	Email string `json:"email" validate:"required,email"`
	Otp   int    `json:"otp" validate:"required,min=4,number"`
}

type Body_SendOtp struct {
	Email string `json:"email" validate:"required,email"`
}

type Body_CreateOrder struct {
	NameOfApp        string `json:"nameOfApp" validate:"required"`
	NameOfRestaurant string `json:"nameOfRestaurant" validate:"required"`
	EstimatedTime    int    `json:"estimated_time" validate:"required"`
	DeliveryPhone    int    `json:"delivery_Phone" validate:"required"`
	Location         string `json:"location" validate:"required"`
	Otp              int    `json:"otp" validate:"required,number"`
}
