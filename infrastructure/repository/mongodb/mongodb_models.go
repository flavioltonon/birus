package mongodb

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"

	"github.com/flavioltonon/birus/domain/entity"
	"github.com/flavioltonon/birus/internal/shingling/classifier"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const _modelsCollection = "models"

// modelsRepository is a repository for Models
type modelsRepository repo

func (r *modelsRepository) getCollection() *mongo.Collection {
	return r.database.Collection(_modelsCollection)
}

type Model struct {
	ID         string                 `bson:"_id"`
	Classifier *classifier.Classifier `bson:"classifier"`
}

func (m Model) MarshalBSON() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(m.Classifier); err != nil {
		return nil, err
	}

	return bson.Marshal(map[string]interface{}{
		"_id":        m.ID,
		"classifier": base64.StdEncoding.EncodeToString(buffer.Bytes()),
	})
}

func (m *Model) UnmarshalBSON(b []byte) error {
	wrapper := struct {
		ID         string `bson:"_id"`
		Classifier string `bson:"classifier"`
	}{}

	if err := bson.Unmarshal(b, &wrapper); err != nil {
		return err
	}

	c, err := base64.StdEncoding.DecodeString(wrapper.Classifier)
	if err != nil {
		return err
	}

	m.ID = wrapper.ID

	return gob.NewDecoder(bytes.NewBuffer(c)).Decode(&m.Classifier)
}

// Get returns a that matches a given ID. If no Models are found, an entity.ErrNotFound will be returned.
func (r *modelsRepository) Get(ctx context.Context, modelID string) (*entity.Model, error) {
	var model Model

	err := r.getCollection().FindOne(ctx, modelID).Decode(&model)
	if err == mongo.ErrNoDocuments {
		return nil, entity.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return entity.NewModel(model.Classifier)
}

// Create creates a Model and returns its ID
func (r *modelsRepository) Create(ctx context.Context, e *entity.Model) (string, error) {
	model := Model{
		ID:         uuid.NewString(),
		Classifier: e.Classifier,
	}

	result, err := r.getCollection().InsertOne(ctx, model)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(string), nil
}

// List returns a set of Models
func (r *modelsRepository) List(ctx context.Context) ([]*entity.Model, error) {
	var (
		pipeline mongo.Pipeline
		models   []Model
	)

	cursor, err := r.getCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	es := make([]*entity.Model, 0, len(models))

	for _, model := range models {
		e, err := entity.NewModel(model.Classifier)
		if err != nil {
			return nil, entity.ErrInvalidEntity
		}

		es = append(es, e)
	}

	return es, nil
}
