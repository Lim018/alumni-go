package database

import (
    "context"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates necessary indexes for better query performance
func CreateIndexes(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Alumni indexes
    alumniCollection := db.Collection("alumni")
    alumniIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "nim", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys: bson.D{{Key: "email", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "nama", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "user_id", Value: 1}},
        },
    }
    
    _, err := alumniCollection.Indexes().CreateMany(ctx, alumniIndexes)
    if err != nil {
        log.Printf("Failed to create alumni indexes: %v", err)
        return err
    }
    log.Println("Alumni indexes created ✅")

    // Users indexes
    userCollection := db.Collection("users")
    userIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true),
        },
    }
    
    _, err = userCollection.Indexes().CreateMany(ctx, userIndexes)
    if err != nil {
        log.Printf("Failed to create user indexes: %v", err)
        return err
    }
    log.Println("User indexes created ✅")

    // Pekerjaan indexes
    pekerjaanCollection := db.Collection("pekerjaan_alumni")
    pekerjaanIndexes := []mongo.IndexModel{
        {
            Keys: bson.D{{Key: "alumni_id", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "nama_perusahaan", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "is_delete", Value: 1}},
        },
    }
    
    _, err = pekerjaanCollection.Indexes().CreateMany(ctx, pekerjaanIndexes)
    if err != nil {
        log.Printf("Failed to create pekerjaan indexes: %v", err)
        return err
    }
    log.Println("Pekerjaan indexes created ✅")

    return nil
}