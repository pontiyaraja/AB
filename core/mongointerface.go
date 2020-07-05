package core

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Create - inserts data into mongo database
func Create(ctx APIContext, db, collectionName string, d interface{}) (interface{}, error) {
	mongo, err := GetS()
	if err != nil {
		return nil, err
	}
	sessOpt := &options.SessionOptions{
		DefaultMaxCommitTime: func(timeout time.Duration) *time.Duration { return &timeout }(time.Duration(5 * time.Second)),
		DefaultReadConcern:   readconcern.Majority(),
		//DefaultWriteConcern:  wri
	}
	session, err := mongo.StartSession(sessOpt)
	if err != nil {
		return nil, err
	}
	ctx1 := context.Background()
	defer session.EndSession(ctx1)

	txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
	return session.WithTransaction(ctx1, func(sessCtx mgo.SessionContext) (interface{}, error) {
		return mongo.Database(db).Collection(collectionName).InsertOne(sessCtx, d)
	}, txnOpts)
}

// ReadOne - inserts data into mongo database
func ReadOne(db, collectionName string, selector, filter bson.M, data interface{}) error {
	mongo, err := GetS()
	if err != nil {
		return err
	}
	sessOpt := &options.SessionOptions{
		DefaultMaxCommitTime: func(timeout time.Duration) *time.Duration { return &timeout }(time.Duration(5 * time.Second)),
		DefaultReadConcern:   readconcern.Majority(),
		//DefaultWriteConcern:  wri
	}
	session, err := mongo.StartSession(sessOpt)
	if err != nil {
		return err
	}
	ctx1 := context.Background()
	defer session.EndSession(ctx1)

	txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
	_, err = session.WithTransaction(ctx1, func(sessCtx mgo.SessionContext) (interface{}, error) {
		return nil, mongo.Database(db).Collection(collectionName).FindOne(ctx1, filter).Decode(data)
	}, txnOpts)
	return err
}

// ReadAll - inserts data into mongo database
func ReadAll(db, collectionName string, selector, filter bson.M, data interface{}) error {
	mongo, err := GetS()
	if err != nil {
		return err
	}
	sessOpt := &options.SessionOptions{
		DefaultMaxCommitTime: func(timeout time.Duration) *time.Duration { return &timeout }(time.Duration(5 * time.Second)),
		DefaultReadConcern:   readconcern.Majority(),
		//DefaultWriteConcern:  wri
	}
	session, err := mongo.StartSession(sessOpt)
	if err != nil {
		return err
	}
	ctx1 := context.Background()
	defer session.EndSession(ctx1)

	txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
	_, err = session.WithTransaction(ctx1, func(sessCtx mgo.SessionContext) (interface{}, error) {
		cur, err := mongo.Database(db).Collection(collectionName).Find(ctx1, filter)
		if err != nil {
			return nil, err
		}
		err = cur.All(ctx1, data)
		return nil, err
	}, txnOpts)
	return err
}

// Count returns number of documents those satisfied findQuery
func Count(db, collectionName string, findQ bson.M) (int64, error) {
	mongo, err := GetS()
	if err != nil {
		err = fmt.Errorf("failed to get mongo session error: %v", err)
		return 0, err
	}
	sessOpt := &options.SessionOptions{
		DefaultMaxCommitTime: func(timeout time.Duration) *time.Duration { return &timeout }(time.Duration(5 * time.Second)),
		DefaultReadConcern:   readconcern.Majority(),
		//DefaultWriteConcern:  wri
	}
	session, err := mongo.StartSession(sessOpt)
	if err != nil {
		return 0, err
	}
	ctx1 := context.Background()
	defer session.EndSession(ctx1)
	txnOpts := options.Transaction().SetReadPreference(readpref.PrimaryPreferred())
	count, err := session.WithTransaction(ctx1, func(sessCtx mgo.SessionContext) (interface{}, error) {
		return mongo.Database(db).Collection(collectionName).CountDocuments(ctx1, findQ)
	}, txnOpts)
	return count.(int64), err
}
