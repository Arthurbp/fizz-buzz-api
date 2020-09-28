package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// FizzBuzzQuery is struct params
type FizzBuzzQuery struct {
	Str1   string `json:"str1" schema:"str1"`
	Str2   string `json:"str2" schema:"str2"`
	Int1   int    `json:"int1" schema:"int1"`
	Int2   int    `json:"int2" schema:"int2"`
	Limit  int    `json:"limit" schema:"limit"`
	NbHits int32  `json:"nbHits" schema:"-"`
}

// Store defines mongo database
type Store struct {
	C *mongo.Client
	l *log.Logger
}

// NewClient return a new mongo client
func NewClient(host string, l *log.Logger) (*Store, error) {

	r := new(Store)
	r.l = l

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))
	if err != nil {
		return nil, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	r.C = client
	return r, nil
}

// InsertFizzBuzzQuery inserts the query in database
func (s *Store) InsertFizzBuzzQuery(ctx context.Context, q FizzBuzzQuery) error {

	if _, err := s.C.Database("app").Collection("queries").InsertOne(ctx, q); err != nil {
		return err
	}

	return nil
}

// AggregateFizzBuzzQueries returns the queries sorted by nbHits
func (s *Store) AggregateFizzBuzzQueries(ctx context.Context) ([]FizzBuzzQuery, error) {

	sortQueries := []FizzBuzzQuery{}

	group := bson.M{
		"_id": bson.M{
			"str1":  "$str1",
			"str2":  "$str2",
			"int1":  "$int1",
			"int2":  "$int2",
			"limit": "$limit",
		},
		"nbHits": bson.M{"$sum": 1},
	}

	sort := bson.M{
		"nbHits": -1,
	}

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: group}},
		{{Key: "$sort", Value: sort}},
	}

	cursor, err := s.C.Database("app").Collection("queries").Aggregate(ctx, pipeline)
	if err != nil {
		return sortQueries, err
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return sortQueries, err
	}

	for _, result := range results {
		var fbQuery FizzBuzzQuery
		bsonBytes, _ := bson.Marshal(result["_id"])
		bson.Unmarshal(bsonBytes, &fbQuery)
		fbQuery.NbHits = result["nbHits"].(int32)
		sortQueries = append(sortQueries, fbQuery)
	}

	return sortQueries, nil
}

// StartTestContainer runs test container for the store
func StartTestContainer() (*Store, error) {
	port := "27017"
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.4",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)),
	}
	mongo, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	ip, err := mongo.Host(ctx)
	if err != nil {
		return nil, err
	}
	mappedPort, err := mongo.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, err
	}
	s, err := NewClient(fmt.Sprintf("mongodb://%s:%s", ip, mappedPort.Port()), nil)
	if err != nil {
		return nil, err
	}
	return s, nil
}
