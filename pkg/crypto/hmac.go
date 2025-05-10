package crypto

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
)

func GenerateHMAC(data string, secret []byte) string {
    h := hmac.New(sha256.New, secret)
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data string, secret []byte, expectedMAC string) bool {
    actualMAC := GenerateHMAC(data, secret)
    return hmac.Equal([]byte(actualMAC), []byte(expectedMAC))
}
