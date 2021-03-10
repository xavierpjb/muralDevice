package artifact

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "muraldevice/migrations"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// IRepositoryHandler defines the methods needed for support artifact handling
type IRepositoryHandler interface {
	Create(RepositoryModel)
	RetrieveList(int64, int64) []RepositoryModel
	Delete(DeleteModel)
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
	var entries []RepositoryModel
	_, err := mongopagination.New(a.collection).Limit(perPage).Page(page).Sort("uploadDateTime", -1).Filter(filter).Decode(&entries).Find()
	if err != nil {
		log.Fatal(err)
	}

	return entries
}

// Delete removes the artifact from mongo db
func (a RepositoryHandler) Delete(artifactReqested DeleteModel) {
	artToDelete := bson.M{"url": artifactReqested.URL, "username": artifactReqested.Username}
	var entry RepositoryModel

	err := a.collection.FindOne(
		context.TODO(),
		artToDelete,
	).Decode(&entry)

	if err != nil {
		log.Fatal(err)
	}

	if entry == (RepositoryModel{}) {
		fmt.Println("Entry not found")
		return
	}

	result, err := a.collection.DeleteOne(context.TODO(), artToDelete)

	fmt.Printf("Remove %v document", result.DeletedCount)

}

// Dbdriver established the connection to our mongo db
func Dbdriver() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
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
	db := client.Database("mvral")
	migrate.SetDatabase(db)
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		fmt.Print(err)
	}

	return client
}
