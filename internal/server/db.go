package server

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoDB struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func NewMongoDB(url string, dbname string, collectionName string) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbname).Collection(collectionName)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	m := &MongoDB{
		collection: collection,
		client:     client,
	}

	go m.Ping()
	log.Println("Connected to MongoDB!")

	return m, nil
}

func (m *MongoDB) Close() {
	err := m.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func (m *MongoDB) Ping() {
	for {
		err := m.client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Minute)
	}
}

func (m *MongoDB) Create(name string) (*mongo.InsertOneResult, error) {
	e := Entity{Name: name}
	return m.collection.InsertOne(context.Background(), e)
}

func (m *MongoDB) UpdateOrCreate(oldName string, newName string) (*Entity, error) {
	filter := bson.D{
		{"name", oldName},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	res := m.collection.FindOneAndUpdate(context.Background(), filter, bson.D{{"$set", bson.D{{"name", newName}}}}, opts)

	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, res.Err()
	}

	return &Entity{Name: newName}, nil
}

func (m *MongoDB) Delete(name string) error {
	filter := bson.D{{"name", name}}
	_, err := m.collection.DeleteOne(context.Background(), filter)
	return err
}

func (m *MongoDB) GetAll() ([]Entity, error) {
	cur, err := m.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	entities := make([]Entity, 0)
	for cur.Next(context.Background()) {
		var entity Entity
		if err := cur.Decode(&entity); err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}
