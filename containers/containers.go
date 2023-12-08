package containers

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/chiyoi/atri/env"
	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/az/identity"
	"github.com/chiyoi/iter/res"
)

const (
	Channels    = "channels"
	QueuedSteps = "queued_steps"
)

var ContainerSchemes = map[string]cosmos.Schema{
	Channels:    {PartitionKeyPath: "/id"},
	QueuedSteps: {PartitionKeyPath: "/id"},
}

func Client(containerID string) (c *azcosmos.ContainerClient, err error) {
	client, err := databaseClient()
	return res.R(containerID, err, client.NewContainer)
}

func databaseClient() (client *azcosmos.DatabaseClient, err error) {
	cred, err := identity.DefaultCredential()
	c, err := res.R(cred, err, cosmos.NewClient(env.EndpointCosmos, nil))
	return res.R(env.Database, err, c.NewDatabase)
}
