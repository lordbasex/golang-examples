# golang-examples
Go Lang - Examples

![golang-examples](https://github.com/lordbasex/golang-examples/blob/main/gophers.png|width=100)

GoLang Example is a GNU project providing a collection of useful functions for those looking to learn GoLang. From basic operations to advanced techniques, this library serves as an invaluable guide to mastering the language. Explore and enhance your Go skills with ease!


## ARM - Oracle

```bash
cd /usr/src/
wget https://go.dev/dl/go1.20.2.linux-arm64.tar.gz -O /usr/src/go1.20.2.linux-arm64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-arm64.tar.gz

cat >> /root/.bashrc <<ENDLINE

#golang
export GOROOT=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go
export GOBIN=/root/go/bin
alias godev='nodemon --exec go run main.go --signal SIGTERM'
ENDLINE

```

## INTEL

```bash
cd /usr/src/
wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz -O /usr/src/go1.20.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz
mkdir -p /root/go/{bin,pkg,src}

cat >> /root/.bashrc <<ENDLINE

#golang
export GOROOT=/usr/local/go
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go
export GOBIN=/root/go/bin
alias godev='nodemon --exec go run main.go --signal SIGTERM'
ENDLINE

```
