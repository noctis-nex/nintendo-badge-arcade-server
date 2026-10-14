package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	nex "github.com/PretendoNetwork/nex-go"
	"os"
)

func passwordFromPID(pid uint32) (string, uint32) {
	secret, err := hex.DecodeString(os.Getenv("PN_NEX_PASSWORD_SECRET"))
	if err != nil || len(secret) < 32 {
		return "", nex.Errors.RendezVous.InvalidUsername
	}
	pidBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(pidBytes, uint64(pid))
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(pidBytes)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), 0
}
