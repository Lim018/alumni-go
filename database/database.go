package database

import (
    "context"
    "log"
    "os"
    "time"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() *mongo.Database {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    mongoURI := os.Getenv("MONGODB_URI")
    dbName := os.Getenv("MONGODB_DATABASE")
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping database
    if err = client.Ping(ctx, nil); err != nil {
        log.Fatal("Database tidak connect:", err)
    }

    log.Println("DB Connected ✅")
    DB = client.Database(dbName)
    return DB
}

func GetCollection(collectionName string) *mongo.Collection {
    return DB.Collection(collectionName)
}