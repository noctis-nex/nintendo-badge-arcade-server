package nex_datastore

import (
	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/database"
	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/globals"
	"log"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func GetPersistenceInfo(err error, client *nex.Client, callID uint32, ownerID uint32, persistenceSlotID uint16) {
	dataID := database.GetDataStorePersistenceInfo(ownerID, persistenceSlotID)
	log.Printf("badge datastore GetPersistenceInfo pid=%d owner_id=%d slot=%d data_id=%d decode_error=%v", client.PID(), ownerID, persistenceSlotID, dataID, err)

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	if dataID != 0 {
		pPersistenceInfo := datastore.NewDataStorePersistenceInfo()
		pPersistenceInfo.OwnerID = ownerID
		pPersistenceInfo.PersistenceSlotID = persistenceSlotID
		pPersistenceInfo.DataID = uint64(dataID)

		rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

		rmcResponseStream.WriteStructure(pPersistenceInfo)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(datastore.MethodGetPersistenceInfo, rmcResponseBody)
	} else {
		rmcResponse.SetError(nex.Errors.DataStore.NotFound)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
