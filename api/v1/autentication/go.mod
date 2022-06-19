module go-save-water/api/v1/autentication

go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/justinas/alice v1.2.0
	go-save-water/pkg/common v0.0.0-00010101000000-000000000000
	go-save-water/pkg/log v0.0.0-00010101000000-000000000000
	go-save-water/pkg/middleware v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
)

require github.com/joho/godotenv v1.4.0 // indirect

replace (
	go-save-water/pkg/common => ../../../pkg/common
	go-save-water/pkg/log => ../../../pkg/log
	go-save-water/pkg/middleware => ../../../pkg/middleware
)
