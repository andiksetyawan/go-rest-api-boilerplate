package reqres

type CreateUserReq struct {
	FirstName string `json:"first_name,omitempty" validate:"required"`
	LastName  string `json:"last_name,omitempty" `
	Email     string `json:"email,omitempty" validate:"email"`
}
