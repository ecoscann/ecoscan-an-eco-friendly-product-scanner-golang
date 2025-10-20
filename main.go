package main

import (
	"ecoscan.com/cmd"
	_ "github.com/lib/pq"
)

/* DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=1212
DB_ENABLE_SSL_MODE=false */

func main() {

	cmd.Serve()
}
