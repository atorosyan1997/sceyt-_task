package repository

import (
	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
	"github.com/uniplaces/carbon"
	"sceyt_task/internal/data"
	"sceyt_task/pkg/logging"
)

// UserRepository is an interface for the storage implementation of the authRepository service
type UserRepository interface {
	Create(user *data.User) error
	Update(user *data.User) error
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

	sqlStr := `INSERT INTO users (id, username, firstname, lastname, createdat, updatedat, status) VALUES ( ?, ?, ?, ?, ?, ?, ?)`

	if err := r.session.Query(sqlStr, user.ID, user.Username, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt, 1).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Update(user *data.User) error {
	user.UpdatedAt = carbon.Now().String()
	sqlStr := `UPDATE users SET firstname = ?, lastname = ?, updatedat = ? WHERE username = ?`

	if err := r.session.Query(sqlStr, user.FirstName, user.LastName, user.UpdatedAt, user.Username).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(userName string) error {
	deletedAt := carbon.Now().String()
	sqlStr := `UPDATE users SET deletedat = ?, status = ? WHERE username = ? `

	if err := r.session.Query(sqlStr, deletedAt, 2, userName).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByUserName(userName string) (*data.User, error) {
	r.logger.Info("user delivered from database")
	sqlStr := `SELECT id,username, firstname, lastname FROM users WHERE username = ? and status = 1`
	user := &data.User{}
	if err := r.session.Query(sqlStr,
		userName).Consistency(gocql.One).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName); err != nil {
		return nil, err
	}

	return user, nil
}
