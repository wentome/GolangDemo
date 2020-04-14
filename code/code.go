// code
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

func main() {
	//mac := GetMacAddr()
	//key, _ := base64.StdEncoding.DecodeString("rbaOfTm+yzZfUG1QK1hMKA==")
	key := []byte{0x5c, 0x9a, 0xac, 0x49, 0x64, 0x8f, 0xde, 0x64, 0xa8, 0x58, 0x43, 0xec, 0x9d, 0x25, 0x58, 0xb5}
	code, err := encryptAES([]byte("hello"), key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(strings.ToUpper(hex.EncodeToString(code)))
}
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}
func encryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

func GetMacAddr() string {
	// var macAddress string
	fd, _ := ioutil.ReadDir("/sys/class/net")
	netName := fd[0].Name()
	macAddrFile := path.Join("/sys/class/net", netName, "address")
	f, err := ioutil.ReadFile(macAddrFile)
	if err != nil {
		return ""
	}
	macAddr := string(f[0:bytes.IndexByte(f, '\n')])
	return macAddr
}
