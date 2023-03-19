package main

import (
	"fmt"
	"log"

	"github.com/kneadCODE/crazycat/apps/golib/app"
)

func main() {
	fmt.Println("Hello from cazycat golib")

	ctx, finish, err := app.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer finish()

	var _ = ctx
}
