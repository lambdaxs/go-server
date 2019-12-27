package code

import "encoding/base64"

// create base64 str
func Base64Str(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}

// parse base64 str
func UnBase64Str(str string) (res string,err error) {
    buf,err := base64.StdEncoding.DecodeString(str)
    return string(buf),err
}
