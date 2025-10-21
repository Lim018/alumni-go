package repository

import (
    "context"
    "errors"
    "time"
    
    "go-fiber/app/model"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const pekerjaanCollection = "pekerjaan_alumni"

type PekerjaanRepository struct {
    DB *mongo.Database
}

func NewPekerjaanRepository(db *mongo.Database) *PekerjaanRepository {
    return &PekerjaanRepository{DB: db}
}

func (r *PekerjaanRepository) CreatePekerjaan(p model.Pekerjaan) (*model.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    p.CreatedAt = time.Now()
    p.UpdatedAt = time.Now()
    
    result, err := collection.InsertOne(ctx, p)
    if err != nil {
        return nil, err
    }
    
    p.ID = result.InsertedID.(primitive.ObjectID)
    return &p, nil
}

func (r *PekerjaanRepository) UpdatePekerjaan(id primitive.ObjectID, p model.Pekerjaan) (*model.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    p.UpdatedAt = time.Now()
    
    update := bson.M{
        "$set": bson.M{
            "alumni_id":             p.AlumniID,
            "nama_perusahaan":       p.NamaPerusahaan,
            "posisi_jabatan":        p.PosisiJabatan,
            "bidang_industri":       p.BidangIndustri,
            "lokasi_kerja":          p.LokasiKerja,
            "gaji_range":            p.GajiRange,
            "tanggal_mulai_kerja":   p.TanggalMulaiKerja,
            "tanggal_selesai_kerja": p.TanggalSelesaiKerja,
            "status_pekerjaan":      p.StatusPekerjaan,
            "deskripsi_pekerjaan":   p.DeskripsiPekerjaan,
            "updated_at":            p.UpdatedAt,
        },
    }
    
    filter := bson.M{"_id": id}
    _, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return nil, err
    }
    
    p.ID = id
    return &p, nil
}

func (r *PekerjaanRepository) FindPekerjaanByAlumniID(alumniID primitive.ObjectID) ([]model.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    filter := bson.M{
        "alumni_id": alumniID,
        "is_delete": bson.M{"$exists": false},
    }
    
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var list []model.Pekerjaan
    if err = cursor.All(ctx, &list); err != nil {
        return nil, err
    }
    
    return list, nil
}

func (r *PekerjaanRepository) GetPekerjaan(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    // Build filter
    filter := bson.M{"is_delete": bson.M{"$exists": false}}
    if search != "" {
        filter["$or"] = []bson.M{
            {"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
            {"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
            {"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
            {"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    // Set sort
    sortOrder := 1
    if order == "desc" {
        sortOrder = -1
    }
    
    allowedSort := map[string]bool{"_id": true, "nama_perusahaan": true, "posisi_jabatan": true, "tanggal_mulai_kerja": true}
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
    
    var list []model.Pekerjaan
    if err = cursor.All(ctx, &list); err != nil {
        return nil, err
    }
    
    return list, nil
}

func (r *PekerjaanRepository) CountPekerjaan(search string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    filter := bson.M{"is_delete": bson.M{"$exists": false}}
    if search != "" {
        filter["$or"] = []bson.M{
            {"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
            {"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
            {"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
            {"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    
    return int(count), nil
}

func (r *PekerjaanRepository) SoftDelete(id primitive.ObjectID, userID primitive.ObjectID, isAdmin bool) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    now := time.Now()
    
    var filter bson.M
    if isAdmin {
        filter = bson.M{"_id": id}
    } else {
        // Check if pekerjaan belongs to alumni owned by user
        alumniCollection := r.DB.Collection("alumni")
        var alumni struct {
            ID primitive.ObjectID `bson:"_id"`
        }
        
        err := alumniCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&alumni)
        if err != nil {
            return errors.New("alumni tidak ditemukan")
        }
        
        filter = bson.M{
            "_id":      id,
            "alumni_id": alumni.ID,
        }
    }
    
    update := bson.M{"$set": bson.M{"is_delete": now}}
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return errors.New("data tidak ditemukan atau tidak memiliki akses")
    }
    
    return nil
}

func (r *PekerjaanRepository) GetTrashPekerjaan(userID primitive.ObjectID, isAdmin bool, search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    // Build filter
    filter := bson.M{"is_delete": bson.M{"$exists": true}}
    
    if !isAdmin {
        // Get alumni for this user
        alumniCollection := r.DB.Collection("alumni")
        var alumni struct {
            ID primitive.ObjectID `bson:"_id"`
        }
        
        err := alumniCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&alumni)
        if err != nil {
            return []model.Pekerjaan{}, nil
        }
        
        filter["alumni_id"] = alumni.ID
    }
    
    if search != "" {
        filter["$or"] = []bson.M{
            {"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
            {"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    // Set sort
    sortOrder := -1
    if order == "asc" {
        sortOrder = 1
    }
    
    allowedSort := map[string]bool{"_id": true, "nama_perusahaan": true, "is_delete": true}
    if !allowedSort[sortBy] {
        sortBy = "is_delete"
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
    
    var list []model.Pekerjaan
    if err = cursor.All(ctx, &list); err != nil {
        return nil, err
    }
    
    return list, nil
}

func (r *PekerjaanRepository) CountTrashPekerjaan(userID primitive.ObjectID, isAdmin bool, search string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    filter := bson.M{"is_delete": bson.M{"$exists": true}}
    
    if !isAdmin {
        alumniCollection := r.DB.Collection("alumni")
        var alumni struct {
            ID primitive.ObjectID `bson:"_id"`
        }
        
        err := alumniCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&alumni)
        if err != nil {
            return 0, nil
        }
        
        filter["alumni_id"] = alumni.ID
    }
    
    if search != "" {
        filter["$or"] = []bson.M{
            {"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
            {"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
        }
    }
    
    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return 0, err
    }
    
    return int(count), nil
}

func (r *PekerjaanRepository) RestorePekerjaan(id, userID primitive.ObjectID, isAdmin bool) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    var filter bson.M
    if isAdmin {
        filter = bson.M{
            "_id":       id,
            "is_delete": bson.M{"$exists": true},
        }
    } else {
        alumniCollection := r.DB.Collection("alumni")
        var alumni struct {
            ID primitive.ObjectID `bson:"_id"`
        }
        
        err := alumniCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&alumni)
        if err != nil {
            return errors.New("alumni tidak ditemukan")
        }
        
        filter = bson.M{
            "_id":       id,
            "alumni_id": alumni.ID,
            "is_delete": bson.M{"$exists": true},
        }
    }
    
    update := bson.M{"$unset": bson.M{"is_delete": ""}}
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    
    if result.MatchedCount == 0 {
        return errors.New("data tidak ditemukan atau tidak memiliki akses")
    }
    
    return nil
}

func (r *PekerjaanRepository) HardDeletePekerjaan(id, userID primitive.ObjectID, isAdmin bool) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := r.DB.Collection(pekerjaanCollection)
    
    var filter bson.M
    if isAdmin {
        filter = bson.M{
            "_id":       id,
            "is_delete": bson.M{"$exists": true},
        }
    } else {
        alumniCollection := r.DB.Collection("alumni")
        var alumni struct {
            ID primitive.ObjectID `bson:"_id"`
        }
        
        err := alumniCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&alumni)
        if err != nil {
            return errors.New("alumni tidak ditemukan")
        }
        
        filter = bson.M{
            "_id":       id,
            "alumni_id": alumni.ID,
            "is_delete": bson.M{"$exists": true},
        }
    }
    
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    
    if result.DeletedCount == 0 {
        return errors.New("data tidak ditemukan atau tidak memiliki akses")
    }
    
    return nil
}