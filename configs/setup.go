package configs

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(ENV_MONGO_URI_LOCAL()))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return client
}

//Client instance
var DB *mongo.Client = ConnectDB()

//getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(ENV_MONGO_DB()).Collection(collectionName)
	return collection
}

func InitRedis() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//Initializing redis
	client := redis.NewClient(&redis.Options{
		Addr: ENV_REDIS_DSN(),
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return client
}

var RDB *redis.Client = InitRedis()

func InitMinio() *minio.Client {
	endpoint := ENV_MINIO_ENDPOINT()
	accessKeyID := ENV_MINIO_ACCESS_KEY()
	secretAccessKey := ENV_MINIO_SECRET_KEY()
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	return minioClient
}

var MinioClient *minio.Client = InitMinio()
