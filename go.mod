module github.com/RileySun/sassy

go 1.22.4

replace api => ./api

replace auth => ./auth

replace admin => ./admin

require (
	admin v0.0.0-00010101000000-000000000000
	api v0.0.0-00010101000000-000000000000
	auth v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/phpdave11/gofpdi v1.0.14-0.20211212211723-1f10f9844311 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/signintech/gopdf v0.29.2 // indirect
	github.com/vicanso/go-charts/v2 v2.6.10 // indirect
	github.com/wcharczuk/go-chart/v2 v2.1.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5 // indirect
)
