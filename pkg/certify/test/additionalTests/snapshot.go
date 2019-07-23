package additionalTests

import (
	. "github.com/onsi/ginkgo"
	"k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/testpatterns"
	"k8s.io/kubernetes/test/e2e/storage/testsuites"
)

type extraSnapshotTestSuite struct {
	tsInfo TestSuiteInfo
}

var _ TestSuite = &extraSnapshotTestSuite{}

func InitVolumesTestSuite() TestSuite {
	return &extraSnapshotTestSuite{
		tsInfo: TestSuiteInfo{
			name: "provisioning",
			testPatterns: []testpatterns.TestPattern{
				testpatterns.DefaultFsDynamicPV,
			},
		},
	}
}

func (t *extraSnapshotTestSuite) getTestSuiteInfo() TestSuiteInfo {
	return t.tsInfo;
}

func (p *extraSnapshotTestSuite) defineTests(driver testsuites.TestDriver, pattern testpatterns.TestPattern) {
	type local struct {
		config      *testsuites.PerTestConfig
		testCleanup func()

		testCase *testsuites.StorageClassTest
		cs       clientset.Interface
		pvc      *v1.PersistentVolumeClaim
		sc       *storage.StorageClass
	}
	var (
		dInfo   = driver.GetDriverInfo()
		dDriver testsuites.DynamicPVTestDriver
		l       local
	)

	BeforeEach(func() {
		// Check preconditions.
		if pattern.VolType != testpatterns.DynamicPV {
			framework.Skipf("Suite %q does not support %v", p.tsInfo.name, pattern.VolType)
		}
		ok := false
		dDriver, ok = driver.(testsuites.DynamicPVTestDriver)
		if !ok {
			framework.Skipf("Driver %s doesn't support %v -- skipping", dInfo.Name, pattern.VolType)
		}
	})

	// This intentionally comes after checking the preconditions because it
	// registers its own BeforeEach which creates the namespace. Beware that it
	// also registers an AfterEach which renders f unusable. Any code using
	// f must run inside an It or Context callback.
	f := framework.NewDefaultFramework("provisioning")

	init := func() {
		l = local{}

		// Now do the more expensive test initialization.
		l.config, l.testCleanup = driver.PrepareTest(f)
		l.cs = l.config.Framework.ClientSet
		claimSize := dDriver.GetClaimSize()
		l.sc = dDriver.GetDynamicProvisionStorageClass(l.config, "")
		if l.sc == nil {
			framework.Skipf("Driver %q does not define Dynamic Provision StorageClass - skipping", dInfo.Name)
		}
		l.pvc = getClaim(claimSize, l.config.Framework.Namespace.Name)
		l.pvc.Spec.StorageClassName = &l.sc.Name
		framework.Logf("In creating storage class object and pvc object for driver - sc: %v, pvc: %v", l.sc, l.pvc)
		l.testCase = &testsuites.StorageClassTest{
			Client:       l.config.Framework.ClientSet,
			Claim:        l.pvc,
			Class:        l.sc,
			ClaimSize:    claimSize,
			ExpectedSize: claimSize,
		}
	}

	cleanup := func() {
		if l.testCleanup != nil {
			l.testCleanup()
			l.testCleanup = nil
		}
	}

	It("should provision storage with defaults", func() {
		init()
		defer cleanup()

		l.testCase.TestDynamicProvisioning()
	})
}
