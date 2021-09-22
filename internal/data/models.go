package data

type (
	// User is the data type for user object
	User struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Password  string `json:"password" validate:"required"`
		Username  string `json:"username" validate:"required"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		TokenHash string `json:"token_hash"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		DeletedAt string `json:"deleted_at"`
	}

	// Balance is the data type for balance object
	Balance struct {
		ID           int64   `json:"id"`
		UserID       string  `json:"userId"`
		IntegerPart  float64 `json:"integerPart"`
		FractionPart float64 `json:"fractionPart"`
		Currency     string  `json:"currency"`
	}

	PaymentHistory struct {
		ID                int64   `json:"id"`
		BalanceID         int64   `json:"balance_d"`
		CreatedAt         string  `json:"created_at"`
		InitialBalance    float64 `json:"initial_balance"`
		FinalBalance      float64 `json:"final_balance"`
		DifferenceBalance float64 `json:"difference_balance"`
	}

	AuthUser struct {
		User     User   `json:"user"`
		AuthUUID string `json:"auth_uuid"`
	}

	AuthDetails struct {
		AuthUuid string
		UserId   string
	}

	Auth struct {
		UserID   string `json:"user_id"`
		AuthUUID string `json:"auth_uuid"`
	}
)
