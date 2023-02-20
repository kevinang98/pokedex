package databases

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func NewMockMongoConn(t *testing.T) (*mongo.Client, error) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	return mt.Client, nil
}
