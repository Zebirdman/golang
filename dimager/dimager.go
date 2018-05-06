package main

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var exMes = "program finished"

func main() {
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
