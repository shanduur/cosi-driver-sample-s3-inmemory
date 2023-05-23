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

package s3fake

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/johannesboyne/gofakes3"
)

// S3Fake is a fake S3 server.
type S3Fake struct {
	Address string
	Backend gofakes3.Backend
	Users   *sync.Map
}

// Run starts the fake S3 server.
func (s *S3Fake) Run(ctx context.Context) error {
	lis := &net.TCPListener{}

	faker := gofakes3.New(s.Backend)

	srv := http.Server{
		Addr:    s.Address,
		Handler: faker.Server(),
	}

	return srv.Serve(lis)
}

// CreateBucket creates a bucket.
func (s *S3Fake) CreateBucket(bucket string) error {
	return s.Backend.CreateBucket(bucket)
}

// BucketExists checks if a bucket exists.
func (s *S3Fake) BucketExists(bucket string) (bool, error) {
	return s.Backend.BucketExists(bucket)
}

// DeleteBucket deletes a bucket.
func (s *S3Fake) DeleteBucket(bucket string) error {
	return s.Backend.DeleteBucket(bucket)
}

// CreateUser creates a user.
func (s *S3Fake) CreateUser(name string) (*User, error) {
	user := &User{
		SecretKey: "secret-key",
		AccessKey: "access-key",
	}

	s.Users.Store(name, user)

	return user, nil
}

// DeleteUser deletes a user.
func (s *S3Fake) DeleteUser(name string) error {
	s.Users.Delete(name)

	return nil
}

// UserExists checks if a user exists.
func (s *S3Fake) UserExists(name string) (bool, error) {
	_, ok := s.Users.Load(name)

	return ok, nil
}
