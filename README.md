# Go Fiber + GORM TODO — JWT + RBAC + OpenAPI + Profile

Tambahan fitur:
- **Profile**: `GET /api/v1/me`, `PUT /api/v1/me` (ubah name), `PATCH /api/v1/me/password` (ganti password), `POST /api/v1/me/avatar` (upload avatar).
- **Static uploads**: file avatar dapat diakses di `/uploads/<filename>` (serve dari `UPLOAD_DIR`).

## Quick Start
### Docker
```bash
cp .env.example .env
docker compose up --build
# Docs: http://localhost:8080/docs
```
### Lokal
```bash
cp .env.example .env
go mod tidy
go run ./cmd/server
```

## Env
- `UPLOAD_DIR` (default `./uploads`), di Docker: `/data/uploads` (otomatis dimount volume).

## Alur Avatar
1. Kirim `POST /api/v1/me/avatar` (multipart) field `avatar` (png|jpg|jpeg|webp, max 2MB).
2. Respon berisi `avatar_url` relatif: `/uploads/...`
3. Akses langsung via browser: `http://localhost:8080/uploads/...`

## Catatan
- Demi keamanan, endpoint update profile hanya mengizinkan `name`. (Email/role tidak bisa diubah via endpoint ini.)
- Password minimal 6 karakter, wajib mengirim `old_password` yang valid.
