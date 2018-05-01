package server

import (
	"fmt"
	"net"
	"net/http"

	auth "github.com/docker/go-plugins-helpers/authorization"
	//"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	//"os"
	//dockerapi "github.com/docker/docker/api"
	//dockerclient "github.com/docker/docker/client"
)

const (
	pluginName   = "zeft-auth"
	pluginFolder = "/run/docker/plugins"
)

// ZeftServer server structure
type ZeftServer struct {
	listener net.Listener
}

// NewZeftServer server structure
func NewZeftServer() *ZeftServer {
	return &ZeftServer{}
}

// Start the auth server
func (s *ZeftServer) Start() error {

	pluginPath := fmt.Sprintf("%s/%s.sock", pluginFolder, pluginName)
	s.listener, _ = net.ListenUnix("unix", &net.UnixAddr{Name: pluginPath, Net: "unix"})
	//if err != nil {
	//  return err
	//}
	return http.Serve(s.listener, nil)
}

// implement the docker access auth interface functions
func (s *ZeftServer) AuthZReq(req auth.Request) auth.Response {

	log.WithFields(log.Fields{
		"URI":  req.RequestURI,
		"User": req.User,
	}).Info("Request recieved")

	return auth.Response{Allow: true, Msg: "Responded", Err: "none"}
}

//
func (s *ZeftServer) AuthZRes(req auth.Request) auth.Response {
	return auth.Response{Allow: true, Msg: "Responded", Err: "none"}
}
