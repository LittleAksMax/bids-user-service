module github.com/LittleAksMax/bids-user-service

go 1.25

require (
	github.com/LittleAksMax/amazon-ads-api-sdk-go v0.0.0-20260321104555-703ca7f29191
	github.com/LittleAksMax/bids-util v0.0.6-0.20260218173913-01f4995c0cfd
	github.com/go-chi/chi/v5 v5.2.5
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/joho/godotenv v1.5.1
	github.com/pressly/goose/v3 v3.26.0
	github.com/redis/go-redis/v9 v9.18.0
)

replace github.com/LittleAksMax/amazon-ads-api-sdk-go => ../amazon-ads-api-go-sdk
replace github.com/LittleAksMax/bids-util => ../bids-util

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-chi/cors v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)
