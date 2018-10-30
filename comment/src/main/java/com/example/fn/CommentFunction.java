package com.example.fn;

import com.fnproject.fn.api.RuntimeContext;

import com.mongodb.ConnectionString;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.model.Filters;
import com.mongodb.client.model.Updates;
import com.mongodb.client.result.UpdateResult;
import org.bson.Document;
import org.bson.types.ObjectId;

public class CommentFunction {

    RuntimeContext ctx;
    MongoCollection<Document> coll;

    public CommentFunction(RuntimeContext ctx) {
        this.ctx = ctx;
    }

    public String handleRequest(CommentDetail commentDetail) {
        String linksDB = ctx.getConfiguration().get("MONGODB_DB");
        String linksCollection = ctx.getConfiguration().get("MONGODB_COLLECTION");

        //only if coll is null i.e. this will NOT be invoked if the same function (Docker) container is being used
        if (coll == null) {
            String mongoConnectString = ctx.getConfiguration().get("MONGODB_CONNECT_STRING");

            System.err.println("Connecting to MongoDB " + mongoConnectString);
            try {
                MongoClient mongoClient = MongoClients.create(new ConnectionString(mongoConnectString));
                coll = mongoClient.getDatabase(linksDB).getCollection(linksCollection);
                System.err.println("Successfully connected to MongoDB");
            } catch (Exception e) {
                String response = "Failed to connect to MongoDB - " + e.getMessage();
                System.err.println(response);
                return response;
            }
        }

        System.err.println("Adding comment " + commentDetail);

        UpdateResult result = null;
        try {
            result = coll.updateOne(Filters.eq("_id", new ObjectId(commentDetail.linkDocId)), Updates.push("comments", new Document("comment", commentDetail.comment)
                    .append("user", commentDetail.user)));
            System.err.println("Result "+ result.getModifiedCount());
        } catch (Exception e) {
            String response = "Error while adding comment - " + e.getMessage();
            System.err.println(response);
            return response;
        }
        
        String finalResp = (result.getModifiedCount() == 0) ? "Failed to add comment. No document with ID "+ commentDetail.linkDocId : "Successfully added comment for document " + commentDetail.linkDocId;

        //System.err.println("Successfully updated comment " + commentDetail.linkDocId);
        return finalResp;
    }

    public static class CommentDetail {

        public String linkDocId;
        public String comment;
        public String user;

        public CommentDetail() {
        }

    }

}
