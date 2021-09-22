package reposytory

import (
	"github.com/gocql/gocql"
	"my-bank-service/internal/data"
	"my-bank-service/pkg/logging"
)

type AuthRepository interface {
	FetchAuth(authD *data.AuthDetails) (*data.Auth, error)
	DeleteAuth(authD *data.AuthDetails) error
	CreateAuth(authD *data.AuthDetails) (*data.Auth, error)
}

// authRepository has the implementation of the db methods.
type authRepository struct {
	session *gocql.Session
	logger  logging.Logger
}

// NewAuthRepository returns a new userRepository instance
func NewAuthRepository(s *gocql.Session, l logging.Logger) AuthRepository {
	return &authRepository{s, l}
}

func (a *authRepository) FetchAuth(authD *data.AuthDetails) (*data.Auth, error) {
	sqlStr := `SELECT * FROM auth WHERE user_id = ? and auth_uuid = ?`

	au := &data.Auth{}
	if err := a.session.Query(sqlStr, authD.UserId, authD.AuthUuid).Consistency(gocql.One).Scan(&au.UserID, &au.AuthUUID); err != nil {
		return nil, err
	}

	return au, nil
}

func (a *authRepository) DeleteAuth(authD *data.AuthDetails) error {
	sqlStr := `DELETE FROM auth WHERE user_id = ? and auth_uuid = ?`
	if err := a.session.Query(sqlStr, authD.UserId, authD.AuthUuid).Exec(); err != nil {
		return err
	}
	return nil
}

func (a *authRepository) CreateAuth(authD *data.AuthDetails) (*data.Auth, error) {
	auth, err := a.FetchAuth(authD)
	if auth != nil {
		return auth, nil
	}
	au := &data.Auth{}
	au.AuthUUID = authD.AuthUuid
	au.UserID = authD.UserId
	sqlStr := `INSERT INTO auth (user_id,auth_uuid) VALUES (?, ?)`
	if err = a.session.Query(sqlStr, au.UserID, au.AuthUUID).Exec(); err != nil {
		return nil, err
	}
	return au, nil
}
