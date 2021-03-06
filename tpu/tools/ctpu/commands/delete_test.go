// Copyright 2017 The TensorFlow Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ==============================================================================

package commands

import (
	"testing"

	"context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/tpu/v1alpha1"
)

func testDeleteWorkflow(t *testing.T, libs *testLibs, expectedGCEAction, expectedTPUAction string) {
	t.Helper()
	c := deleteCmd{
		cfg: libs.cfg,
		gce: libs.gce,
		tpu: libs.tpu,
	}
	c.tpuCmd.skipConfirmation = true

	exit := c.Execute(context.Background(), nil)
	if exit != 0 {
		t.Fatalf("Exit code incorrect: %d", exit)
	}

	verifySingleOperation(t, libs.gce.OperationsPerformed, expectedGCEAction)
	verifySingleOperation(t, libs.tpu.OperationsPerformed, expectedTPUAction)
}

func TestDeleteNotExistent(t *testing.T) {
	libs := newTestLibs()
	testDeleteWorkflow(t, libs, "", "")
}

func TestDeleteNotRunning(t *testing.T) {
	libs := newTestLibs()
	libs.gce.instance = &compute.Instance{Status: "STOPPED"}
	libs.tpu.instance = &tpu.Node{State: "CREATING"}
	testDeleteWorkflow(t, libs, "DELETE", "DELETE")
}

func TestDelete(t *testing.T) {
	libs := newTestLibs()
	libs.gce.instance = &compute.Instance{Status: "RUNNING"}
	libs.tpu.instance = &tpu.Node{State: "READY"}
	testDeleteWorkflow(t, libs, "DELETE", "DELETE")
}
