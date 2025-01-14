module github.com/RileySun/sassy

go 1.22.4

replace api => ./api

require (
	api v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	golang.org/x/crypto v0.29.0 // indirect
)
