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

var (
	// save current env variables for restoration at end if needed
	dockHost     = os.Getenv("DOCKER_HOST")
	dockCertPath = os.Getenv("DOCKER_CERT_PATH")
	dockTLS      = os.Getenv("DOCKER_TLS_VERIFY")

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
	// option to remove a existing prefix
	remove = newOption("x", true, nil)

	hp = `dim: Docker Image Manager
	allows for the easy renaming of docker image tags prefix's
  usefull if we want to retag images to use with a pivate registry
  author: Ben Futterleib

Usage: dimager [OPTION]... [-h] DOCKER_HOST [-p] DOCKER_CERT_PATH
Can pass the DOCKER_HOST and DOCKER_CERT_PATH values using
  the given flags, DOCKER_TLS_VERIFY will be set automatically if certpath
	is set

Arguments:
  -h 	[DOCKER_HOST] specify the connection for the docker host
  -c 	[DOCKER_CERT_PATH] specify the path to the directory holding the client
	  certs for docker
	-p [NEW PREFIX] tag all images with a new prefix
	-r [EXISTING TAG]specify a existing tag prefix to replace, designate the new
	  prefix with the '-d' flag
	-d set the old tags for deletion so only the newly modified tags are kept
	-v verbose output
	-x [PREFIX] specfiy a prefix to remove
`
)

func main() {
	initArgs(name, hp)
	op, err := checkArgs(os.Args)
	if err != nil || help.Enabled {
		showErrors(op, err)
		os.Exit(1)
	}
	// initialize context and a new client
	assignVariables()
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if remove.Enabled {
		clean.Enabled = true
	}

	// get images list
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	// verbose output headers
	vOutput("%-50s %-50s", "EXISTING TAG", "NEW TAG")
	if clean.Enabled {
		vOutput("DELETE EXISTING")
	}
	vOutput("\n")

	// cycle through images
	for _, image := range images {
		// cycle through the tag names not the actual images
		for _, oTag := range image.RepoTags {

			// TODO: optimize this if possible
			if addP.Enabled && !repP.Enabled {
				nt := addPrefix(addP.Operand, oTag)
				processTags(ctx, cli, oTag, nt)
			}

			if addP.Enabled && repP.Enabled {
				if pos := strings.IndexRune(oTag, '/'); pos > -1 {
					rem := strings.SplitAfterN(oTag, "/", 2)
					if trim(rem[0]) == trim(repP.Operand) {
						nt := strings.Join([]string{addP.Operand, rem[1]}, "/")
						processTags(ctx, cli, oTag, nt)
					}
				}
			}
			if remove.Enabled {
				if pos := strings.IndexRune(oTag, '/'); pos > -1 {
					rem := strings.SplitAfterN(oTag, "/", 2)
					if trim(rem[0]) == trim(remove.Operand) {
						processTags(ctx, cli, oTag, rem[1])
					}
				}
			}
		}
	}
	restoreVariables()
	ctx.Done()
}

// tag the images and clean up old ones if specified
func processTags(ctx context.Context, cli *client.Client, t, nt string) {
	err := cli.ImageTag(ctx, t, nt)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	vOutput("%-50s %-50s", t, nt)
	if clean.Enabled {
		ro := types.ImageRemoveOptions{true, false}
		_, err := cli.ImageRemove(ctx, t, ro)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			vOutput("succesfull")
		}
	}
	vOutput("\n")
}

// wraps Printf so info is only displayed when verbose mode enabled
func vOutput(a string, s ...string) {
	if verb.Enabled {
		switch len(s) {
		case 0:
			fmt.Printf(a)
			break
		case 1:
			fmt.Printf(a, s[0])
			break
		case 2:
			fmt.Printf(a, s[0], s[1])
			break
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

// set env variables if option specified
func assignVariables() {
	if host.Enabled {
		os.Setenv("DOCKER_HOST", host.Operand)
	}
	if path.Enabled {
		os.Setenv("DOCKER_CERT_PATH", path.Operand)
		os.Setenv("DOCKER_TLS_VERIFY", "1")
	}
}

// restore the env variables from before
func restoreVariables() {
	os.Setenv("DOCKER_HOST", dockHost)
	os.Setenv("DOCKER_CERT_PATH", dockCertPath)
	os.Setenv("DOCKER_TLS_VERIFY", dockTLS)
}
