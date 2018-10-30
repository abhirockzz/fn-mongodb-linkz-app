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

	fdk.Handle(fdk.HandlerFunc(deleteHandler))
}

func deleteHandler(ctx context.Context, in io.Reader, out io.Writer) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	linkID := buf.String()

	log.Println("Link to be deleted ", linkID)

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

		log.Println("Successfully connected to MongoDB..")
	}

	linkObjectID, err := objectid.FromHex(linkID)

	if err != nil {
		resp := Response{"FAILED", err.Error()}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
	}
	toBeDeleted := bson.EC.ObjectID("_id", linkObjectID)

	linkDb := fdxCtx.Config()["MONGODB_DB"]
	linkColl := fdxCtx.Config()["MONGODB_COLLECTION"]

	delResp, delErr := client.Database(linkDb).Collection(linkColl).DeleteOne(nil, bson.NewDocument(toBeDeleted))

	if delErr != nil {
		resp := Response{"FAILED", delErr.Error()}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
	}

	if delResp.DeletedCount == 0 {
		resp := Response{"FAILED", "Unable to delete document with ID " + linkID + " maybe because it does not exist"}
		json.NewEncoder(out).Encode(resp)
		log.Println(resp.failure())
		return
	}

	success := "Deleted successfully"
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
