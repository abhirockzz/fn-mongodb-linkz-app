package main

import (
	"bytes"
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

	fdk.Handle(fdk.HandlerFunc(upvoteHandler))
}

func upvoteHandler(ctx context.Context, in io.Reader, out io.Writer) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	linkID := buf.String()

	log.Println("Link to be updated ", linkID)

	if linkID == "" {
		resp := Response{Status: "FAILED", Message: "Link ID cannot be blank"}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
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

	linkObjectID, err := objectid.FromHex(linkID)

	if err != nil {
		resp := Response{"FAILED", err.Error()}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
	}
	toBeUpvoted := bson.EC.ObjectID("_id", linkObjectID)
	upvote := bson.NewDocument(bson.EC.SubDocumentFromElements("$inc", bson.EC.Int32("upvotes", 1)))

	linkDb := fdxCtx.Config()["MONGODB_DB"]
	linkColl := fdxCtx.Config()["MONGODB_COLLECTION"]

	_, updateErr := client.Database(linkDb).Collection(linkColl).UpdateOne(nil, bson.NewDocument(toBeUpvoted), upvote)
	if updateErr != nil {
		resp := Response{"FAILED", updateErr.Error()}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
	}

	success := "Upvoted successfully"
	resp := Response{"SUCCESS", success}
	json.NewEncoder(out).Encode(resp)
	log.Println(success)

}

//Response ...
type Response struct {
	Status, Message string
}

func (r Response) failure() string {
	return "FAILED due to " + r.Message
}
