package session

import (
	"fmt"
	"github.com/gocql/gocql"
)

type SessionFactory struct {
	session *gocql.Session
}

// NewSessionFactory creates a session factory
func NewSessionFactory(driverName string) (*SessionFactory, error) {
	/*logger := logging.GetLogger()
	dataSource := config.LoadConfig(logger)
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}*/
	factory := new(SessionFactory)
	return factory, nil
}

// GetSession get a session
func (sf *SessionFactory) GetSession() *gocql.Session {
	if sf.session == nil {

		cluster := gocql.NewCluster("127.0.0.1")
		cluster.ProtoVersion = 4
		cluster.Keyspace = "taskdb"
		cluster.CQLVersion = "3.4.4"
		cluster.Consistency = gocql.Quorum
		var err error
		sf.session, err = cluster.CreateSession()
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	return sf.session
}
