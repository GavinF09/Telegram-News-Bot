package connect_mongodb

import (
	"crypto/tls"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(certFile string, keyFile string, dbURI string) *mongo.Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	opts := options.Client().ApplyURI(dbURI).SetTLSConfig(tlsConfig)
	dbClient, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	return dbClient
}
