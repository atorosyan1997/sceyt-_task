package reposytory

import (
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
	"github.com/uniplaces/carbon"
	"my-bank-service/internal/data"
	"my-bank-service/pkg/logging"
)

// UserRepository is an interface for the storage implementation of the authRepository service
type UserRepository interface {
	Create(user *data.User) error
	Update(user *data.User, userName string) error
	Delete(userName string) error
	GetUserByUserName(userName string) (*data.User, error)
}

// userRepository has the implementation of the db methods.
type userRepository struct {
	session *gocql.Session
	logger  logging.Logger
}

// NewUserRepository returns a new userRepository instance
func NewUserRepository(s *gocql.Session, l logging.Logger) UserRepository {
	return &userRepository{s, l}
}

func (r *userRepository) Create(user *data.User) error {
	user.ID = uuid.NewV4().String()
	user.CreatedAt = carbon.Now().String()
	user.UpdatedAt = carbon.Now().String()

	sqlStr := `INSERT INTO users (id, email, userName, firstName, lastName, password, tokenhash, createdat, updatedat) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if err := r.session.Query(sqlStr, user.ID, user.Email, user.Username, user.FirstName, user.LastName, user.Password,
		user.TokenHash, user.CreatedAt, user.UpdatedAt).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(user *data.User, userName string) error {
	panic("implement me")
}

func (r *userRepository) Delete(userName string) error {
	panic("implement me")
}

func (r *userRepository) GetUserByUserName(userName string) (*data.User, error) {
	sqlStr := `SELECT id, email, userName, firstName, lastName, password, tokenhash, createdat, updatedat FROM users WHERE userName = ?`
	user := &data.User{}
	if err := r.session.Query(sqlStr,
		userName).Consistency(gocql.One).Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Password, &user.TokenHash,
		&user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return user, nil
}
