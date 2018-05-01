package main

import (
	//  "bytes"
	//  "fmt"
	"os"
	//  "os/user"
	//  "strconv"
	//  "net"
	//  "net/http"
	srv "socket-listener/server"
	//auth  "github.com/docker/go-plugins-helpers/authorization"
	log "github.com/sirupsen/logrus"
)

const (
	//defaultDockerHost = "unix:///var/run/docker.sock"
	pluginSocket = "/run/docker/plugins/zeft-auth.sock"
)

// initialize unix socket and logger
func init() {
	//buffer := bytes.Buffer

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})
	log.SetOutput(os.Stdout)
}

func main() {

	server := srv.NewZeftServer()
	server.Start()

}
