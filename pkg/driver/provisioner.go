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

	"github.com/shanduur/cosi-driver-sample-s3-inmemory/pkg/s3fake"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	cosi "sigs.k8s.io/container-object-storage-interface-spec"
)

type ProvisionerServer struct {
	Provisioner string
	S3URL       string
	S3          *s3fake.S3Fake
}

func (s *ProvisionerServer) DriverCreateBucket(ctx context.Context,
	req *cosi.DriverCreateBucketRequest) (*cosi.DriverCreateBucketResponse, error) {

	if exists, err := s.S3.BucketExists(req.Name); err != nil {
		klog.V(3).ErrorS(err, "DriverCreateBucket: failed to check if bucket exists", "bucket", req.Name)
		return nil, status.Error(codes.Internal, err.Error())
	} else if !exists {
		klog.V(3).InfoS("DriverCreateBucket: creating bucket", "bucket", req.Name)

		err := s.S3.CreateBucket(req.Name)
		if err != nil {
			klog.V(3).ErrorS(err, "DriverCreateBucket: failed to create bucket", "bucket", req.Name)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	klog.V(3).InfoS("DriverCreateBucket: bucket created", "bucket", req.Name)
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
		klog.V(3).ErrorS(err, "DriverDeleteBucket: failed to check if bucket exists", "bucket", req.BucketId)
		return nil, status.Error(codes.Internal, err.Error())
	} else if exists {
		klog.V(3).InfoS("DriverDeleteBucket: deleting bucket", "bucket", req.BucketId)

		err := s.S3.DeleteBucket(req.BucketId)
		if err != nil {
			klog.V(3).ErrorS(err, "DriverDeleteBucket: failed to delete bucket", "bucket", req.BucketId)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	klog.V(3).Info("DriverDeleteBucket: bucket deleted", "bucket", req.BucketId)
	return &cosi.DriverDeleteBucketResponse{}, nil
}

func (s *ProvisionerServer) DriverGrantBucketAccess(ctx context.Context,
	req *cosi.DriverGrantBucketAccessRequest) (*cosi.DriverGrantBucketAccessResponse, error) {

	user, err := s.S3.UserExists(req.Name)
	if err != nil {
		klog.V(3).ErrorS(err, "DriverGrantBucketAccess: failed to check if user exists", "user", req.Name)
		return nil, status.Error(codes.Internal, err.Error())
	} else if user == nil {
		klog.V(3).Info("DriverGrantBucketAccess: creating user", "user", req.Name)

		user, err = s.S3.CreateUser(req.Name)
		if err != nil {
			klog.V(3).ErrorS(err, "DriverGrantBucketAccess: failed to create user", "user", req.Name)
			return nil, status.Error(codes.Internal, err.Error())
		}

		klog.V(5).InfoS("DriverGrantBucketAccess: user created", "user", user)
	}

	klog.V(3).InfoS("DriverGrantBucketAccess: access granted", "user", req.Name, "bucket", req.BucketId)
	return &cosi.DriverGrantBucketAccessResponse{
		AccountId: user.Name,
		Credentials: map[string]*cosi.CredentialDetails{
			"s3": &cosi.CredentialDetails{
				Secrets: map[string]string{
					"accessKeyID":     user.AccessKey,
					"accessSecretKey": user.SecretKey,
				},
			},
		},
	}, nil
}

func (s *ProvisionerServer) DriverRevokeBucketAccess(ctx context.Context,
	req *cosi.DriverRevokeBucketAccessRequest) (*cosi.DriverRevokeBucketAccessResponse, error) {

	user, err := s.S3.UserExists(req.AccountId)
	if err != nil {
		klog.V(3).ErrorS(err, "DriverRevokeBucketAccess: failed to check if user exists", "user", req.AccountId)
		return nil, status.Error(codes.Internal, err.Error())
	} else if user != nil {
		klog.V(3).InfoS("DriverRevokeBucketAccess: deleting user", "user", req.AccountId)

		if err := s.S3.DeleteUser(req.AccountId); err != nil {
			klog.V(3).ErrorS(err, "DriverRevokeBucketAccess: failed to delete user", "user", req.AccountId)
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	klog.V(3).InfoS("DriverRevokeBucketAccess: access revoked", "user", req.AccountId, "bucket", req.BucketId)
	return &cosi.DriverRevokeBucketAccessResponse{}, nil
}
