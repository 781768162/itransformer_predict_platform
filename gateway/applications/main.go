package main

import (
	"gateway/applications/router"
)

func main() {
	r := router.MustNewRouter()

	r.Run("localhost:8080")
}
