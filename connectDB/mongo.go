package connectDB

import (
	"os"

	mongo "github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/bhupeshbhatia/go-agg-monitoring/model"
	"github.com/pkg/errors"
)

type Db struct {
	Collection *mongo.Collection
}

func ConfirmDbExists() (*Db, error) {
	hosts := os.Getenv("MONGO_HOSTS")
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("MONGO_DATABASE")
	collection := os.Getenv("MONGO_COLLECTION")

	clientConfig := mongo.ClientConfig{
		Hosts:               []string{hosts},
		Username:            username,
		Password:            password,
		TimeoutMilliseconds: 3000,
	}

	// ====> MongoDB Client
	client, err := mongo.NewClient(clientConfig)
	if err != nil {
		err = errors.Wrap(err, "Mongo client not available")
		return nil, err
	}

	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: 5000,
	}

	// Index Configuration
	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name:        "item_id",
					IsDescOrder: false,
				},
			},
			IsUnique: true,
			Name:     "item_id_index",
		},
	}

	// ====> Create New Collection
	colConfig := &mongo.Collection{
		Connection:   conn,
		Name:         collection,
		Database:     database,
		SchemaStruct: &model.Metric{},
		Indexes:      indexConfigs,
	}
	c, err := mongo.EnsureCollection(colConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}

	return &Db{
		Collection: c,
	}, nil

}
