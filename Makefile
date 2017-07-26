install:
	mkdir server
	go build -buildmode=plugin -o server/http.so plugins/transports/http/main.go
	go build -buildmode=plugin -o server/base32.so plugins/encoders/base32/main.go
	go build -buildmode=plugin -o server/base64.so plugins/encoders/base64/main.go
	go build -buildmode=plugin -o server/hex.so plugins/encoders/hex/main.go
	go build -buildmode=plugin -o server/des.so plugins/crypto/des/main.go
	go build -buildmode=plugin -o server/rc4.so plugins/crypto/rc4/main.go
	go build -buildmode=plugin -o server/aes_ctr.so plugins/crypto/aes_ctr/main.go
	go build -buildmode=plugin -o server/basic.so plugins/info/basic/main.go
	mkdir server/plugins
	go build -buildmode=plugin -o server/plugins/shellcode.so plugins/shellcode/main.go
	cd cmd/menu && go get -v && go build -v
	mv cmd/menu/menu server/
	cp config.toml server/
	cp http.toml server/
clean:
	rm -r server
