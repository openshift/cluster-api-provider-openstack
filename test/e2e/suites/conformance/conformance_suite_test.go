//go:build e2e
// +build e2e

/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conformance

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-openstack/test/e2e/shared"
)

var e2eCtx *shared.E2EContext

func init() {
	e2eCtx = shared.NewE2EContext()
	shared.CreateDefaultFlags(e2eCtx)
}

func TestConformance(t *testing.T) {
	RegisterFailHandler(Fail)
	ctrl.SetLogger(GinkgoLogr)
	RunSpecs(t, "capo-conformance")
}

var _ = SynchronizedBeforeSuite(func(ctx context.Context) []byte {
	return shared.Node1BeforeSuite(ctx, e2eCtx)
}, func(data []byte) {
	shared.AllNodesBeforeSuite(e2eCtx, data)
})

var _ = SynchronizedAfterSuite(func() {
	shared.AllNodesAfterSuite(e2eCtx)
}, func(ctx context.Context) {
	shared.Node1AfterSuite(ctx, e2eCtx)
})
