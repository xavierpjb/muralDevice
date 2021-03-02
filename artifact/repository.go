package artifact

import (
	"context"
	"fmt"
	"log"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// IRepositoryHandler defines the methods needed for support artifact handling
type IRepositoryHandler interface {
	Create(RepositoryModel)
	RetrieveList(int64, int64) []RepositoryModel
}

// RepositoryModel represents the db and json models to be retrieved and sent to clients
type RepositoryModel struct {
	// ID             primitive.ObjectID `bson: "_id,omitempty"`
	URL            string    `json:"url" bson:"url,omitempty"`
	FileType       string    `json:"fileType" bson:"fileType,omitempty"`
	UploadDateTime time.Time `json:"uploadDateTime" bson:"uploadDateTime,omitempty"`
	Username       string    `json:"username" bson:"username,omitempty"`
}

// RepositoryHandler takes the collection of artifacts to perform CRUD ops
type RepositoryHandler struct {
	collection *mongo.Collection
}

// NewRH is the instantiation of an artifact repository handler
func NewRH(client *mongo.Client) RepositoryHandler {

	col := client.Database("mvral").Collection("artifact")
	a := RepositoryHandler{col}
	return a
}

// Create makes a new entry in the collection of artifacts
func (a RepositoryHandler) Create(artifactPersisted RepositoryModel) {
	// this will be correctly filled in once feature for creating entry done
	_, err := a.collection.InsertOne(context.TODO(), artifactPersisted)
	if err != nil {
		log.Fatalln(err)
	}
}

// RetrieveList get the artifacts and metadata from the db
func (a RepositoryHandler) RetrieveList(page int64, perPage int64) []RepositoryModel {
	filter := bson.M{}
	paginatedData, err := mongopagination.New(a.collection).Limit(perPage).Page(page).Sort("uploadDateTime", -1).Filter(filter).Find()
	if err != nil {
		log.Fatal(err)
	}

	var entries []RepositoryModel
	for _, raw := range paginatedData.Data {
		var art *RepositoryModel
		if marshallErr := bson.Unmarshal(raw, &art); marshallErr == nil {
			entries = append(entries, *art)
		}
	}

	fmt.Println("entries found")
	fmt.Println(entries)
	return entries
}

// Dbdriver established the connection to our mongo db
func Dbdriver() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client
}
