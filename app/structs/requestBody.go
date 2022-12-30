package structs

type UserVerifyOtp struct {
	Email string `json:"email" validate:"required,email"`
	Otp   int    `json:"otp" validate:"required,min=4,number"`
}

type UserSendOtp struct {
	Email string `json:"email" validate:"required,email"`
}

type UserCreateOrder struct {
	NameOfApp        string `json:"nameOfApp" validate:"required"`
	NameOfRestaurant string `json:"nameOfRestaurant" validate:"required"`
	EstimatedTime    int    `json:"estimated_time" validate:"required"`
	DeliveryPhone    int    `json:"delivery_Phone" validate:"required"`
	Location         string `json:"location" validate:"required"`
	Otp              int    `json:"otp" validate:"required,number"`
}

type UserUpdateUser struct {
	RegistrationNumber string `json:"registrationNumber"`
	Phone              string `json:"phone"`
	Name               string `json:"name"`
	DefaultAddress     string `json:"defaultAddress"`
}

type RunnerSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RunnerSignUp struct {
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required,min=10"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RunnerDeliverOrder struct {
	Otp int `json:"otp" validate:"required"`
}

type RunnerChangeStatus struct {
	Status string `json:"status" validate:"required,oneof='waiting for delivery' 'pickedup' 'doorstep'"`
}
