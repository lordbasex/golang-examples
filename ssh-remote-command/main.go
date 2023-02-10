package main

import (
        "io/ioutil"
        "log"
        "os"
        "time"

        "golang.org/x/crypto/ssh"
)

func failOnError(err error, msg string) {
        if err != nil {
                log.Printf("%s: %s\n", msg, err)
        }
}

func track(name string) func() {
        start := time.Now()
        return func() {
                log.Printf("%s, execution time %s\n", name, time.Since(start))
        }
}

func main() {
        defer track("main")()

        privateKey, err := ioutil.ReadFile("/root/.ssh/id_rsa")
        if err != nil {
                log.Fatal("Failed to load private key: ", err)
        }

        signer, err := ssh.ParsePrivateKey(privateKey)
        if err != nil {
                log.Fatal("Failed to parse private key: ", err)
        }

        config := &ssh.ClientConfig{
                User: "root",
                Auth: []ssh.AuthMethod{
                        ssh.PublicKeys(signer),
                },
                HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        }

        conn, err := ssh.Dial("tcp", "192.168.0.10:22", config)
        if err != nil {
                log.Fatal("Failed to dial: ", err)
        }

        session, err := conn.NewSession()
        if err != nil {
                log.Fatal("Failed to create session: ", err)
        }
        defer session.Close()

        var b []byte
        b, err = session.Output("ls -la")
        if err != nil {
                log.Fatal("Failed to run: " + err.Error())
        }

        log.Print(string(b))

}
