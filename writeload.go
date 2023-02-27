package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	"math/rand"
	"time"
)

func writeToVault(client *api.Client, nspath string, count int) int {
	errs := 0
	for i := 0; i < count; i++ {
		fullPath := fmt.Sprintf("%s/secret-%d", nspath, i)

		// Define the secret data and metadata
		secretData := map[string]interface{}{
			"Key":             randSeq(128),
			"UsagePlanName":   randSeq(68),
			"Username":        randSeq(32),
			"Email":           randSeq(32),
			"ApiProductName":  randSeq(68),
			"EnvironmentName": randSeq(68),
		}

		_, err := client.KVv2("secret").Put(context.TODO(), fullPath, secretData)
		if err != nil {
			fmt.Println(err.Error())
			errs++
			i--
			time.Sleep(20 * time.Millisecond)
			continue
		}

		//fmt.Printf("Secret '%s' created successfully: %d\n", fullPath, count)
		//print(".")
	}
	println("")
	return errs
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
