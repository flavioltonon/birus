package mongodb

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"

	"birus/domain/entity/shingling/classifier"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const _classifiersCollection = "classifiers"

// classifierRepository is a repository for Models
type classifierRepository repo

func (r *classifierRepository) getCollection() *mongo.Collection {
	return r.database.Collection(_classifiersCollection)
}

type Classifier struct {
	data *classifier.Classifier
}

func (c Classifier) MarshalBSON() ([]byte, error) {
	var buffer bytes.Buffer

	if err := gob.NewEncoder(&buffer).Encode(c.data); err != nil {
		return nil, err
	}

	return bson.Marshal(map[string]interface{}{
		"_id":        c.data.ID(),
		"name":       c.data.Name(),
		"classifier": base64.StdEncoding.EncodeToString(buffer.Bytes()),
	})
}

func (c *Classifier) UnmarshalBSON(b []byte) error {
	wrapper := struct {
		ID         string `bson:"_id"`
		Name       string `bson:"name"`
		Classifier string `bson:"classifier"`
	}{}

	if err := bson.Unmarshal(b, &wrapper); err != nil {
		return err
	}

	b, err := base64.StdEncoding.DecodeString(wrapper.Classifier)
	if err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(&c.data)
}

// Create creates a Classifier
func (r *classifierRepository) CreateClassifier(ctx context.Context, classifier *classifier.Classifier) error {
	if _, err := r.getCollection().InsertOne(ctx, Classifier{data: classifier}); err != nil {
		return err
	}

	return nil
}

// List returns a set of Classifiers
func (r *classifierRepository) ListClassifiers(ctx context.Context) ([]*classifier.Classifier, error) {
	var (
		pipeline    mongo.Pipeline
		classifiers []Classifier
	)

	cursor, err := r.getCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &classifiers); err != nil {
		return nil, err
	}

	es := make([]*classifier.Classifier, 0, len(classifiers))

	for _, classifier := range classifiers {
		es = append(es, classifier.data)
	}

	return es, nil
}
