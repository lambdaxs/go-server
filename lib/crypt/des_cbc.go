package crypt

import (
    "bytes"
    "crypto/cipher"
    "crypto/des"
)

//加密数据
func DES_CBC_Encrypt(data, key, iv []byte) ([]byte, error) {
    aesBlockEncrypter, err := des.NewCipher(key[:8])
    if err != nil {
        return nil, err
    }
    content := PKCS5Padding(data, aesBlockEncrypter.BlockSize())
    encrypted := make([]byte, len(content))
    aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, iv)
    aesEncrypter.CryptBlocks(encrypted, content)
    return encrypted, nil
}

//解密数据
func DES_CBC_Decrypt(src ,key , iv []byte) (data []byte, err error) {
    decrypted := make([]byte, len(src))
    var aesBlockDecrypter cipher.Block
    aesBlockDecrypter, err = des.NewCipher(key[:8])
    if err != nil {
        return nil, err
    }
    aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, iv)
    aesDecrypter.CryptBlocks(decrypted, src)
    return PKCS5Trimming(decrypted), nil
}

/**
PKCS5包装
*/
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(cipherText, padText...)
}

/*
解包装
*/
func PKCS5Trimming(encrypt []byte) []byte {
    padding := encrypt[len(encrypt)-1]
    return encrypt[:len(encrypt)-int(padding)]
}
