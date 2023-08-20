package dto

type CreateUserRequest struct {
	Username             string `json:"username" validate:"required,alphanum"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=Password"`
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,currency"`
}

type CreateAccountRequest struct {
	Currency string `validate:"required,currency"`
}

type GetAccountRequest struct {
	ID int64 `validate:"required,min=1"`
}

type GetListAccountsRequest struct {
	Page    int32 `validate:"min=1"`
	PerPage int32 `query:"per_page" validate:"min=1"`
}
