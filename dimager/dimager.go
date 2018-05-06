package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	name = "dimager"
)

//testing one tow three

// save current env variables for restoration at end if needed
var (
	dockHost     = os.Getenv("DOCKER_HOST")
	dockCertPath = os.Getenv("DOCKER_CERT_PATH")
	dockTLS      = os.Getenv("DOCKER_TLS_VERIFY")
	newPrefix    string
	oldPrefix    string

	// option for specifying a docker host
	host = newOption("h", true, []string{"-p"})
	// option for docker cert path when using tls
	path = newOption("c", true, []string{"-h"})
	// option for specifying the image prefix to add
	addP = newOption("p", true, nil)
	// option to specify that existing image tag prefix will be replaced
	repP = newOption("r", true, nil)
	// option to specify that the existing image tags be removed after rename
	clean = newOption("d", false, nil)
	// option for verbose output from operation
	verb = newOption("v", false, nil)

	hp = `dimager: allows for the easy renaming of docker image tags prefix's
  usefull if we want to retag images to use with a pivate registry
  author: Ben Futterleib

Usage: dimager [OPTION]... [-s] SCRIPT_NAME (1st form)
  or: dimager [OPTION]... [-h] DOCKER_HOST [-p] DOCKER_CERT_PATH (2nd form)
In the first form specify a path to a script containing exports for the
  DOCKER_HOST, DOCKER_CERT_PATH, DOCKER_TLS_VERIFY
environment variables, must be valid executable script
In the second form pass the DOCKER_HOST and DOCKER_CERT_PATH values using
  the given flags, DOCKER_TLS_VERIFY will be set automatically

Arguments:
  -h 	specify the connection for the docker host
  -p 	specify the path to the directory holding the client certs for docker
  -s 	give a path to a script containing the relevant export env variables
  --help 	display this help and exit
`
)

func assignVariables() {
	os.Setenv("DOCKER_HOST", host.Operand)
	os.Setenv("DOCKER_CERT_PATH", path.Operand)
	newPrefix = addP.Operand

}
func main() {
	initArgs(name, hp)
	op, err := checkArgs(os.Args)
	if err != nil || help.Enabled {
		showErrors(op, err)
		os.Exit(1)
	}
	// initialize context and a new client
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	// show enabled options and Arguments
	for _, opt := range cmdOptions {
		fmt.Printf("Enabled arguments:\n")
		if opt.Enabled {
			fmt.Printf("Name: %s  Operand: %s\n", opt.Name, opt.Operand)
		}
	}

	// get images list
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {

		for _, tag := range image.RepoTags {
			//fmt.Printf("%s %s 	\n", image.ID, tag)
			if tag == "test:latest" {
				ro := types.ImageRemoveOptions{true, false}
				_, err := cli.ImageRemove(ctx, tag, ro)
				if err != nil {
					panic(err)
				}
				//for _, r := range resp {
				//fmt.Printf("Deleted = %s\nUntagged = %s\n", r.Deleted, r.Untagged)
				//}
			}
		}
	}

	ctx.Done()
}
