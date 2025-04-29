package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithVersion("1.47"))
	if err != nil {
		panic(err)
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "easynode-watchtower-1")
	filterArgs.Add("name", "easynode-easynode-1")
	containers, err := cli.ContainerList(ctx, container.ListOptions{Filters: filters.NewArgs()})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		log.Println("No containers found with the specified label.")
		return
	}

	opt := container.StopOptions{}
	for _, container := range containers {
		log.Printf("Restarting container %s (%s)\n", container.Names[0], container.ID[:12])
		if err := cli.ContainerRestart(ctx, container.ID, opt); err != nil {
			log.Printf("Failed to restart container %s: %v", container.ID[:12], err)
		}
	}

	log.Println("All matching containers have been restarted.")
}
