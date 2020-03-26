package shadow

import (
	"encoding/hex"
	"strings"
)

func UnParseToken(token string) (data, sign []byte) {
	split := strings.Split(token, ".")
	data, _ = hex.DecodeString(split[0])
	sign, _ = hex.DecodeString(split[1])
	return
}
