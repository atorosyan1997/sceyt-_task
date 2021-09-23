package session

import (
	"github.com/gocql/gocql"
	"log"
	"sceyt_task/internal/config"
)

type SessionFactory struct {
	session *gocql.Session
}

// NewSessionFactory creates a session factory
func NewSessionFactory() (*SessionFactory, error) {
	dbConfig := config.LoadConfig()
	cluster := gocql.NewCluster(dbConfig.Address)
	cluster.ProtoVersion = dbConfig.ProtoVersion
	cluster.Keyspace = dbConfig.Keyspace
	cluster.CQLVersion = dbConfig.CQLVersion
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: dbConfig.Username,
		Password: dbConfig.Password,
	}
	var err error
	session, err := cluster.CreateSession()
	if err != nil {
		log.Panic(err)
	}
	factory := &SessionFactory{
		session: session,
	}
	return factory, nil
}

// GetSession get a session
func (sf *SessionFactory) GetSession() *gocql.Session {
	return sf.session
}
