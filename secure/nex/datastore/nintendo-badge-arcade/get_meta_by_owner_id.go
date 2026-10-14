package nex_datastore_nintendo_badge_arcade

import (
	"fmt"

	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/globals"
	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/utility"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
	datastore_nintendo_badge_arcade "github.com/PretendoNetwork/nex-protocols-go/datastore/nintendo-badge-arcade"
)

func GetMetaByOwnerID(err error, client *nex.Client, callID uint32, param *datastore_nintendo_badge_arcade.DataStoreGetMetaByOwnerIDParam) {
	pMetaInfo := make([]*datastore.DataStoreMetaInfo, 0)
	var pHasNext bool = false // Unknown

	if len(param.OwnerIDs) != len(param.DataTypes) {
		// Not sure if this is possible in the first place
		fmt.Println("WARNING: Owner ID and DataType length mismatch")
	}

	for i, _ := range param.OwnerIDs {
		switch param.DataTypes[i] {
		case 100: // Free Play Data
			freePlayDataMetaInfo := utility.FreePlayDataToDataStoreMetaInfo(param.OwnerIDs[i], param.DataTypes[i])
			if freePlayDataMetaInfo != nil {
				pMetaInfo = append(pMetaInfo, freePlayDataMetaInfo)
			}
		default:
			fmt.Printf("WARNING: Unknown DataType: %d\n", param.DataTypes[i])
		}
	}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteListStructure(pMetaInfo)
	rmcResponseStream.WriteBool(pHasNext)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_nintendo_badge_arcade.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_nintendo_badge_arcade.MethodGetMetaByOwnerID, rmcResponseBody)

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
