package utility

import (
	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/database"
	"github.com/PretendoNetwork/nex-go"
)

func ChangeFreePlayDataMeta(dataID uint64, metaBinary []byte) {
	dateTime := nex.NewDateTime(0)
	updatedTime := dateTime.Now()

	database.UpdateFreePlayDataMetaBinary(uint32(dataID), metaBinary, updatedTime)
}
