package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/DarkAnHell/FastPhish/api"
	"github.com/DarkAnHell/FastPhish/pkg/db/redis"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("missing JSON config file path")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open config file %s: %v", os.Args[1], err)
	}

	var db redis.Redis

	err = db.Load(f)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("Enter key: ")
		key, _ := reader.ReadString('\n')
		fmt.Print("Enter score: ")
		score, _ := reader.ReadString('\n')
		score = strings.TrimSuffix(score, "\n")

		scoreInt, _ := strconv.Atoi(score)

		d := api.DomainScore{Name: key, Score: uint32(scoreInt)}

		fmt.Printf("%v\n", db.Store(d))
	}

}
