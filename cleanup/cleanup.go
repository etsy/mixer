package main

import (
	"fmt"
	. "github.com/etsy/mixer/db"
)

func main() {
	fmt.Printf("running cleanup....\n")
	CleanupAlumni()
}
