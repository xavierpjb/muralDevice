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

type IArtifactRepositoryHandler interface {
	Create(ArtifactRepositoryModel)
	RetrieveList(int64) []ArtifactRepositoryModel
}

type ArtifactRepositoryModel struct {
	// ID             primitive.ObjectID `bson: "_id,omitempty"`
	URL            string    `json:"url" bson:"url,omitempty"`
	FileType       string    `json:"fileType" bson:"fileType,omitempty"`
	UploadDateTime time.Time `json:"uploadDateTime" bson:"uploadDateTime,omitempty"`
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

func (a ArtifactRepositoryHandler) RetrieveList(page int64) []ArtifactRepositoryModel {
	filter := bson.M{}
	paginatedData, err := mongopagination.New(a.collection).Limit(5).Page(page).Sort("uploadDateTime", -1).Filter(filter).Find()
	if err != nil {
		panic(err)
	}

	var entries []ArtifactRepositoryModel
	for _, raw := range paginatedData.Data {
		var art *ArtifactRepositoryModel
		if marshallErr := bson.Unmarshal(raw, &art); marshallErr == nil {
			entries = append(entries, *art)
		}
	}

	fmt.Println("entries found")
	fmt.Println(entries)
	return entries
}

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
