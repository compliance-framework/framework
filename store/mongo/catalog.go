package mongo

import (
	"context"
	"log"

	"github.com/compliance-framework/configuration-service/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CatalogStoreMongo struct {
	collection *mongo.Collection
}

func NewCatalogStore() *CatalogStoreMongo {
	return &CatalogStoreMongo{
		collection: Collection("catalog"),
	}
}

func (c *CatalogStoreMongo) CreateCatalog(catalog *domain.Catalog) (interface{}, error) {
	result, err := c.collection.InsertOne(context.TODO(), catalog)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (store *CatalogStoreMongo) GetCatalog(id string) (*domain.Catalog, error) {
	var catalog domain.Catalog
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	err = store.collection.FindOne(context.Background(), filter).Decode(&catalog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &catalog, nil
}

func (store *CatalogStoreMongo) UpdateCatalog(id string, catalog *domain.Catalog) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": catalog}
	_, err = store.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (store *CatalogStoreMongo) DeleteCatalog(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = store.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (store *CatalogStoreMongo) CreateControl(catalogId string, control *domain.Control) (interface{}, error) {
	log.Println("CreateControl called with catalogId:", catalogId)

	catalogObjID, err := primitive.ObjectIDFromHex(catalogId)
	if err != nil {
		log.Println("Error converting catalogId to ObjectID:", err)
		return nil, err
	}

	control.Uuid = uuid.New() // Assign a new UUID to the control

	filter := bson.M{"_id": catalogObjID}
	update := bson.M{"$push": bson.M{"controls": *control}}
	result, err := store.collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Println("Error updating collection:", err)
		return nil, err
	}

	log.Println("CreateControl successful, updated count:", result.ModifiedCount)
	return control.Uuid, nil // Return the UUID of the control
}
