// alert

package utils

import (
	"bytes"
	"crypto"
	"encoding/hex"

	//"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func HostsWithoutThis(hostThis string, hosts []string) []string {
	// hostWithoutThis := make([]string)
	var hostWithoutThis []string
	for _, host := range hosts {
		if host != hostThis {
			hostWithoutThis = append(hostWithoutThis, host)
		}
	}
	return hostWithoutThis
}

func FindProcessPidByName(processName string) []int {
	var pids []int
	fd, _ := ioutil.ReadDir("/proc")
	for _, fi := range fd {
		fiName := fi.Name()
		pid, err := strconv.Atoi(fiName)
		if err == nil {
			statusFile := path.Join("/proc", fiName, "status")
			f, err := ioutil.ReadFile(statusFile)
			if err != nil {
				continue
			}

			name := string(f[6:bytes.IndexByte(f, '\n')])
			if name == processName {
				pids = append(pids, pid)
			}

		} else {
			continue
		}
	}
	return pids
}
func KillProcess(pid int) {
	proc, _ := os.FindProcess(pid)
	proc.Kill()
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

func GetCode() string {
	mac := GetMacAddr()
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(mac))
	code := Sha1Inst.Sum([]byte(""))
	codeString := hex.EncodeToString(code)
	return codeString
}
func VerifySign(path string, publicKeyString string, signText []byte, plainText []byte) bool {
	//首先从文件中提取公钥
	var buf []byte
	if len(publicKeyString) > 0 {
		buf = []byte(publicKeyString)
	} else {
		fp, _ := os.Open(path)
		defer fp.Close()
		//测量文件长度以便于保存
		fileinfo, _ := fp.Stat()
		buf := make([]byte, fileinfo.Size())
		fp.Read(buf)

	}
	block, _ := pem.Decode(buf)
	//x509解码,得到一个interface类型的pub
	pub, _ := x509.ParsePKIXPublicKey(block.Bytes)
	//签名函数中需要的数据散列值
	hash := sha256.Sum256(plainText)
	//验证签名
	err := rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash[:], signText)
	if err != nil {
		return false
	} else {
		return true //success
	}

}
