package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var (
	connection *DBConnection
	cOnce      sync.Once
)

type DBConnection struct {
	m *mongo.Client
	l *sync.RWMutex
}

func (c *DBConnection) set(value *mongo.Client) {
	c.l.Lock()
	defer c.l.Unlock()
	c.m = value
}

func (c *DBConnection) get() (*mongo.Client, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	if c.m == nil {
		return c.m, fmt.Errorf("empty session")
	}
	return c.m, nil
}

func (c *DBConnection) connect() error {
	c.l.Lock()
	defer c.l.Unlock()

	clientOptions := &options.ClientOptions{
		// Addrs:     []string{"localhost:27017"},
		// Username:  "",
		// Password:  "",
		ConnectTimeout: func(timeout time.Duration) *time.Duration { return &timeout }(time.Duration(15 * time.Second)),
		MaxPoolSize:    func(maxPool uint64) *uint64 { return &maxPool }(256),
		Direct:         func(isDirect bool) *bool { return &isDirect }(true),
	}

	clientOptions = clientOptions.ApplyURI("mongodb://localhost:27017/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err == nil {
		// set the session object as the value in the map for the respective tenantName as key
		c.m = client
	} else {
		return errors.Wrap(err, "error dialing mongo")
	}
	return nil
}

// Singleton pattern
func newConection() *DBConnection {
	cOnce.Do(func() {
		connection = &DBConnection{
			l: new(sync.RWMutex),
		}
	})
	return connection
}

func init() {
	connection = newConection()
}

func connectMongo() error {
	err := connection.connect()
	if err != nil {
		fmt.Println("error connecting mongo")
	}
	return err
}

// GetS returns session for database, if session is already created for database it returns session copy.
func GetS() (*mongo.Client, error) {
	session, err := connection.get()
	if err != nil {
		if err := connectMongo(); err == nil {
			return GetS() //singleton recursion to again call GetS
		}
		return nil, err
	}
	return session, err
}
