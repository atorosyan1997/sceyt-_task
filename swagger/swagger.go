package swagger

type (
	UserSearchDelete struct {
		Username string `json:"username" validate:"required"`
	}

	UserAddUpdate struct {
		Username  string `json:"username" validate:"required"`
		FirstName string `json:"firstname" validate:"required"`
		LastName  string `json:"lastname" validate:"required"`
	}
)
