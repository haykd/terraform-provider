package protocol5

import (
	"context"

	"github.com/apparentlymart/terraform-provider/internal/tfplugin5"
	"github.com/apparentlymart/terraform-provider/tfprovider/internal/common"
)

type DataResourceType struct {
	client   tfplugin5.ProviderClient
	typeName string
	schema   *common.DataResourceTypeSchema
}

func (dst *DataResourceType) Read(ctx context.Context, req common.DataResourceReadRequest) (common.DataResourceReadResponse, common.Diagnostics) {
	resp := common.DataResourceReadResponse{}
	// Maybe
	dv, diags := encodeDynamicValue(req.Config, dst.schema.Content)
	if diags.HasErrors() {
		return resp, diags
	}
	rawResp, err := dst.client.ReadDataSource(ctx, &tfplugin5.ReadDataSource_Request{
		TypeName: dst.typeName,
		Config:   dv,
	})
	diags = append(diags, common.RPCErrorDiagnostics(err)...)
	if err != nil {
		return resp, diags
	}
	diags = append(diags, decodeDiagnostics(rawResp.Diagnostics)...)

	if raw := rawResp.State; raw != nil {
		v, moreDiags := decodeDynamicValue(raw, dst.schema.Content)
		resp.State = v
		diags = append(diags, moreDiags...)
	}

	return resp, diags
}

func (dst *DataResourceType) Sealed() common.Sealed {
	return common.Sealed{}
}
