package migrations

import (
	"context"
	"fmt"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Migration script for adding username field to non updated fields
func init() {
	migrate.Register(func(db *mongo.Database) error {
		update := bson.M{
			"$set": bson.M{
				"username": "Anonymous",
			},
		}
		filter := bson.M{}
		_, err := db.Collection("artifact").UpdateMany(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		fmt.Println("Done with mig")

		return nil
	}, func(db *mongo.Database) error {
		fmt.Println("err in mig")
		update := bson.M{
			"$unset": bson.M{
				"username": "",
			},
		}
		filter := bson.M{}
		_, err := db.Collection("artifact").UpdateMany(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		return nil

	})
}
