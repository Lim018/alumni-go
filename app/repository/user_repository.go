package repository

import (
    "context"
    "errors"
    "time"
    
    "go-fiber/app/model"
    
    "go.mongodb.org/mongo-driver/bson"
    // "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"

type UserRepository struct {
    DB *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) FindUserByUsernameOrEmail(identifier string) (*model.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(userCollection)
    
    filter := bson.M{
        "$or": []bson.M{
            {"username": identifier},
            {"email": identifier},
        },
    }
    
    var user model.User
    err := collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}

func (r *UserRepository) GetUsers(search, sortBy, order string, limit, offset int) ([]model.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(userCollection)
    
    // Build filter
    filter := bson.M{}
    if search != "" {
        filter["$or"] = []bson.M{
            {"username": bson.M{"$regex": search, "$options": "i"}},
            {"email": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    // Set sort
    sortOrder := 1
    if order == "DESC" || order == "desc" {
        sortOrder = -1
    }
    
    // Map sortBy to MongoDB fields
    sortByMap := map[string]string{
        "id":         "_id",
        "username":   "username",
        "email":      "email",
        "role":       "role",
        "created_at": "created_at",
    }
    
    mongoSortBy := sortByMap[sortBy]
    if mongoSortBy == "" {
        mongoSortBy = "_id"
    }
    
    opts := options.Find().
        SetSort(bson.D{{Key: mongoSortBy, Value: sortOrder}}).
        SetLimit(int64(limit)).
        SetSkip(int64(offset))
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var users []model.User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    
    return users, nil
}

func (r *UserRepository) CountUsers(search string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(userCollection)
    
    filter := bson.M{}
    if search != "" {
        filter["$or"] = []bson.M{
            {"username": bson.M{"$regex": search, "$options": "i"}},
            {"email": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    
    return int(count), nil
}

// Legacy functions for compatibility
func FindUserByUsernameOrEmail(db *mongo.Database, identifier string) (*model.User, string, error) {
    repo := NewUserRepository(db)
    user, err := repo.FindUserByUsernameOrEmail(identifier)
    if err != nil {
        return nil, "", err
    }
    return user, user.PasswordHash, nil
}

func GetUsersRepo(db *mongo.Database, search, sortBy, order string, limit, offset int) ([]model.User, error) {
    repo := NewUserRepository(db)
    return repo.GetUsers(search, sortBy, order, limit, offset)
}

func CountUsersRepo(db *mongo.Database, search string) (int, error) {
    repo := NewUserRepository(db)
    return repo.CountUsers(search)
}