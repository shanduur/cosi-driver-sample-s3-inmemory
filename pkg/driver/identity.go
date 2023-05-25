package driver

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"

	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

type IdentityServer struct {
	Provisioner string
}

func (id *IdentityServer) DriverGetInfo(ctx context.Context,
	req *cosi.DriverGetInfoRequest) (*cosi.DriverGetInfoResponse, error) {

	if id.Provisioner == "" {
		klog.ErrorS(errors.New("provisioner name cannot be empty"), "Invalid argument")
		return nil, status.Error(codes.InvalidArgument, "ProvisionerName is empty")
	}

	return &cosi.DriverGetInfoResponse{
		Name: id.Provisioner,
	}, nil
}
