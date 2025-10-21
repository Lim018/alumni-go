package database

import (
    "context"
    "log"
    "time"

    "go-fiber/utils"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// SeedData seeds initial data into the database
func SeedData(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    log.Println("🌱 Starting seeding...")

    // Seed Users first
    userIDs, err := seedUsers(ctx, db)
    if err != nil {
        return err
    }

    // Seed Alumni
    alumniIDs, err := seedAlumni(ctx, db, userIDs)
    if err != nil {
        return err
    }

    // Seed Pekerjaan
    if err := seedPekerjaan(ctx, db, alumniIDs); err != nil {
        return err
    }

    log.Println("✅ Seeding completed successfully!")
    return nil
}

// seedUsers seeds user data
func seedUsers(ctx context.Context, db *mongo.Database) (map[string]primitive.ObjectID, error) {
    collection := db.Collection(UsersCollection)

    // Check if users already exist
    count, err := collection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    if count > 0 {
        log.Println("  ⏭️  Users already seeded, fetching existing data...")
        
        // Fetch existing users
        cursor, err := collection.Find(ctx, bson.M{})
        if err != nil {
            return nil, err
        }
        defer cursor.Close(ctx)

        userIDs := make(map[string]primitive.ObjectID)
        for cursor.Next(ctx) {
            var user bson.M
            if err := cursor.Decode(&user); err != nil {
                return nil, err
            }
            userIDs[user["username"].(string)] = user["_id"].(primitive.ObjectID)
        }
        
        return userIDs, nil
    }

    // Hash passwords
    adminPassword, _ := utils.HashPassword("admin123")
    userPassword, _ := utils.HashPassword("user123")
    dosenPassword, _ := utils.HashPassword("dosen123")

    adminID := primitive.NewObjectID()
    user1ID := primitive.NewObjectID()
    user2ID := primitive.NewObjectID()
    dosenID := primitive.NewObjectID()

    users := []interface{}{
        bson.M{
            "_id":           adminID,
            "username":      "admin",
            "email":         "admin@university.ac.id",
            "password_hash": adminPassword,
            "role":          "admin",
            "created_at":    time.Now(),
        },
        bson.M{
            "_id":           user1ID,
            "username":      "johndoe",
            "email":         "john.doe@university.ac.id",
            "password_hash": userPassword,
            "role":          "user",
            "created_at":    time.Now(),
        },
        bson.M{
            "_id":           user2ID,
            "username":      "janesmith",
            "email":         "jane.smith@university.ac.id",
            "password_hash": userPassword,
            "role":          "user",
            "created_at":    time.Now(),
        },
        bson.M{
            "_id":           dosenID,
            "username":      "dosen",
            "email":         "dosen@university.ac.id",
            "password_hash": dosenPassword,
            "role":          "user",
            "created_at":    time.Now(),
        },
    }

    _, err = collection.InsertMany(ctx, users)
    if err != nil {
        log.Printf("❌ Failed to seed users: %v", err)
        return nil, err
    }

    log.Println("  ✓ Users seeded (4 users)")
    log.Println("    - admin / admin123")
    log.Println("    - johndoe / user123")
    log.Println("    - janesmith / user123")
    log.Println("    - dosen / dosen123")

    return map[string]primitive.ObjectID{
        "admin":     adminID,
        "johndoe":   user1ID,
        "janesmith": user2ID,
        "dosen":     dosenID,
    }, nil
}

// seedAlumni seeds alumni data
func seedAlumni(ctx context.Context, db *mongo.Database, userIDs map[string]primitive.ObjectID) ([]primitive.ObjectID, error) {
    collection := db.Collection(AlumniCollection)

    // Check if alumni already exist
    count, err := collection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    if count > 0 {
        log.Println("  ⏭️  Alumni already seeded, fetching existing data...")
        
        // Fetch existing alumni IDs
        cursor, err := collection.Find(ctx, bson.M{})
        if err != nil {
            return nil, err
        }
        defer cursor.Close(ctx)

        var alumniIDs []primitive.ObjectID
        for cursor.Next(ctx) {
            var alumni bson.M
            if err := cursor.Decode(&alumni); err != nil {
                return nil, err
            }
            alumniIDs = append(alumniIDs, alumni["_id"].(primitive.ObjectID))
        }
        
        return alumniIDs, nil
    }

    now := time.Now()
    alumni1ID := primitive.NewObjectID()
    alumni2ID := primitive.NewObjectID()
    alumni3ID := primitive.NewObjectID()
    alumni4ID := primitive.NewObjectID()
    alumni5ID := primitive.NewObjectID()
    alumni6ID := primitive.NewObjectID()
    alumni7ID := primitive.NewObjectID()
    alumni8ID := primitive.NewObjectID()

    alumni := []interface{}{
        bson.M{
            "_id":         alumni1ID,
            "nim":         "1234567001",
            "nama":        "John Doe",
            "jurusan":     "Teknik Informatika",
            "angkatan":    2018,
            "tahun_lulus": 2022,
            "email":       "john.doe@university.ac.id",
            "no_telepon":  "081234567001",
            "alamat":      "Jl. Merdeka No. 123, Jakarta",
            "user_id":     userIDs["johndoe"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni2ID,
            "nim":         "1234567002",
            "nama":        "Jane Smith",
            "jurusan":     "Sistem Informasi",
            "angkatan":    2019,
            "tahun_lulus": 2023,
            "email":       "jane.smith@university.ac.id",
            "no_telepon":  "081234567002",
            "alamat":      "Jl. Sudirman No. 456, Bandung",
            "user_id":     userIDs["janesmith"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni3ID,
            "nim":         "1234567003",
            "nama":        "Ahmad Rizki",
            "jurusan":     "Teknik Informatika",
            "angkatan":    2017,
            "tahun_lulus": 2021,
            "email":       "ahmad.rizki@university.ac.id",
            "no_telepon":  "081234567003",
            "alamat":      "Jl. Gatot Subroto No. 789, Surabaya",
            "user_id":     userIDs["dosen"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni4ID,
            "nim":         "1234567004",
            "nama":        "Siti Nurhaliza",
            "jurusan":     "Teknik Komputer",
            "angkatan":    2018,
            "tahun_lulus": 2022,
            "email":       "siti.nurhaliza@university.ac.id",
            "no_telepon":  "081234567004",
            "alamat":      "Jl. Ahmad Yani No. 321, Medan",
            "user_id":     userIDs["dosen"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni5ID,
            "nim":         "1234567005",
            "nama":        "Budi Santoso",
            "jurusan":     "Sistem Informasi",
            "angkatan":    2019,
            "tahun_lulus": 2023,
            "email":       "budi.santoso@university.ac.id",
            "no_telepon":  "081234567005",
            "alamat":      "Jl. Diponegoro No. 654, Semarang",
            "user_id":     userIDs["dosen"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni6ID,
            "nim":         "1234567006",
            "nama":        "Dewi Lestari",
            "jurusan":     "Teknik Informatika",
            "angkatan":    2020,
            "tahun_lulus": 2024,
            "email":       "dewi.lestari@university.ac.id",
            "no_telepon":  "081234567006",
            "alamat":      "Jl. Pahlawan No. 987, Yogyakarta",
            "user_id":     userIDs["dosen"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni7ID,
            "nim":         "1234567007",
            "nama":        "Rudi Hartono",
            "jurusan":     "Teknik Komputer",
            "angkatan":    2017,
            "tahun_lulus": 2021,
            "email":       "rudi.hartono@university.ac.id",
            "no_telepon":  "081234567007",
            "alamat":      "Jl. Veteran No. 147, Makassar",
            "user_id":     userIDs["dosen"],
            "created_at":  now,
            "updated_at":  now,
        },
        bson.M{
            "_id":         alumni8ID,
            "nim":         "1234567008",
            "nama":        "Linda Wijaya",
            "jurusan":     "Sistem Informasi",
            "angkatan":    2018,
            "tahun_lulus": 2022,
            "email":       "linda.wijaya@university.ac.id",
            "no_telepon":  "081234567008",
            "alamat":      "Jl. Proklamasi No. 258, Denpasar",
            "user_id":     userIDs["dosen"],
            "created_at":  now,
            "updated_at":  now,
        },
    }

    _, err = collection.InsertMany(ctx, alumni)
    if err != nil {
        log.Printf("❌ Failed to seed alumni: %v", err)
        return nil, err
    }

    log.Println("  ✓ Alumni seeded (8 alumni)")

    return []primitive.ObjectID{
        alumni1ID, alumni2ID, alumni3ID, alumni4ID,
        alumni5ID, alumni6ID, alumni7ID, alumni8ID,
    }, nil
}

// seedPekerjaan seeds pekerjaan data
func seedPekerjaan(ctx context.Context, db *mongo.Database, alumniIDs []primitive.ObjectID) error {
    collection := db.Collection(PekerjaanCollection)

    // Check if pekerjaan already exist
    count, err := collection.CountDocuments(ctx, bson.M{})
    if err != nil {
        return err
    }

    if count > 0 {
        log.Println("  ⏭️  Pekerjaan already seeded, skipping...")
        return nil
    }

    now := time.Now()
    oneYearAgo := now.AddDate(-1, 0, 0)
    twoYearsAgo := now.AddDate(-2, 0, 0)
    sixMonthsAgo := now.AddDate(0, -6, 0)
    threeMonthsAgo := now.AddDate(0, -3, 0)

    pekerjaan := []interface{}{
        // Alumni 1 (John Doe) - Teknik Informatika - 2 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[0],
            "nama_perusahaan":       "PT Teknologi Nusantara",
            "posisi_jabatan":        "Backend Developer",
            "bidang_industri":       "Technology",
            "lokasi_kerja":          "Jakarta",
            "gaji_range":            "8-12 juta",
            "tanggal_mulai_kerja":   twoYearsAgo,
            "tanggal_selesai_kerja": oneYearAgo,
            "status_pekerjaan":      "resign",
            "deskripsi_pekerjaan":   "Mengembangkan API menggunakan Go dan Node.js",
            "created_at":            now,
            "updated_at":            now,
        },
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[0],
            "nama_perusahaan":       "PT Digital Indonesia",
            "posisi_jabatan":        "Senior Backend Engineer",
            "bidang_industri":       "Financial Technology",
            "lokasi_kerja":          "Jakarta",
            "gaji_range":            "15-20 juta",
            "tanggal_mulai_kerja":   oneYearAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Lead backend development untuk aplikasi fintech",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 2 (Jane Smith) - Sistem Informasi - 1 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[1],
            "nama_perusahaan":       "PT Solusi Bisnis Digital",
            "posisi_jabatan":        "Business Analyst",
            "bidang_industri":       "Consulting",
            "lokasi_kerja":          "Bandung",
            "gaji_range":            "10-15 juta",
            "tanggal_mulai_kerja":   sixMonthsAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Melakukan analisis bisnis dan requirement gathering",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 3 (Ahmad Rizki) - Teknik Informatika - 3 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[2],
            "nama_perusahaan":       "PT Maju Jaya",
            "posisi_jabatan":        "Junior Programmer",
            "bidang_industri":       "Manufacturing",
            "lokasi_kerja":          "Surabaya",
            "gaji_range":            "5-7 juta",
            "tanggal_mulai_kerja":   twoYearsAgo.AddDate(-1, 0, 0),
            "tanggal_selesai_kerja": twoYearsAgo,
            "status_pekerjaan":      "resign",
            "deskripsi_pekerjaan":   "Maintenance sistem ERP perusahaan",
            "created_at":            now,
            "updated_at":            now,
        },
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[2],
            "nama_perusahaan":       "PT Startup Indonesia",
            "posisi_jabatan":        "Full Stack Developer",
            "bidang_industri":       "E-commerce",
            "lokasi_kerja":          "Jakarta",
            "gaji_range":            "10-15 juta",
            "tanggal_mulai_kerja":   twoYearsAgo,
            "tanggal_selesai_kerja": sixMonthsAgo,
            "status_pekerjaan":      "kontrak_habis",
            "deskripsi_pekerjaan":   "Develop marketplace platform",
            "created_at":            now,
            "updated_at":            now,
        },
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[2],
            "nama_perusahaan":       "PT Global Tech Solutions",
            "posisi_jabatan":        "Tech Lead",
            "bidang_industri":       "Information Technology",
            "lokasi_kerja":          "Jakarta",
            "gaji_range":            "18-25 juta",
            "tanggal_mulai_kerja":   sixMonthsAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Memimpin tim development dan arsitektur sistem",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 4 (Siti Nurhaliza) - Teknik Komputer - 2 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[3],
            "nama_perusahaan":       "PT Telekomunikasi Indonesia",
            "posisi_jabatan":        "Network Engineer",
            "bidang_industri":       "Telecommunications",
            "lokasi_kerja":          "Medan",
            "gaji_range":            "8-12 juta",
            "tanggal_mulai_kerja":   oneYearAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Maintenance dan monitoring infrastruktur jaringan",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 5 (Budi Santoso) - Sistem Informasi - 1 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[4],
            "nama_perusahaan":       "PT Bank Digital Indonesia",
            "posisi_jabatan":        "Data Analyst",
            "bidang_industri":       "Banking",
            "lokasi_kerja":          "Semarang",
            "gaji_range":            "9-13 juta",
            "tanggal_mulai_kerja":   threeMonthsAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Analisis data transaksi dan customer behavior",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 6 (Dewi Lestari) - Teknik Informatika - 1 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[5],
            "nama_perusahaan":       "PT Inovasi Digital",
            "posisi_jabatan":        "Frontend Developer",
            "bidang_industri":       "Technology",
            "lokasi_kerja":          "Yogyakarta",
            "gaji_range":            "7-10 juta",
            "tanggal_mulai_kerja":   threeMonthsAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Develop UI/UX menggunakan React dan Vue",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 7 (Rudi Hartono) - Teknik Komputer - 2 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[6],
            "nama_perusahaan":       "PT Cybersecurity Indonesia",
            "posisi_jabatan":        "Security Analyst",
            "bidang_industri":       "Cybersecurity",
            "lokasi_kerja":          "Makassar",
            "gaji_range":            "12-17 juta",
            "tanggal_mulai_kerja":   oneYearAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Monitor dan analisis keamanan sistem",
            "created_at":            now,
            "updated_at":            now,
        },

        // Alumni 8 (Linda Wijaya) - Sistem Informasi - 2 pekerjaan
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[7],
            "nama_perusahaan":       "PT Media Online",
            "posisi_jabatan":        "Product Manager",
            "bidang_industri":       "Media",
            "lokasi_kerja":          "Denpasar",
            "gaji_range":            "11-16 juta",
            "tanggal_mulai_kerja":   sixMonthsAgo,
            "tanggal_selesai_kerja": nil,
            "status_pekerjaan":      "aktif",
            "deskripsi_pekerjaan":   "Manage product development lifecycle",
            "created_at":            now,
            "updated_at":            now,
        },

        // Soft deleted example
        bson.M{
            "_id":                   primitive.NewObjectID(),
            "alumni_id":             alumniIDs[0],
            "nama_perusahaan":       "PT Legacy Systems (DELETED)",
            "posisi_jabatan":        "Junior Developer",
            "bidang_industri":       "Technology",
            "lokasi_kerja":          "Jakarta",
            "gaji_range":            "5-7 juta",
            "tanggal_mulai_kerja":   twoYearsAgo.AddDate(-2, 0, 0),
            "tanggal_selesai_kerja": twoYearsAgo.AddDate(-1, 0, 0),
            "status_pekerjaan":      "resign",
            "deskripsi_pekerjaan":   "First job after graduation",
            "is_delete":             now.AddDate(0, -1, 0), // Deleted 1 month ago
            "created_at":            now,
            "updated_at":            now,
        },
    }

    _, err = collection.InsertMany(ctx, pekerjaan)
    if err != nil {
        log.Printf("❌ Failed to seed pekerjaan: %v", err)
        return err
    }

    log.Println("  ✓ Pekerjaan seeded (13 records, including 1 soft-deleted)")

    return nil
}

// SeedSummary prints seeding summary
func SeedSummary(db *mongo.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    log.Println("\n📊 Database Summary:")
    log.Println("==========================================")

    // Count users
    userCount, _ := db.Collection(UsersCollection).CountDocuments(ctx, bson.M{})
    log.Printf("  Users: %d", userCount)

    // Count alumni
    alumniCount, _ := db.Collection(AlumniCollection).CountDocuments(ctx, bson.M{})
    log.Printf("  Alumni: %d", alumniCount)

    // Count alumni by jurusan
    pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$jurusan"},
            {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
        }}},
        {{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
    }

    cursor, _ := db.Collection(AlumniCollection).Aggregate(ctx, pipeline)
    defer cursor.Close(ctx)

    log.Println("  Alumni by Jurusan:")
    for cursor.Next(ctx) {
        var result bson.M
        cursor.Decode(&result)
        log.Printf("    - %s: %v", result["_id"], result["count"])
    }

    // Count pekerjaan
    pekerjaanCount, _ := db.Collection(PekerjaanCollection).CountDocuments(ctx, bson.M{})
    pekerjaanActive, _ := db.Collection(PekerjaanCollection).CountDocuments(ctx, bson.M{"is_delete": bson.M{"$exists": false}})
    pekerjaanDeleted, _ := db.Collection(PekerjaanCollection).CountDocuments(ctx, bson.M{"is_delete": bson.M{"$exists": true}})

    log.Printf("  Pekerjaan Total: %d", pekerjaanCount)
    log.Printf("    - Active: %d", pekerjaanActive)
    log.Printf("    - Soft Deleted: %d", pekerjaanDeleted)

    log.Println("==========================================")

    return nil
}