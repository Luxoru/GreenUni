package redis

import (
	"backend/internal/db"
	"backend/internal/utils/concurrency"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	var redisContainer = Container{}
	var config = Configurations{
		authentications: db.AuthenticationConfigurations{
			Host:     "localhost",
			Port:     6379,
			Username: "",
			Password: "",
		},
		DatabaseIndex: 0,
	}
	err := redisContainer.Connect(config)
	if err != nil {
		panic(err)
		return
	}

	repo := Repository{
		Database: redisContainer,
	}

	now := time.Now()

	for i := 0; i < 200; i++ {
		key := GenerateRandomString(5)
		value := GenerateRandomString(5)
		repo.Insert(key, value, &concurrency.Callback[string]{
			Success: func(value string) {
				getValue(key, &repo, now)

			},
			Error: func(errorType error) {
				fmt.Printf("Error occured %s\n", errorType)
			},
		})
	}

	time.Sleep(1000 * time.Millisecond)

}

func getValue(Key string, repository *Repository, t time.Time) {
	repository.Fetch(Key, &concurrency.Callback[string]{
		Error: func(errorType error) {
			panic(errorType)
		},
		Success: func(_ string) {

		},
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano())) //

	// Create a slice of random characters
	result := make([]byte, length)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))] // Select a random character from letterBytes
	}
	return string(result)
}
