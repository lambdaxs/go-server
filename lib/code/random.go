package code

import (
    "encoding/base64"
    "fmt"
    "math/rand"
    "time"
    cryptorand "crypto/rand"
)

func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := cryptorand.Read(b)
    if err != nil {
        return nil, err
    }
    return b, nil
}

func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}

// create 4 captcha code
func CaptchaCode() string {
    return fmt.Sprintf("%04d", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}

