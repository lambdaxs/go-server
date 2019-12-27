package code

import (
    "crypto/md5"
    "encoding/hex"
)

// create md5 str
func MD5Str(str string) string  {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}
