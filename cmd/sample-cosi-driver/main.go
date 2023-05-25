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

package main

import (
	"context"
	"os/signal"
	"syscall"

	"k8s.io/klog/v2"
)

func main() {
	ctx, cancel := setupSignalHandler(context.Background())
	defer cancel()

	if err := cmd.ExecuteContext(ctx); err != nil {
		klog.ErrorS(err, "Exiting on error")
	}
}

// setupSignalHandler creates a context that is canceled when SIGTERM or SIGINT is received.
func setupSignalHandler(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
}
