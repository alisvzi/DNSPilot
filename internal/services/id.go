package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func newID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s_%s_%s", prefix, time.Now().UTC().Format("20060102150405"), hex.EncodeToString(b))
}
