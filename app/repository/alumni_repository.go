package repository

import (
    "context"
    "time"
    
    "go-fiber/app/model"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const alumniCollection = "alumni"

type AlumniRepository struct {
    DB *mongo.Database
}

func NewAlumniRepository(db *mongo.Database) *AlumniRepository {
    return &AlumniRepository{DB: db}
}

func (r *AlumniRepository) CreateAlumni(alumni model.Alumni) (*model.Alumni, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(alumniCollection)
    
    alumni.CreatedAt = time.Now()
    alumni.UpdatedAt = time.Now()
    
    result, err := collection.InsertOne(ctx, alumni)
    if err != nil {
        return nil, err
    }
    
    alumni.ID = result.InsertedID.(primitive.ObjectID)
    return &alumni, nil
}

func (r *AlumniRepository) UpdateAlumni(id primitive.ObjectID, alumni model.Alumni) (*model.Alumni, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(alumniCollection)
    
    alumni.UpdatedAt = time.Now()
    
    update := bson.M{
        "$set": bson.M{
            "nim":         alumni.NIM,
            "nama":        alumni.Nama,
            "jurusan":     alumni.Jurusan,
            "angkatan":    alumni.Angkatan,
            "tahun_lulus": alumni.TahunLulus,
            "email":       alumni.Email,
            "no_telepon":  alumni.NoTelepon,
            "alamat":      alumni.Alamat,
            "updated_at":  alumni.UpdatedAt,
        },
    }
    
    filter := bson.M{"_id": id}
    _, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return nil, err
    }
    
    alumni.ID = id
    return &alumni, nil
}

func (r *AlumniRepository) DeleteAlumni(id primitive.ObjectID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(alumniCollection)
    
    _, err := collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *AlumniRepository) GetAlumni(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(alumniCollection)
    
    // Build filter
    filter := bson.M{}
    if search != "" {
        filter["$or"] = []bson.M{
            {"nama": bson.M{"$regex": search, "$options": "i"}},
            {"nim": bson.M{"$regex": search, "$options": "i"}},
            {"email": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    // Set sort
    sortOrder := 1
    if order == "desc" {
        sortOrder = -1
    }
    
    allowedSort := map[string]bool{"_id": true, "nama": true, "angkatan": true, "tahun_lulus": true}
    if !allowedSort[sortBy] {
        sortBy = "_id"
    }
    
    opts := options.Find().
        SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
        SetLimit(int64(limit)).
        SetSkip(int64(offset))
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var list []model.Alumni
    if err = cursor.All(ctx, &list); err != nil {
        return nil, err
    }
    
    return list, nil
}

func (r *AlumniRepository) CountAlumni(search string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(alumniCollection)
    
    filter := bson.M{}
    if search != "" {
        filter["$or"] = []bson.M{
            {"nama": bson.M{"$regex": search, "$options": "i"}},
            {"nim": bson.M{"$regex": search, "$options": "i"}},
            {"email": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    
    return int(count), nil
}

func (r *AlumniRepository) GetAlumniStatsByJurusan() ([]model.AlumniStatsByJurusanResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(alumniCollection)
    
    pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$jurusan"},
            {Key: "total", Value: bson.D{{Key: "$sum", Value: 1}}},
        }}},
        {{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}},
        {{Key: "$project", Value: bson.D{
            {Key: "jurusan", Value: "$_id"},
            {Key: "total", Value: 1},
            {Key: "_id", Value: 0},
        }}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var stats []model.AlumniStatsByJurusanResponse
    if err = cursor.All(ctx, &stats); err != nil {
        return nil, err
    }
    
    return stats, nil
}

// Legacy functions for compatibility (will be deprecated)
func CreateAlumni(db *mongo.Database, alumni model.Alumni) (*model.Alumni, error) {
    repo := NewAlumniRepository(db)
    return repo.CreateAlumni(alumni)
}

func UpdateAlumni(db *mongo.Database, id primitive.ObjectID, alumni model.Alumni) (*model.Alumni, error) {
    repo := NewAlumniRepository(db)
    return repo.UpdateAlumni(id, alumni)
}

func DeleteAlumni(db *mongo.Database, id primitive.ObjectID) error {
    repo := NewAlumniRepository(db)
    return repo.DeleteAlumni(id)
}

func GetAlumniRepo(db *mongo.Database, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
    repo := NewAlumniRepository(db)
    return repo.GetAlumni(search, sortBy, order, limit, offset)
}

func CountAlumniRepo(db *mongo.Database, search string) (int, error) {
    repo := NewAlumniRepository(db)
    return repo.CountAlumni(search)
}

func GetAlumniStatsByJurusan(db *mongo.Database) ([]model.AlumniStatsByJurusanResponse, error) {
    repo := NewAlumniRepository(db)
    return repo.GetAlumniStatsByJurusan()
}