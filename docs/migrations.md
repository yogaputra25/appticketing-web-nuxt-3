# Migrations

Menggunakan [golang-migrate](https://github.com/golang-migrate/migrate).

## Install CLI

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/

# Windows (scoop)
scoop install golang-migrate
```

## Pakai

```bash
# Apply semua migration
make migrate-up

# Rollback 1 step
make migrate-down

# Atau langsung:
migrate -path apps/api/migrations -database "$DATABASE_URL" up
```

## Convention

- File naming: `<version>_<description>.up.sql` & `.down.sql`
- Version: 6 digit zero-padded (e.g. `000001`)
- Setiap `up` harus punya `down`配套
