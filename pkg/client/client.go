package client

import (
	"context"

	// official docker client leaks memory with
	// the curent code, so use fsouza's client
	// dc "github.com/docker/docker/client"
	dc "github.com/fsouza/go-dockerclient"
)

type Client struct {
	ctx context.Context
	cli *dc.Client
}

func NewClient() (*Client, error) {
	// cli, err := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	cli, err := dc.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return &Client{
		ctx: ctx,
		cli: cli,
	}, nil
}

// func (c *Client) GetContainerList() []types.Container {
// res, err := c.cli.ContainerList(c.ctx, container.ListOptions{
func (c *Client) GetContainerList() []dc.APIContainers {
	res, err := c.cli.ListContainers(dc.ListContainersOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	return res
}

func (c *Client) Restart(cId string) bool {
	// err := c.cli.ContainerRestart(c.ctx, cId, container.StopOptions{})
	err := c.cli.RestartContainer(cId, 10)
	return err == nil
}
