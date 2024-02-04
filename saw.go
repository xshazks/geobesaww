package geobesaww

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/Befous/BackendGin/models"
)

func ReturnStruct(DataStuct any) (result string) {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func ReturnString(geojson []FullGeoJson) string {
	var names []string
	for _, geojson := range geojson {
		names = append(names, geojson.Properties.Name)
	}
	result := strings.Join(names, ", ")
	return result
}

// ----------------------------------------------------------------------- User -----------------------------------------------------------------------

func Authorization(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response CredentialUser
	var auth User
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenname := DecodeGetName(os.Getenv(publickey), header)
	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	auth.Username = tokenusername

	if tokenname == "" || tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, auth) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	response.Message = "Berhasil decode token"
	response.Status = true
	response.Data.Name = tokenname
	response.Data.Username = tokenusername
	response.Data.Role = tokenrole

	return ReturnStruct(response)
}

func Registrasi(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	if UsernameExists(mconn, collname, datauser) {
		response.Message = "Username telah dipakai"
		return ReturnStruct(response)
	}

	hash, hashErr := HashPassword(datauser.Password)
	if hashErr != nil {
		response.Message = "Gagal hash password: " + hashErr.Error()
		return ReturnStruct(response)
	}

	datauser.Password = hash

	InsertUser(mconn, collname, datauser)
	response.Status = true
	response.Message = "Berhasil input data"

	return ReturnStruct(response)
}

func Login(privatekey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, datauser) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if !IsPasswordValid(mconn, collname, datauser) {
		response.Message = "Password Salah"
		return ReturnStruct(response)
	}

	user := FindUser(mconn, collname, datauser)

	tokenstring, tokenerr := Encode(user.Name, user.Username, user.Role, os.Getenv(privatekey))
	if tokenerr != nil {
		response.Message = "Gagal encode token: " + tokenerr.Error()
		return ReturnStruct(response)
	}

	response.Status = true
	response.Message = "Berhasil login"
	response.Token = tokenstring

	return ReturnStruct(response)
}

func AmbilSemuaUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	datauser := GetAllUser(mconn, collname)
	response.Status = true
	response.Message = "Berhasil mengambil data"
	response.Data = datauser
	return ReturnStruct(response)
}

func EditUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	if datauser.Username == "" {
		response.Message = "Parameter dari function ini adalah username"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, datauser) {
		response.Message = "Akun yang ingin diedit tidak ditemukan"
		return ReturnStruct(response)
	}

	if datauser.Password != "" {
		hash, hashErr := HashPassword(datauser.Password)
		if hashErr != nil {
			response.Message = "Gagal Hash Password: " + hashErr.Error()
			return ReturnStruct(response)
		}
		datauser.Password = hash
	} else {
		user := FindUser(mconn, collname, datauser)
		datauser.Password = user.Password
	}

	UpdateUser(mconn, collname, datauser)

	response.Status = true
	response.Message = "Berhasil update " + datauser.Username + " dari database"
	return ReturnStruct(response)
}

func HapusUser(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		response.Message = "Anda tidak memiliki akses"
		return ReturnStruct(response)
	}

	if datauser.Username == "" {
		response.Message = "Parameter dari function ini adalah username"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, datauser) {
		response.Message = "Akun yang ingin dihapus tidak ditemukan"
		return ReturnStruct(response)
	}

	DeleteUser(mconn, collname, datauser)

	response.Status = true
	response.Message = "Berhasil hapus " + datauser.Username + " dari database"
	return ReturnStruct(response)
}

// ---------------------------------------------------------------------- Geojson ----------------------------------------------------------------------

func MembuatGeojsonPoint(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var geojsonpoint GeoJsonPoint
	err := json.NewDecoder(r.Body).Decode(&geojsonpoint)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		if tokenrole != "dosen" {
			response.Message = "Anda tidak memiliki akses"
			return ReturnStruct(response)
		}
	}

	PostPoint(mconn, collname, geojsonpoint)
	response.Status = true
	response.Message = "Data point berhasil masuk"

	return ReturnStruct(response)
}

func MembuatGeojsonPolyline(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var geojsonpolyline GeoJsonLineString
	err := json.NewDecoder(r.Body).Decode(&geojsonpolyline)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		if tokenrole != "dosen" {
			response.Message = "Anda tidak memiliki akses"
			return ReturnStruct(response)
		}
	}

	PostLinestring(mconn, collname, geojsonpolyline)
	response.Status = true
	response.Message = "Data polyline berhasil masuk"

	return ReturnStruct(response)
}

