module go-save-water/cmd/main

go 1.18

replace (
	go-save-water/pkg/common => ../../pkg/common
	go-save-water/pkg/log => ../../pkg/log
	go-save-water/pkg/validator => ../../pkg/validator
)

require (
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/gorilla/mux v1.8.0
	go-save-water/pkg/common v0.0.0-00010101000000-000000000000
	go-save-water/pkg/log v0.0.0-00010101000000-000000000000
	go-save-water/pkg/validator v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	gopkg.in/mail.v2 v2.3.1
)

require (
	github.com/joho/godotenv v1.4.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
