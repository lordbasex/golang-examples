package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {
	// Definir parámetros de conexión SSH
	host := "192.168.0.10"
	port := "22"
	user := "root"
	pass := "mypassword"

	// Crear la configuración del cliente SSH
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second, // Tiempo máximo para establecer la conexión
	}

	// Conectar a través de SSH
	client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	// Crear una sesión SSH
	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()

	// Convertir la entrada de la terminal a un formato reconocido por el destino de la conexión
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(fd, state)

	// Obtener el tamaño del terminal
	w, h, err := term.GetSize(fd)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Configurar modos del terminal SSH
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", h, w, modes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Redirigir la entrada/salida estándar de la sesión SSH a la terminal local
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// Iniciar la sesión de shell SSH
	err = session.Shell()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Detectar y procesar cambios en el tamaño del terminal local
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGWINCH)
	go func() {
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGWINCH:
				fd := int(os.Stdout.Fd())
				w, h, _ = term.GetSize(fd)
				session.WindowChange(h, w)
			}
		}
	}()

	// Esperar a que la sesión SSH termine
	err = session.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
