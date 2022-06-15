module go-save-water/cmd/main

go 1.18

replace (
	go-save-water/pkg/common => ../../pkg/common
	go-save-water/pkg/log => ../../pkg/log
)

require (
	github.com/gorilla/mux v1.8.0
	go-save-water/pkg/common v0.0.0-00010101000000-000000000000
	go-save-water/pkg/log v0.0.0-00010101000000-000000000000
)

require github.com/joho/godotenv v1.4.0 // indirect
