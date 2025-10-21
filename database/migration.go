package database

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Collection names
const (
    UsersCollection      = "users"
    AlumniCollection     = "alumni"
    PekerjaanCollection  = "pekerjaan_alumni"
    MigrationsCollection = "migrations"
)

// Migration represents a database migration
type Migration struct {
    Name      string    `bson:"name"`
    AppliedAt time.Time `bson:"applied_at"`
}

// RunMigrations executes all database migrations
func RunMigrations(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    log.Println("🔄 Starting migrations...")

    // Check if migrations collection exists
    migrationsColl := db.Collection(MigrationsCollection)

    migrations := []struct {
        name string
        fn   func(*mongo.Database) error
    }{
        {"create_users_collection", createUsersCollection},
        {"create_alumni_collection", createAlumniCollection},
        {"create_pekerjaan_collection", createPekerjaanCollection},
        {"create_indexes", createAllIndexes},
    }

    for _, migration := range migrations {
        // Check if migration already applied
        count, err := migrationsColl.CountDocuments(ctx, bson.M{"name": migration.name})
        if err != nil {
            return err
        }

        if count > 0 {
            log.Printf("⏭️  Migration '%s' already applied, skipping...", migration.name)
            continue
        }

        // Run migration
        log.Printf("▶️  Running migration: %s", migration.name)
        if err := migration.fn(db); err != nil {
            log.Printf("❌ Migration '%s' failed: %v", migration.name, err)
            return err
        }

        // Record migration
        _, err = migrationsColl.InsertOne(ctx, Migration{
            Name:      migration.name,
            AppliedAt: time.Now(),
        })
        if err != nil {
            return err
        }

        log.Printf("✅ Migration '%s' completed", migration.name)
    }

    log.Println("✅ All migrations completed successfully!")
    return nil
}

// createUsersCollection creates users collection with validation
func createUsersCollection(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Create collection with schema validation
    validator := bson.M{
        "$jsonSchema": bson.M{
            "bsonType": "object",
            "required": []string{"username", "email", "password_hash", "role", "created_at"},
            "properties": bson.M{
                "username": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                    "minLength":   3,
                    "maxLength":   50,
                },
                "email": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and match email pattern",
                    "pattern":     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
                },
                "password_hash": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                },
                "role": bson.M{
                    "bsonType":    "string",
                    "description": "must be either admin or user",
                    "enum":        []string{"admin", "user"},
                },
                "created_at": bson.M{
                    "bsonType":    "date",
                    "description": "must be a date and is required",
                },
            },
        },
    }

    opts := options.CreateCollection().SetValidator(validator)
    return db.CreateCollection(ctx, UsersCollection, opts)
}

// createAlumniCollection creates alumni collection with validation
func createAlumniCollection(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    validator := bson.M{
        "$jsonSchema": bson.M{
            "bsonType": "object",
            "required": []string{"nim", "nama", "jurusan", "angkatan", "tahun_lulus", "email", "no_telepon", "created_at", "updated_at"},
            "properties": bson.M{
                "nim": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                    "minLength":   5,
                    "maxLength":   20,
                },
                "nama": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                    "minLength":   3,
                    "maxLength":   100,
                },
                "jurusan": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                },
                "angkatan": bson.M{
                    "bsonType":    "int",
                    "description": "must be an integer and is required",
                    "minimum":     1900,
                    "maximum":     2100,
                },
                "tahun_lulus": bson.M{
                    "bsonType":    "int",
                    "description": "must be an integer and is required",
                    "minimum":     1900,
                    "maximum":     2100,
                },
                "email": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and match email pattern",
                    "pattern":     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
                },
                "no_telepon": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                },
                "alamat": bson.M{
                    "bsonType":    []string{"string", "null"},
                    "description": "must be a string or null",
                },
                "user_id": bson.M{
                    "bsonType":    []string{"objectId", "null"},
                    "description": "must be an objectId reference to users or null",
                },
                "created_at": bson.M{
                    "bsonType":    "date",
                    "description": "must be a date and is required",
                },
                "updated_at": bson.M{
                    "bsonType":    "date",
                    "description": "must be a date and is required",
                },
            },
        },
    }

    opts := options.CreateCollection().SetValidator(validator)
    return db.CreateCollection(ctx, AlumniCollection, opts)
}

func createPekerjaanCollection(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    validator := bson.M{
        "$jsonSchema": bson.M{
            "bsonType": "object",
            "required": []string{"alumni_id", "nama_perusahaan", "posisi_jabatan", "bidang_industri", "lokasi_kerja", "tanggal_mulai_kerja", "status_pekerjaan", "created_at", "updated_at"},
            "properties": bson.M{
                "alumni_id": bson.M{
                    "bsonType":    "objectId",
                    "description": "must be an objectId reference to alumni and is required",
                },
                "nama_perusahaan": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                    "minLength":   2,
                    "maxLength":   200,
                },
                "posisi_jabatan": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                },
                "bidang_industri": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                },
                "lokasi_kerja": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                },
                "gaji_range": bson.M{
                    "bsonType":    []string{"string", "null"},
                    "description": "must be a string or null",
                },
                "tanggal_mulai_kerja": bson.M{
                    "bsonType":    "date",
                    "description": "must be a date and is required",
                },
                "tanggal_selesai_kerja": bson.M{
                    "bsonType":    []string{"date", "null"}, // Allow null values
                    "description": "must be a date or null",
                },
                "status_pekerjaan": bson.M{
                    "bsonType":    "string",
                    "description": "must be a string and is required",
                    "enum":        []string{"aktif", "resign", "kontrak_habis"},
                },
                "deskripsi_pekerjaan": bson.M{
                    "bsonType":    []string{"string", "null"},
                    "description": "must be a string or null",
                },
                "is_delete": bson.M{
                    "bsonType":    []string{"date", "null"}, // Allow null values
                    "description": "soft delete timestamp",
                },
                "created_at": bson.M{
                    "bsonType":    "date",
                    "description": "must be a date and is required",
                },
                "updated_at": bson.M{
                    "bsonType":    "date",
                    "description": "must be a date and is required",
                },
            },
        },
    }

    opts := options.CreateCollection().SetValidator(validator)
    return db.CreateCollection(ctx, PekerjaanCollection, opts)
}

