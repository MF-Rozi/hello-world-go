module dev.mfr/main

go 1.24.5

replace dev.mfr/db => ../db

require (
	dev.mfr/db v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
)
