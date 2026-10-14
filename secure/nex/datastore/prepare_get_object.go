package nex_datastore

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"time"

	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/database"
	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func PrepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *datastore.DataStorePrepareGetParam) {
	pReqGetInfo := datastore.NewDataStoreReqGetInfo()

	dataVersion := database.GetVersionByDataID(uint32(dataStorePrepareGetParam.DataID))

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("PN_NBA_CONFIG_S3_PATH"), dataStorePrepareGetParam.DataID, dataVersion)
	dataSize, err := globals.S3ObjectSize(os.Getenv("PN_NBA_CONFIG_S3_BUCKET"), key)
	if err != nil {
		globals.Logger.Error(err.Error())
	}

	presigner := s3.NewPresignClient(globals.S3Client)
	presigned, signErr := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("PN_NBA_CONFIG_S3_BUCKET")),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) { opts.Expires = 15 * time.Minute })
	if signErr != nil {
		globals.Logger.Error(signErr.Error())
		return
	}
	pReqGetInfo.URL = presigned.URL
	pReqGetInfo.RequestHeaders = []*datastore.DataStoreKeyValue{}
	pReqGetInfo.Size = uint32(dataSize)
	pReqGetInfo.RootCA = []byte{}
	pReqGetInfo.DataID = dataStorePrepareGetParam.DataID

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pReqGetInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPrepareGetObject, rmcResponseBody)

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
