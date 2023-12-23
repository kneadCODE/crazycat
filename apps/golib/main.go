package main

import (
	"log"
	"os"

	"github.com/kneadCODE/crazycat/apps/golib/app2"
)

func main() {
	log.New(os.Stdout, "", log.LstdFlags).Println("Hello from cazycat golib")

	ctx, finish, err := app2.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer finish()

	var _ = ctx
}
