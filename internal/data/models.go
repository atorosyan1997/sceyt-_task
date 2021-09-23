package data

type (
	// User is the data type for user object
	User struct {
		ID        string `json:"id"`
		Username  string `json:"username" validate:"required"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		DeletedAt string `json:"deleted_at"`
	}
)
