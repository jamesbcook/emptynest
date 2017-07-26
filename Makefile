install:
        mkdir server
        go build -buildmode=plugin -o server/http.so transports/http/main.go
        go build -buildmode=plugin -o server/base32.so encoders/base32/main.go
        go build -buildmode=plugin -o server/base64.so encoders/base64/main.go
        go build -buildmode=plugin -o server/hex.so encoders/hex/main.go
        go build -buildmode=plugin -o server/des.so crypto/des/main.go
        go build -buildmode=plugin -o server/rc4.so crypto/rc4/main.go
        go build -buildmode=plugin -o server/aes_ctr.so crypto/aes_ctr/main.go
        go build -buildmode=plugin -o server/basic.so info/basic/main.go
        mkdir server/plugins
        go build -buildmode=plugin -o server/plugins/shellcode.so plugins/shellcode/main.go
        cd cmd/menu && go get -v && go build -v
        mv cmd/menu/menu server/
        cp config.toml server/
        cp http.toml server/
clean:
        rm -r server
