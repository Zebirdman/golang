package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
	"strings"
)

const (
	name = "dimager"
)

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
	oldPrefix = repP.Operand
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
	fmt.Printf("Enabled arguments:\n")
	for _, opt := range cmdOptions {
		if opt.Enabled {
			fmt.Printf("Name: %s  Operand: %s\n", opt.Name, opt.Operand)
		}
	}
	// get images list
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	// cycle through images
	for _, image := range images {
		// cycle through the tag names not the actual images
		for _, oTag := range image.RepoTags {

			if addP.Enabled && !repP.Enabled {
				nt := addPrefix(addP.Operand, oTag)
				err := cli.ImageTag(ctx, oTag, nt)
				if err != nil {
					fmt.Printf("%s\n", err)
					//os.Exit(1)
				} else {
					checkClean(ctx, cli, oTag)
				}
			}

			if addP.Enabled && repP.Enabled {
				if pos := strings.IndexRune(oTag, '/'); pos > -1 {
					rem := strings.SplitAfterN(oTag, "/", 2)
					if trim(rem[0]) == trim(repP.Operand) {
						nt := strings.Join([]string{addP.Operand, rem[1]}, "/")
						fmt.Printf("%s ====> %s\n", oTag, nt)
						err := cli.ImageTag(ctx, oTag, nt)
						if err != nil {
							fmt.Printf("%s\n", err)
							//os.Exit(1)
						} else {
							checkClean(ctx, cli, oTag)
						}
					}
				}
			}
		}
	}
	ctx.Done()
}

func checkClean(ctx context.Context, cli *client.Client, t string) {
	// if deletion of old tags is set
	if clean.Enabled {
		ro := types.ImageRemoveOptions{true, false}
		_, err := cli.ImageRemove(ctx, t, ro)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("Old tag %s removed\n", t)
		}
	}
}

// trims the / from the end of a string if it exists
func trim(s string) string {
	if pos := strings.IndexRune(s, '/'); pos > -1 {
		s = strings.Trim(s, "/")
	}
	return s
}

// preappends prefix and returns the new tag
func addPrefix(p, t string) string {
	p = trim(p)
	ar := []string{p, t}
	return strings.Join(ar, "/")
}