// createAllIndexes creates all necessary indexes
func createAllIndexes(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Users indexes
    userCollection := db.Collection(UsersCollection)
    userIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "username", Value: 1}},
            Options: options.Index().SetUnique(true).SetName("idx_username"),
        },
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetUnique(true).SetName("idx_email"),
        },
        {
            Keys:    bson.D{{Key: "role", Value: 1}},
            Options: options.Index().SetName("idx_role"),
        },
    }

    _, err := userCollection.Indexes().CreateMany(ctx, userIndexes)
    if err != nil {
        return err
    }
    log.Println("  ✓ Users indexes created")

    // Alumni indexes
    alumniCollection := db.Collection(AlumniCollection)
    alumniIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "nim", Value: 1}},
            Options: options.Index().SetUnique(true).SetName("idx_nim"),
        },
        {
            Keys:    bson.D{{Key: "email", Value: 1}},
            Options: options.Index().SetName("idx_alumni_email"),
        },
        {
            Keys:    bson.D{{Key: "nama", Value: 1}},
            Options: options.Index().SetName("idx_nama"),
        },
        {
            Keys:    bson.D{{Key: "jurusan", Value: 1}},
            Options: options.Index().SetName("idx_jurusan"),
        },
        {
            Keys:    bson.D{{Key: "angkatan", Value: 1}},
            Options: options.Index().SetName("idx_angkatan"),
        },
        {
            Keys:    bson.D{{Key: "tahun_lulus", Value: 1}},
            Options: options.Index().SetName("idx_tahun_lulus"),
        },
        {
            Keys:    bson.D{{Key: "user_id", Value: 1}},
            Options: options.Index().SetName("idx_user_id"),
        },
        {
            Keys: bson.D{
                {Key: "nama", Value: "text"},
                {Key: "nim", Value: "text"},
                {Key: "email", Value: "text"},
            },
            Options: options.Index().SetName("idx_text_search"),
        },
    }

    _, err = alumniCollection.Indexes().CreateMany(ctx, alumniIndexes)
    if err != nil {
        return err
    }
    log.Println("  ✓ Alumni indexes created")

    // Pekerjaan indexes
    pekerjaanCollection := db.Collection(PekerjaanCollection)
    pekerjaanIndexes := []mongo.IndexModel{
        {
            Keys:    bson.D{{Key: "alumni_id", Value: 1}},
            Options: options.Index().SetName("idx_alumni_id"),
        },
        {
            Keys:    bson.D{{Key: "nama_perusahaan", Value: 1}},
            Options: options.Index().SetName("idx_nama_perusahaan"),
        },
        {
            Keys:    bson.D{{Key: "bidang_industri", Value: 1}},
            Options: options.Index().SetName("idx_bidang_industri"),
        },
        {
            Keys:    bson.D{{Key: "status_pekerjaan", Value: 1}},
            Options: options.Index().SetName("idx_status_pekerjaan"),
        },
        {
            Keys:    bson.D{{Key: "is_delete", Value: 1}},
            Options: options.Index().SetName("idx_is_delete"),
        },
        {
            Keys:    bson.D{{Key: "tanggal_mulai_kerja", Value: -1}},
            Options: options.Index().SetName("idx_tanggal_mulai_kerja"),
        },
        {
            Keys: bson.D{
                {Key: "nama_perusahaan", Value: "text"},
                {Key: "posisi_jabatan", Value: "text"},
                {Key: "lokasi_kerja", Value: "text"},
            },
            Options: options.Index().SetName("idx_pekerjaan_text_search"),
        },
        // Compound index for common queries
        {
            Keys: bson.D{
                {Key: "alumni_id", Value: 1},
                {Key: "is_delete", Value: 1},
            },
            Options: options.Index().SetName("idx_alumni_active"),
        },
    }

    _, err = pekerjaanCollection.Indexes().CreateMany(ctx, pekerjaanIndexes)
    if err != nil {
        return err
    }
    log.Println("  ✓ Pekerjaan indexes created")

    return nil
}

// DropAllCollections drops all collections (for testing/reset)
func DropAllCollections(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collections := []string{
        UsersCollection,
        AlumniCollection,
        PekerjaanCollection,
        MigrationsCollection,
    }

    for _, coll := range collections {
        if err := db.Collection(coll).Drop(ctx); err != nil {
            log.Printf("Warning: Failed to drop collection %s: %v", coll, err)
        } else {
            log.Printf("✓ Dropped collection: %s", coll)
        }
    }

    return nil
}