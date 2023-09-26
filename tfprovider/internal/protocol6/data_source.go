package protocol6

import (
	"context"

	"github.com/apparentlymart/terraform-provider/internal/tfplugin6"
	"github.com/apparentlymart/terraform-provider/tfprovider/internal/common"
)

type DataResourceType struct {
	client   tfplugin6.ProviderClient
	typeName string
	schema   *common.DataResourceTypeSchema
}

func (dst *DataResourceType) Read(ctx context.Context, req common.DataResourceReadRequest) (common.DataResourceReadResponse, common.Diagnostics) {
	resp := common.DataResourceReadResponse{}
	return resp, nil
}

func (dst *DataResourceType) Sealed() common.Sealed {
	return common.Sealed{}
}
