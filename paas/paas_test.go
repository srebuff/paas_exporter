package paas

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

var db *sql.DB
var dbClient *mongo.Client

type MysqlConnect struct {
	Host string
	Port string
}

type MongoDBConnect struct {
	Host string
	Port string
}

var mysqlConnect MysqlConnect
var mongoConnect MongoDBConnect

func initMysql(pool *dockertest.Pool) *dockertest.Resource {
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		mysqlConnect.Host = "localhost"
		mysqlConnect.Port = resource.GetPort("3306/tcp")
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	return resource
}

func initMongoDB(pool *dockertest.Pool) *dockertest.Resource {
	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "4.4",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		mongoConnect.Host = "localhost"
		mongoConnect.Port = resource.GetPort("27017/tcp")
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return resource
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	var resourceMysql *dockertest.Resource
	MYSQL_ENABLE := os.Getenv("MYSQL_ENABLE")
	if MYSQL_ENABLE == "true" {
		resourceMysql = initMysql(pool)
	}
	var resourceMongoDB *dockertest.Resource
	MONGODB_ENABLE := os.Getenv("MONGODB_ENABLE")
	if MONGODB_ENABLE == "true" {
		resourceMongoDB = initMongoDB(pool)
	}
	code := m.Run()
	if MYSQL_ENABLE == "true" {
		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resourceMysql); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}

	if MONGODB_ENABLE == "true" {
		// disconnect mongodb client
		if err = dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}

		// When you're done, kill and remove the container
		if err = pool.Purge(resourceMongoDB); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}

	}

	os.Exit(code)
}

func TestConnectMysql(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				dsn: fmt.Sprintf("root:secret@tcp(%s:%s)/mysql", mysqlConnect.Host, mysqlConnect.Port),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConnectMysql(tt.args.dsn); (err != nil) != tt.wantErr {
				t.Errorf("ConnectMysql() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnectMongoDB(t *testing.T) {
	type args struct {
		dsn string
	}
	var mongodsn string
	mongodsn = fmt.Sprintf("mongodb://localhost:%s", mongoConnect.Port)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				dsn: mongodsn,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConnectMongoDB(tt.args.dsn); (err != nil) != tt.wantErr {
				t.Errorf("ConnectMongoDB=%s, error = %v, wantErr %v", mongodsn, err, tt.wantErr)
			}
		})
	}
}
