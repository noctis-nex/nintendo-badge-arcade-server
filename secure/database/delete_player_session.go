package database

import (
	"context"

	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/globals"
	"go.mongodb.org/mongo-driver/bson"
)

func DeletePlayerSession(pid uint32) {
	_, err := sessionsCollection.DeleteOne(context.TODO(), bson.D{{"pid", pid}})
	if err != nil {
		globals.Logger.Error(err.Error())
	}
}
