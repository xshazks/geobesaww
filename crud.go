package geobesaww

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ------------------------------------------------------------------ Set Connection ------------------------------------------------------------------

func SetConnection(mongoenv, dbname string) *mongo.Database {
	var DBmongoinfo = DBInfo{
		DBString: os.Getenv(mongoenv),
		DBName:   dbname,
	}
	return MongoConnect(DBmongoinfo)
}

func SetConnection2dsphere(mongoenv, dbname, collname string) *mongo.Database {
	var DBmongoinfo = DBInfo{
		DBString:       os.Getenv(mongoenv),
		DBName:         dbname,
		CollectionName: collname,
	}
	return Create2dsphere(DBmongoinfo)
}

// ----------------------------------------------------------------------- User -----------------------------------------------------------------------

// Create

func InsertUser(mongoenv *mongo.Database, collname string, datauser User) interface{} {
	return InsertOneDoc(mongoenv, collname, datauser)
}

// Read

func GetAllUser(mconn *mongo.Database, collname string) []User {
	return GetAllDoc[[]User](mconn, collname)
}

func FindUser(mconn *mongo.Database, collname string, userdata User) User {
	filter := bson.M{"username": userdata.Username}
	return GetOneDoc[User](mconn, collname, filter)
}

func IsPasswordValid(mconn *mongo.Database, collname string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	res := GetOneDoc[User](mconn, collname, filter)
	hashChecker := CheckPasswordHash(userdata.Password, res.Password)
	return hashChecker
}

func UsernameExists(mconn *mongo.Database, collname string, userdata User) bool {
	filter := bson.M{"username": userdata.Username}
	return DocExists[User](mconn, collname, filter, userdata)
}

// Update

func UpdateUser(mconn *mongo.Database, collname string, datauser User) interface{} {
	filter := bson.M{"username": datauser.Username}
	return ReplaceOneDoc(mconn, collname, filter, datauser)
}

// Delete

func DeleteUser(mconn *mongo.Database, collname string, userdata User) interface{} {
	filter := bson.M{"username": userdata.Username}
	return DeleteOneDoc(mconn, collname, filter)
}

// ---------------------------------------------------------------------- Geojson ----------------------------------------------------------------------

// Create

func PostPoint(mconn *mongo.Database, collection string, pointdata GeoJsonPoint) interface{} {
	return InsertOneDoc(mconn, collection, pointdata)
}

func PostLinestring(mconn *mongo.Database, collection string, linestringdata GeoJsonLineString) interface{} {
	return InsertOneDoc(mconn, collection, linestringdata)
}

func PostPolygon(mconn *mongo.Database, collection string, polygondata GeoJsonPolygon) interface{} {
	return InsertOneDoc(mconn, collection, polygondata)
}

// Read

func GetAllBangunan(mconn *mongo.Database, collname string) []GeoJson {
	return GetAllDoc[[]GeoJson](mconn, collname)
}

// Update

// Delete

func DeleteGeojson(mconn *mongo.Database, collname string, userdata User) interface{} {
	filter := bson.M{"username": userdata.Username}
	return DeleteOneDoc(mconn, collname, filter)
}

func GeoIntersects(mconn *mongo.Database, collname string, coordinates Point) string {
	return GetGeoIntersectsDoc(mconn, collname, coordinates)
}

func GeoWithin(mconn *mongo.Database, collname string, coordinates Polygon) string {
	return GetGeoWithinDoc(mconn, collname, coordinates)
}

func Near(mconn *mongo.Database, collname string, coordinates Point) string {
	return GetNearDoc(mconn, collname, coordinates)
}

func NearSphere(mconn *mongo.Database, collname string, coordinates Point) string {
	return GetNearSphereDoc(mconn, collname, coordinates)
}

func Box(mconn *mongo.Database, collname string, coordinates Polyline) string {
	return GetBoxDoc(mconn, collname, coordinates)
}

func Center(mconn *mongo.Database, collname string, coordinates Point) string {
	return GetCenterDoc(mconn, collname, coordinates)
}

func CenterSphere(mconn *mongo.Database, collname string, coordinates Point) string {
	return GetCenterSphereDoc(mconn, collname, coordinates)
}