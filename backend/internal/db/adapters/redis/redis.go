package redis

import (
	"backend/internal/db"
	"backend/internal/utils/concurrency"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// Configurations Make sure Configurations implements db.DatabaseConfigurations
type Configurations struct {
	authentications db.AuthenticationConfigurations
	DatabaseIndex   int
}

// GetAuthenticationConfigurations This method ensures Configurations implements db.DatabaseConfigurations
func (config Configurations) GetAuthenticationConfigurations() db.AuthenticationConfigurations {
	return config.authentications
}

// Container used to store client andd thread pool implements Database
type Container struct {
	redis *redis.Client
	pool  *concurrency.ThreadPool
}

// Name Returns the name of the database type (Redis)
func (redisDatabase *Container) Name() string {
	return "Redis"
}

var emptyConfig = db.AuthenticationConfigurations{}

// Connect Initializes the connection to Redis and starts the thread pool
func (redisDatabase *Container) Connect(config Configurations) error {

	if config.authentications == emptyConfig {
		return fmt.Errorf("authentication must be provided")
	}

	host := config.authentications.Host
	port := config.authentications.Port

	if host == "" || port < -1 {
		return fmt.Errorf("invalid connection string. Host: %s, Port: %d", host, port)
	}

	connectionString := host + ":" + strconv.Itoa(port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: config.authentications.Password,
		DB:       config.DatabaseIndex,
	})

	if rdb == nil {
		return fmt.Errorf("instantiation of client returned null")
	}

	//ping to actually check we connected
	ping := rdb.Ping(context.Background())

	_, err := ping.Result()
	if err != nil {
		return fmt.Errorf("ping returned error")
	}

	redisDatabase.redis = rdb
	redisDatabase.pool = concurrency.NewThreadPool(500, 100)
	redisDatabase.pool.Start()

	return nil
}

// Close Closes the Redis client connection
func (redisDatabase *Container) Close() error {
	return redisDatabase.redis.Close()
}

// GetThreadPool Returns the thread pool for concurrency management
func (redisDatabase *Container) GetThreadPool() *concurrency.ThreadPool {
	return redisDatabase.pool
}

type Repository struct {
	Database Container
}

// Insert Adds a new key-value pair to Redis asynchronously
func (repo *Repository) Insert(Key string, Value string, callback *concurrency.Callback[string]) {
	repo.Database.pool.Submit(func() {
		repo.insertIntoDb(Key, Value, callback)
	})
}

// Fetch Retrieves a value by key from Redis asynchronously
func (repo *Repository) Fetch(Key string, callback *concurrency.Callback[string]) {
	repo.Database.pool.Submit(func() {
		repo.fetchFromDb(Key, callback)
	})
}

// fetchFromDb Helper method for fetching data from Redis
func (repo *Repository) fetchFromDb(Key string, callback *concurrency.Callback[string]) {
	client, err := setupClient(repo, callback)
	if err != nil {
		return
	}

	ctx := context.Background()

	result, err := client.Get(ctx, Key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			callback.Error(fmt.Errorf("key %s doesn't exist", Key))
			return
		}
		callback.Error(fmt.Errorf("error occured whilst fetching key (%s)\nError: %s", Key, err))
		return
	}

	if callback.Success == nil {
		return
	}

	callback.Success(result)
}

// insertIntoDb Helper method for inserting data into Redis
func (repo *Repository) insertIntoDb(Key string, Value string, callback *concurrency.Callback[string]) {

	client, err := setupClient(repo, callback)
	if err != nil {
		return
	}

	ctx := context.Background()

	err = client.Set(ctx, Key, Value, 0).Err()

	if err != nil {
		callback.Error(err)
		return
	}

	callback.Success(Value)
}

// setupClient Sets up the Redis client for the repository
func setupClient[T any](repo *Repository, callback *concurrency.Callback[T]) (*redis.Client, error) {
	database := repo.Database

	if database == (Container{}) {
		err := fmt.Errorf("database is nill")
		callback.Error(err)
		return nil, err
	}

	client := database.redis

	if client == nil {
		err := fmt.Errorf("client is nill")
		callback.Error(err)
		return nil, err
	}

	return client, nil
}
