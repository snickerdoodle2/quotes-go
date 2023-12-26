package main

import (
	"fmt"
	"os"
)

func main() {
    db_url := os.Getenv("DATABASE_URL");
    fmt.Printf("DB_URL: %s", db_url);
}
