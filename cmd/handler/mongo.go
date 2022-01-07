package handler

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DbHandler interface {
	AddPoi(poi PoiDbEntry) (id string, err error)
	GetPoi(id string) (poi PoiDbEntry, err error)
	UpdatePoi(id string, poi PoiDbEntry) (err error)
	DeletePoi(id string) (err error)
	SearchByRadius(location Location, distanceInMeter uint64) (result PoiDbEntries, err error)
	GetAllPois() (result PoiDbEntries, err error)
}

func NewDbHandler(url string) (DbHandler, error) {
	log.Info().Str("url", url).Msg("db connection")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Connecting to mongodb failed")
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Ping to mongodb")
		return nil, err
	}

	log.Info().Str("url", url).Msg("Connecting to mongodb done")

	handler := &dbHandler{
		dbClient:   client,
		dbName:     "poiDb",
		collection: "poi",
	}

	handler.createIndex()

	return handler, nil
}

type dbHandler struct {
	dbClient   *mongo.Client
	dbName     string
	collection string
}

type PoiDbEntry struct {
	Id       string   `json:"id" bson:"_id"`
	Name     string   `json:"name" bson:"name"`
	Location Location `json:"location" bson:"location"`
}

type PoiDbEntries []PoiDbEntry

// We need this type so we can store it in our mongodb db and do geospatial queries
// https://docs.mongodb.com/manual/geospatial-queries/
type Location struct {
	GeoJSONType string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

func NewLocation(lat, long float64) Location {
	return Location{
		"Point",
		[]float64{long, lat},
	}
}

func (c *dbHandler) getMongoDbCollection() (collection *mongo.Collection) {
	collection = c.dbClient.Database(c.dbName).Collection(c.collection)
	return
}

func (c *dbHandler) AddPoi(poi PoiDbEntry) (id string, err error) {
	insertResult, err := c.getMongoDbCollection().InsertOne(context.TODO(), poi)
	if err != nil {
		log.Printf("Could not insert new Point. Id")
		return "", err
	}
	id = poi.Id
	log.Info().Str("poi id", poi.Id).Msg("Inserted new Point")
	log.Info().Str("id", fmt.Sprint(insertResult.InsertedID)).Msg("created unique id")

	return
}

func (c *dbHandler) GetPoi(id string) (poi PoiDbEntry, err error) {
	filter := bson.M{"_id": bson.M{"$eq": id}}
	if err = c.getMongoDbCollection().FindOne(context.TODO(), filter).Decode(&poi); err != nil {
		//fmt.Println(err)
		return
	}
	return
}

func (c *dbHandler) GetAllPois() (result PoiDbEntries, err error) {
	cur, err := c.getMongoDbCollection().Find(context.TODO(), bson.M{})
	if err != nil {
		log.Warn().Err(err).Msg("GetAllPois failed")
		return
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem PoiDbEntry
		err := cur.Decode(&elem)
		if err != nil {
			log.Warn().Err(err).Msg("GetAllPois failed")
		}

		result = append(result, elem)
	}

	return
}

func (c *dbHandler) UpdatePoi(id string, poi PoiDbEntry) (err error) {
	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{
		"$set": bson.M{"name": poi.Name, "location": poi.Location},
	}
	_, err = c.getMongoDbCollection().UpdateOne(
		context.Background(),
		filter,
		update,
	)

	return
}

func (c *dbHandler) DeletePoi(id string) (err error) {
	filter := bson.M{"_id": bson.M{"$eq": id}}
	_, err = c.getMongoDbCollection().DeleteOne(context.TODO(), filter)
	return
}

func (c *dbHandler) SearchByRadius(location Location, distanceInMeter uint64) (result PoiDbEntries, err error) {
	// connect to mongo
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to db")
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// query the database
	col := session.DB(c.dbName).C(c.collection)
	err = col.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{location.Coordinates[0], location.Coordinates[1]},
				},
				"$maxDistance": distanceInMeter,
			},
		},
	}).All(&result)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (c *dbHandler) createIndex() (err error) {
	pointIndexModel := mongo.IndexModel{
		Keys: bsonx.MDoc{"location": bsonx.String("2dsphere")},
	}

	_, err = c.getMongoDbCollection().Indexes().CreateOne(context.TODO(), pointIndexModel)
	return
}
