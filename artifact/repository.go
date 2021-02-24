package artifact

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type IArtifactRepositoryHandler interface {
	Create(ArtifactRepositoryModel)
	RetrieveList()
}

type ArtifactRepositoryModel struct {
	// ID             primitive.ObjectID `bson: "_id,omitempty"`
	URL            string    `bson: "url,omitempty"`
	FileType       string    `bson: "fileType,omitempty"`
	UploadDateTime time.Time `bson: "uploadDateTime,omitempty"`
}

type ArtifactRepositoryHandler struct {
	collection *mongo.Collection
}

func NewARH(client *mongo.Client) ArtifactRepositoryHandler {

	col := client.Database("mvral").Collection("artifact")
	a := ArtifactRepositoryHandler{col}
	return a
}

func (a ArtifactRepositoryHandler) Create(artifactPersisted ArtifactRepositoryModel) {
	// this will be correctly filled in once feature for creating entry done
	_, err := a.collection.InsertOne(context.TODO(), artifactPersisted)
	if err != nil {
		log.Fatalln(err)
	}
}

func (a ArtifactRepositoryHandler) RetrieveList() {
	cursor, err := a.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var entries []bson.M
	if err = cursor.All(context.TODO(), &entries); err != nil {
		log.Fatal(err)
	}
	fmt.Println("entries found")
	fmt.Println(entries)
}

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

	return client
}
