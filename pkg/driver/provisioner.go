// Copyright 2023 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
	"sigs.k8s.io/cosi-driver-sample/pkg/s3fake"
)

type ProvisionerServer struct {
	Provisioner string
	S3URL       string
	S3          *s3fake.S3Fake
}

// DriverCreateBucket is an idempotent method for creating buckets
// It is expected to create the same bucket given a bucketName and protocol
// If the bucket already exists, then it MUST return codes.AlreadyExists
// Return values
//
//	nil -                   Bucket successfully created
//	codes.AlreadyExists -   Bucket already exists. No more retries
//	non-nil err -           Internal error                                [requeue'd with exponential backoff]
func (s *ProvisionerServer) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest) (*cosi.DriverCreateBucketResponse, error) {

	if exists, err := s.S3.BucketExists(req.Name); err != nil {
		klog.V(3).Error(err, "DriverCreateBucket: failed to check if bucket exists", "bucket", req.Name)
		return nil, status.Error(codes.Internal, err.Error())
	} else if !exists {
		klog.V(3).Info("DriverCreateBucket: creating bucket", "bucket", req.Name)

		err := s.S3.CreateBucket(req.Name)
		if err != nil {
			klog.V(3).Error(err, "DriverCreateBucket: failed to create bucket", "bucket", req.Name)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	klog.V(3).Info("DriverCreateBucket: bucket created", "bucket", req.Name)
	return &cosi.DriverCreateBucketResponse{
		BucketId: req.Name,
		BucketInfo: &cosi.Protocol{
			Type: &cosi.Protocol_S3{
				S3: &cosi.S3{
					Region:           "eu-central-1",
					SignatureVersion: cosi.S3SignatureVersion_S3V2,
				},
			},
		},
	}, nil
}

func (s *ProvisionerServer) DriverDeleteBucket(ctx context.Context,
	req *cosi.DriverDeleteBucketRequest) (*cosi.DriverDeleteBucketResponse, error) {
	if exists, err := s.S3.BucketExists(req.BucketId); err != nil {
		klog.V(3).Error(err, "DriverDeleteBucket: failed to check if bucket exists", "bucket", req.BucketId)
		return nil, status.Error(codes.Internal, err.Error())
	} else if exists {
		klog.V(3).Info("DriverDeleteBucket: deleting bucket", "bucket", req.BucketId)

		err := s.S3.DeleteBucket(req.BucketId)
		if err != nil {
			klog.V(3).Error(err, "DriverDeleteBucket: failed to delete bucket", "bucket", req.BucketId)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cosi.DriverDeleteBucketResponse{}, nil
}

func (s *ProvisionerServer) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}

func (s *ProvisionerServer) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest) (*cosi.DriverRevokeBucketAccessResponse, error) {

	return nil, status.Error(codes.Unimplemented, "DriverCreateBucket: not implemented")
}
