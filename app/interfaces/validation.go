package interfaces

type Body_VerifyOtp struct {
	Email string `json:"email" validate:"required,email"`
	Otp   int    `json:"otp" validate:"required,min=4,number"`
}

type Body_SendOtp struct {
	Email string `json:"email" validate:"required,email"`
}
