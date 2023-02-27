package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"time"
)

func main() {
	config := api.DefaultConfig()
	config.Address = "http://127.0.0.1:8200"
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	errs := 0

	client.SetToken("password")
	max := 1000000
	steps := max / 10000

	for i := 0; i < steps; i++ {
		path := fmt.Sprintf("ns%d", i)
		errs += writeToVault(client, path, max/steps)

		time.Sleep(200 * time.Millisecond)

		milisecondsPerRequest, errsOnRead := measureAveragePerformance(client, max, steps, i)
		fmt.Printf("Average response time: %f ms, errors: %d\n", milisecondsPerRequest, errsOnRead)
	}

	fmt.Printf("Total errors adding %d keys was: %d\n", max, errs)
}
