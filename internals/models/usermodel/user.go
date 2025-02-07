package usermodel






/*Declare a struct that holds the user  login credentials*/

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

type SignupRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ConfirmPass string `json:"confirm_password"`
}

type SignupResponse struct {
	Message string `json:"message"`
}