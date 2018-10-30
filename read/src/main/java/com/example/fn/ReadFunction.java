package com.example.fn;

import com.fasterxml.jackson.annotation.JsonInclude;

import com.fnproject.fn.api.RuntimeContext;

import com.mongodb.ConnectionString;
import com.mongodb.client.FindIterable;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.model.Filters;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.function.Consumer;

import org.bson.Document;

public class ReadFunction {

    RuntimeContext ctx;
    MongoCollection<Document> coll;

    public ReadFunction(RuntimeContext ctx) {
        this.ctx = ctx;
    }
    
    public static class Filter {
        public String name;
        public String value;

        public Filter() {
        }
    }

    public List<Link> handleRequest(Filter filter) {
        String linksDB = ctx.getConfiguration().get("MONGODB_DB");
        String linksCollection = ctx.getConfiguration().get("MONGODB_COLLECTION");

        //only if coll is null i.e. this will NOT be invoked if the same function (Docker) container is being used
        if (coll == null) {
            String mongoConnectString = ctx.getConfiguration().get("MONGODB_CONNECT_STRING");

            System.err.println("Connecting to MongoDB... " + mongoConnectString);
            try {
                MongoClient mongoClient = MongoClients.create(new ConnectionString(mongoConnectString));
                coll = mongoClient.getDatabase(linksDB).getCollection(linksCollection);
                System.err.println("Successfully connected to MongoDB");
            } catch (Exception e) {
                String response = "Failed to connect to MongoDB " + e.getMessage();
                System.err.println(response);
                return Collections.emptyList();
            }
        }

        List<Link> links = new ArrayList<>();
        FindIterable<Document> linkDocs = null;

        if (filter.name.equals("ALL")) {
            //find ALL
            System.err.println("Listing all links");
            linkDocs = coll.find();
        } else if (filter.name.equals("users")) {
            System.err.println("Finding links posted by users " + filter.value);
            //find by user
            linkDocs = coll.find(Filters.in("postedBy", Arrays.asList(filter.value.split(","))));
        } else if (filter.name.equals("tags")) {
            System.err.println("Finding links containing tags " + filter.value);
            //find by tags
            linkDocs = coll.find(Filters.in("tags", Arrays.asList(filter.value.split(","))));
        }

        if (linkDocs!= null) {
            linkDocs.forEach(new Consumer<Document>() {
                @Override
                public void accept(Document d) {
                    links.add(convert(d));
                }
            });
        }
        System.err.println("Found " + links.size() + " links......");
        return links;
    }

    static Link convert(Document link) {

        Link _link = new Link();
        _link.link = link.getString("link");
        _link.headline = link.getString("headline");
        _link.postedBy = link.getString("postedBy");

        if (link.containsKey("upvotes")) {
            _link.upvotes = link.getInteger("upvotes");
        }

        if (link.containsKey("comments")) {
            _link.comments = link.get("comments", List.class);
        }
        
         if (link.containsKey("tags")) {
            _link.tags = link.get("tags", List.class);
        }

        System.err.println("Found link " + _link);
        return _link;
    }

    @JsonInclude(JsonInclude.Include.NON_NULL) 
    public static class Link {

        public String link;
        public String headline;
        public String postedBy;
        public List<Comment> comments;
        public Integer upvotes;
        public List<String> tags;

        public Link() {
        }

    }

    public static class Comment {

        public String comment;
        public String user;

        public Comment() {
        }

        public Comment(String comment, String user) {
            this.comment = comment;
            this.user = user;
        }

    }

}
