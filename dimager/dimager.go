package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var exMes = "program finished"

func main() {
	printRunes(os.Args)
	op, err := checkArgs(os.Args)
	if help.Enabled == true {
		helpPage()
	} else if err != nil {
		errorPage(op.Name, err)
	} else {
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
}

func printRunes(a []string) {
	ar := a[1:]
	for _, arg := range ar {
		for index, runVal := range arg {
			fmt.Printf("%#U starts at byte position %d\n", runVal, index)
		}
	}
}
