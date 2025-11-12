package repository

import (
    "context"
    // "errors"
    "time"
    
    "go-fiber/app/model"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const fileCollection = "files"

type FileRepository struct {
    DB *mongo.Database
}

func NewFileRepository(db *mongo.Database) *FileRepository {
    return &FileRepository{DB: db}
}

func (r *FileRepository) Create(file *model.File) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    
    file.UploadedAt = time.Now()
    
    result, err := collection.InsertOne(ctx, file)
    if err != nil {
        return err
    }
    
    file.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (r *FileRepository) FindAll(search string, fileCategory string, limit, offset int) ([]model.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    
    // Build filter
    filter := bson.M{}
    if search != "" {
        filter["$or"] = []bson.M{
            {"original_name": bson.M{"$regex": search, "$options": "i"}},
            {"file_type": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    if fileCategory != "" {
        filter["file_category"] = fileCategory
    }
    
    opts := options.Find().
        SetSort(bson.D{{Key: "uploaded_at", Value: -1}}).
        SetLimit(int64(limit)).
        SetSkip(int64(offset))
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var files []model.File
    if err = cursor.All(ctx, &files); err != nil {
        return nil, err
    }

    return files, nil
}

func (r *FileRepository) FindByAlumniID(alumniID primitive.ObjectID, fileCategory string) ([]model.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    
    filter := bson.M{"alumni_id": alumniID}
    if fileCategory != "" {
        filter["file_category"] = fileCategory
    }
    
    opts := options.Find().SetSort(bson.D{{Key: "uploaded_at", Value: -1}})
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var files []model.File
    if err = cursor.All(ctx, &files); err != nil {
        return nil, err
    }

    return files, nil
}

func (r *FileRepository) Count(search string, fileCategory string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    
    filter := bson.M{}
    if search != "" {
        filter["$or"] = []bson.M{
            {"original_name": bson.M{"$regex": search, "$options": "i"}},
            {"file_type": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    if fileCategory != "" {
        filter["file_category"] = fileCategory
    }
    
    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    
    return int(count), nil
}

func (r *FileRepository) FindByID(id primitive.ObjectID) (*model.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    
    var file model.File
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
    if err != nil {
        return nil, err
    }

    return &file, nil
}

func (r *FileRepository) Delete(id primitive.ObjectID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    
    _, err := collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

// Check if user owns the file (for authorization)
func (r *FileRepository) IsFileOwnedByUser(fileID, userID primitive.ObjectID) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(fileCollection)
    alumniCollection := r.DB.Collection("alumni")
    
    // Get file
    var file model.File
    err := collection.FindOne(ctx, bson.M{"_id": fileID}).Decode(&file)
    if err != nil {
        return false, err
    }
    
    // Check if alumni belongs to user
    var alumni struct {
        UserID primitive.ObjectID `bson:"user_id"`
    }
    err = alumniCollection.FindOne(ctx, bson.M{"_id": file.AlumniID}).Decode(&alumni)
    if err != nil {
        return false, err
    }
    
    return alumni.UserID == userID, nil
}