package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBRepository is a MongoDB repository
type MongoDBRepository struct {
	client *mongo.Client

	common repository

	databaseName string
	opts         []*options.ClientOptions
}

type repository struct {
	database *mongo.Database
}

// NewMongoDBRepository creates a new MongoDBRepository
func NewMongoDBRepository(databaseName string, opts ...*options.ClientOptions) *MongoDBRepository {
	return &MongoDBRepository{
		databaseName: databaseName,
		opts:         opts,
	}
}

// Connect initializes the repository
func (r *MongoDBRepository) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, r.opts...)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	r.client = client
	r.common.database = r.client.Database(r.databaseName)
	return nil
}

// Disconnect interrupts the connection of the repository with MongoDB servers
func (r *MongoDBRepository) Disconnect(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
