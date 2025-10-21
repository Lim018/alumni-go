package model

import "time"

// ToAlumniResponse - Convert Alumni to AlumniResponse
func (a *Alumni) ToAlumniResponse() AlumniResponse {
    return AlumniResponse{
        ID:         a.ID.Hex(),
        NIM:        a.NIM,
        Nama:       a.Nama,
        Jurusan:    a.Jurusan,
        Angkatan:   a.Angkatan,
        TahunLulus: a.TahunLulus,
        Email:      a.Email,
        NoTelepon:  a.NoTelepon,
        Alamat:     a.Alamat,
        CreatedAt:  a.CreatedAt,
        UpdatedAt:  a.UpdatedAt,
        UserID:     a.UserID.Hex(),
    }
}

// ToUserResponse - Convert User to UserResponse
func (u *User) ToUserResponse() UserResponse {
    return UserResponse{
        ID:       u.ID.Hex(),
        Username: u.Username,
        Email:    u.Email,
        Role:     u.Role,
    }
}

// ToPekerjaanResponse converts Pekerjaan to PekerjaanResponse
func (p *Pekerjaan) ToPekerjaanResponse() PekerjaanResponse {
    return PekerjaanResponse{
        ID:                  p.ID.Hex(),
        AlumniID:            p.AlumniID.Hex(),
        NamaPerusahaan:      p.NamaPerusahaan,
        PosisiJabatan:       p.PosisiJabatan,
        BidangIndustri:      p.BidangIndustri,
        LokasiKerja:         p.LokasiKerja,
        GajiRange:           p.GajiRange,
        TanggalMulaiKerja:   p.TanggalMulaiKerja,
        TanggalSelesaiKerja: p.TanggalSelesaiKerja,
        StatusPekerjaan:     p.StatusPekerjaan,
        DeskripsiPekerjaan:  p.DeskripsiPekerjaan,
        CreatedAt:           p.CreatedAt,
        UpdatedAt:           p.UpdatedAt,
    }
}

// ToPekerjaanTrashResponse converts Pekerjaan to PekerjaanTrashResponse
func (p *Pekerjaan) ToPekerjaanTrashResponse() PekerjaanTrashResponse {
    deletedAt := time.Time{}
    if p.IsDelete != nil {
        deletedAt = *p.IsDelete
    }
    
    return PekerjaanTrashResponse{
        ID:                  p.ID.Hex(),
        AlumniID:            p.AlumniID.Hex(),
        NamaPerusahaan:      p.NamaPerusahaan,
        PosisiJabatan:       p.PosisiJabatan,
        BidangIndustri:      p.BidangIndustri,
        LokasiKerja:         p.LokasiKerja,
        GajiRange:           p.GajiRange,
        TanggalMulaiKerja:   p.TanggalMulaiKerja,
        TanggalSelesaiKerja: p.TanggalSelesaiKerja,
        StatusPekerjaan:     p.StatusPekerjaan,
        DeskripsiPekerjaan:  p.DeskripsiPekerjaan,
        CreatedAt:           p.CreatedAt,
        UpdatedAt:           p.UpdatedAt,
        DeletedAt:           deletedAt,
    }
}