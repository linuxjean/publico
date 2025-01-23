sudo apt update
sudo apt install -y golang-go

go mod init bip39-address-generator
go mod tidy
go run main.go
