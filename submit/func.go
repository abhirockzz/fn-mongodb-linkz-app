package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"

	fdk "github.com/fnproject/fdk-go"
)

var client *mongo.Client

func main() {

	fdk.Handle(fdk.HandlerFunc(createLinkHandler))
}

func createLinkHandler(ctx context.Context, in io.Reader, out io.Writer) {

	var linkJSON linkInfo
	json.NewDecoder(in).Decode(&linkJSON)
	log.Println("Link info ", linkJSON)
	linkDoc := bson.NewDocument(bson.EC.String("link", linkJSON.Link), bson.EC.String("headline", linkJSON.Headline), bson.EC.String("postedBy", linkJSON.PostedBy))

	if len(linkJSON.Tags) > 0 {
		var tags []*bson.Value
		log.Println("processing tag info")
		for _, tag := range linkJSON.Tags {
			tags = append(tags, bson.VC.String(tag))
			log.Println("processed tag ", tag)
		}

		linkDoc = linkDoc.Append(bson.EC.ArrayFromElements("tags", tags...))
	}
	fdxCtx := fdk.GetContext(ctx)

	if client == nil {
		connString := fdxCtx.Config()["MONGODB_CONNECT_STRING"]
		log.Println("Connecting to MongoDB ", connString)

		client, _ = mongo.NewClient(connString)
		err := client.Connect(context.Background())

		if err != nil {
			resp := Response{"FAILED", err.Error()}
			json.NewEncoder(out).Encode(resp)
			log.Println(resp.failure())
			return
		}

		log.Println("Successfully connected to MongoDB")
	}

	linkDb := fdxCtx.Config()["MONGODB_DB"]
	linkColl := fdxCtx.Config()["MONGODB_COLLECTION"]

	insertResult, insertErr := client.Database(linkDb).Collection(linkColl).InsertOne(context.Background(), linkDoc)
	if insertErr != nil {
		resp := Response{"FAILED", insertErr.Error()}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
	}
	insertedILinkID := insertResult.InsertedID.(objectid.ObjectID).Hex()
	success := "Successfully inserted link " + insertedILinkID
	resp := Response{"SUCCESS", insertedILinkID}
	json.NewEncoder(out).Encode(resp)
	log.Println(success)

}

type linkInfo struct {
	Link     string   `json:"link"`
	Headline string   `json:"headline"`
	PostedBy string   `json:"postedBy"`
	Tags     []string `json:"tags,omitempty"`
}

//Response ...
type Response struct {
	Status, Message string
}

func (r Response) failure() string {
	return "FAILED due to " + r.Message
}
