/*
Copyright 2018 The Kubernetes Authors.

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

package testsuites

import (
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/testpatterns"
)

// TestDriver represents an interface for a driver to be tested in TestSuite
type TestDriver interface {
	// GetDriverInfo returns DriverInfo for the TestDriver
	GetDriverInfo() *DriverInfo
	// CreateDriver creates all driver resources that is required for TestDriver method
	// except CreateVolume. May be called more than once and should only do something on
	// the first call.
	CreateDriver(config *TestConfig)
	// CreateDriver cleanup all the resources that is created in CreateDriver. There is
	// no guarantee that CreateDriver succeeded or even was called at all, so the test driver
	// has to track resources.
	CleanupDriver()
}

// FilterTestDriver is an optional interface that drivers can
// implement to filter out unsuitable tests while tests get defined.
type FilterTestDriver interface {
	// IsTestSupported returns true if the Testpattern is
	// suitable to test with the TestDriver. This will be called
	// already while defining tests and unsupported tests will not even
	// be added to the test suite.
	IsTestSupported(testpatterns.TestPattern) bool
}

// TestVolume is the result of PreprovisionedVolumeTestDriver.CreateVolume.
// The only common functionality is to delete it. Individual driver interfaces
// have additional methods that work with volumes created by them.
type TestVolume interface {
	DeleteVolume()
}

// PreprovisionedVolumeTestDriver represents an interface for a TestDriver that has pre-provisioned volume
type PreprovisionedVolumeTestDriver interface {
	TestDriver
	// CreateVolume creates a pre-provisioned volume of the desired volume type.
	CreateVolume(config *TestConfig, volumeType testpatterns.TestVolType) TestVolume
}

// InlineVolumeTestDriver represents an interface for a TestDriver that supports InlineVolume
type InlineVolumeTestDriver interface {
	PreprovisionedVolumeTestDriver

	// GetVolumeSource returns a volumeSource for inline volume.
	// It will set readOnly and fsType to the volumeSource, if TestDriver supports both of them.
	// It will return nil, if the TestDriver doesn't support either of the parameters.
	GetVolumeSource(config *TestConfig, readOnly bool, fsType string, testVolume TestVolume) *v1.VolumeSource
}

// PreprovisionedPVTestDriver represents an interface for a TestDriver that supports PreprovisionedPV
type PreprovisionedPVTestDriver interface {
	PreprovisionedVolumeTestDriver
	// GetPersistentVolumeSource returns a PersistentVolumeSource for pre-provisioned Persistent Volume.
	// It will set readOnly and fsType to the PersistentVolumeSource, if TestDriver supports both of them.
	// It will return nil, if the TestDriver doesn't support either of the parameters.
	GetPersistentVolumeSource(config *TestConfig, readOnly bool, fsType string, testVolume TestVolume) *v1.PersistentVolumeSource
}

// DynamicPVTestDriver represents an interface for a TestDriver that supports DynamicPV
type DynamicPVTestDriver interface {
	TestDriver
	// GetDynamicProvisionStorageClass returns a StorageClass dynamic provision Persistent Volume.
	// It will set fsType to the StorageClass, if TestDriver supports it.
	// It will return nil, if the TestDriver doesn't support it.
	GetDynamicProvisionStorageClass(config *TestConfig, fsType string) *storagev1.StorageClass

	// GetClaimSize returns the size of the volume that is to be provisioned ("5Gi", "1Mi").
	// The size must be chosen so that the resulting volume is large enough for all
	// enabled tests and within the range supported by the underlying storage.
	GetClaimSize() string
}

// Capability represents a feature that a volume plugin supports.
// Unless noted otherwise, capabilities are assumed to be not
// supported when not set explicitly.
type Capability string

const (
	CapPersistence Capability = "persistence" // data is persisted across pod restarts
	CapBlock       Capability = "block"       // raw block mode
	CapFsGroup     Capability = "fsGroup"     // volume ownership via fsGroup
	CapExec        Capability = "exec"        // exec a file in the volume
	CapMultiPODs   Capability = "multipods"   // multiple pods on a node can use the same volume concurrently; enabled by default
)

// DriverInfo represents static information about a TestDriver.
type DriverInfo struct {
	Name       string // Name of the driver
	FeatureTag string // FeatureTag for the driver

	MaxFileSize          int64               // Max file size to be tested for this driver
	SupportedFsType      sets.String         // Map of string for supported fs type
	SupportedMountOption sets.String         // Map of string for supported mount option
	RequiredMountOption  sets.String         // Map of string for required mount option (Optional)
	Capabilities         map[Capability]bool // Map that represents plugin capabilities

	Config TestConfig // Test configuration for the current test.
}

// TestConfig represents parameters that control test execution.
// One instance gets allocated for each test and is then passed
// via pointer to functions involved in the test. These functions,
// in particular CreateDriver, can update fields if needed for
// the remaining test execution.
type TestConfig struct {
	// Some short word that gets inserted into dynamically
	// generated entities (pods, paths) as first part of the name
	// to make debugging easier. Can be the same for different
	// tests inside the test suite.
	Prefix string

	// The framework instance allocated for the current test.
	Framework *framework.Framework

	// If non-empty, then pods using a volume will be scheduled
	// onto the node with this name. Otherwise Kubernetes will
	// pick a node.
	ClientNodeName string

	// Some tests also support scheduling pods onto nodes with
	// these label/value pairs. As not all tests use this field,
	// a driver that absolutely needs the pods on a specific
	// node must use ClientNodeName.
	ClientNodeSelector map[string]string

	// Some test drivers initialize a storage server. This is
	// the configuration that then has to be used to run tests.
	// The values above are ignored for such tests.
	ServerConfig *framework.VolumeTestConfig
}