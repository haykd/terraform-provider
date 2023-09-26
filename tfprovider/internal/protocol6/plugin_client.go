package protocol6

import (
	"context"

	"github.com/apparentlymart/terraform-provider/internal/tfplugin6"
	"google.golang.org/grpc"
)

type PluginClient struct{}

func (c PluginClient) ClientProxy(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
	return tfplugin6.NewProviderClient(conn), nil
}