func MembuatGeojsonPolygon(publickey, mongoenv, dbname, collname string, r *http.Request) string {
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)
	var geojsonpolygon GeoJsonPolygon
	err := json.NewDecoder(r.Body).Decode(&geojsonpolygon)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	header := r.Header.Get("token")
	if header == "" {
		response.Message = "Header login tidak ditemukan"
		return ReturnStruct(response)
	}

	tokenusername := DecodeGetUsername(os.Getenv(publickey), header)
	tokenrole := DecodeGetRole(os.Getenv(publickey), header)

	if tokenusername == "" || tokenrole == "" {
		response.Message = "Hasil decode tidak ditemukan"
		return ReturnStruct(response)
	}

	if !UsernameExists(mconn, collname, User{Username: tokenusername}) {
		response.Message = "Akun tidak ditemukan"
		return ReturnStruct(response)
	}

	if tokenrole != "owner" {
		if tokenrole != "dosen" {
			response.Message = "Anda tidak memiliki akses"
			return ReturnStruct(response)
		}
	}

	PostPolygon(mconn, collname, geojsonpolygon)
	response.Status = true
	response.Message = "Data polygon berhasil masuk"

	return ReturnStruct(response)
}

func AmbilDataGeojson(mongoenv, dbname, collname string, r *http.Request) string {
	mconn := SetConnection(mongoenv, dbname)
	datagedung := GetAllBangunan(mconn, collname)
	return ReturnStruct(datagedung)
}

func PostGeoIntersects(mongoenv, dbname, collname string, r *http.Request) string {
	var geospatial models.Geospatial
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	err := json.NewDecoder(r.Body).Decode(&geospatial)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	geointersects, err := GeoIntersects(mconn, collname, geospatial)

	if err != nil {
		response.Message = "Error : " + err.Error()
		return ReturnStruct(response)
	}

	result := ReturnString(geointersects)

	response.Status = true
	response.Message = result

	return ReturnStruct(response)
}

func PostGeoWithin(mongoenv, dbname, collname string, r *http.Request) string {
	var coordinate Polygon
	var response Pesan
	response.Status = false
	mconn := SetConnection(mongoenv, dbname)

	err := json.NewDecoder(r.Body).Decode(&coordinate)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	response.Status = true
	response.Message = GeoWithin(mconn, collname, coordinate)

	return ReturnStruct(response)
}

func PostNear(mongoenv, dbname, collname string, r *http.Request) string {
	var coordinate Point
	var response Pesan
	response.Status = false
	mconn := SetConnection2dsphere(mongoenv, dbname, collname)

	err := json.NewDecoder(r.Body).Decode(&coordinate)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	response.Status = true
	response.Message = Near(mconn, collname, coordinate)

	return ReturnStruct(response)
}

func PostNearSphere(mongoenv, dbname, collname string, r *http.Request) string {
	var coordinate Point
	var response Pesan
	response.Status = false
	mconn := SetConnection2dsphere(mongoenv, dbname, collname)

	err := json.NewDecoder(r.Body).Decode(&coordinate)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	response.Status = true
	response.Message = NearSphere(mconn, collname, coordinate)

	return ReturnStruct(response)
}

func PostBox(mongoenv, dbname, collname string, r *http.Request) string {
	var coordinate Polyline
	var response Pesan
	response.Status = false
	mconn := SetConnection2dsphere(mongoenv, dbname, collname)

	err := json.NewDecoder(r.Body).Decode(&coordinate)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	response.Status = true
	response.Message = Box(mconn, collname, coordinate)

	return ReturnStruct(response)
}

func PostCenter(mongoenv, dbname, collname string, r *http.Request) string {
	var coordinate Point
	var response Pesan
	response.Status = false
	mconn := SetConnection2dsphere(mongoenv, dbname, collname)

	err := json.NewDecoder(r.Body).Decode(&coordinate)

	if err != nil {
		response.Message = "Error parsing application/json: " + err.Error()
		return ReturnStruct(response)
	}

	response.Status = true
	response.Message = Center(mconn, collname, coordinate)

	return ReturnStruct(response)
}
