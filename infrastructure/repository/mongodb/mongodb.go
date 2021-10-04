package mongodb

import (
	"context"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository is a MongoDB repository
type Repository struct {
	client *mongo.Client

	common               repo
	ClassifierRepository *classifierRepository

	options *Options
}

type Options struct {
	DatabaseName  string
	ClientOptions []*options.ClientOptions
}

func (o Options) validate() error {
	return ozzo.ValidateStruct(&o,
		ozzo.Field(&o.DatabaseName, ozzo.Required),
	)
}

type repo struct {
	database *mongo.Database
}

// NewRepository creates a new Repository
func NewRepository(options *Options) (*Repository, error) {
	if err := options.validate(); err != nil {
		return nil, err
	}

	return &Repository{options: options}, nil
}

// Connect initializes the repository
func (r *Repository) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, r.options.ClientOptions...)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	r.client = client
	r.common.database = r.client.Database(r.options.DatabaseName)
	r.ClassifierRepository = (*classifierRepository)(&r.common)
	return nil
}

// Disconnect interrupts the connection of the repository with MongoDB servers
func (r *Repository) Disconnect(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
