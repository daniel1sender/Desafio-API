#%/bin/bash
export DB_URL="postgres://postgres:1234@localhost:5432/desafio"
export API_PORT=":3000"
export TOKEN_SECRET="PU3lcBUKmE"
export EXP_TIME="5m"
go run main.go
