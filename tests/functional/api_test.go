//go:build functional

package functional

import (
	"os"
	"testing"

	testsuite "starter-boilerplate/tests/suite"

	"github.com/stretchr/testify/suite"
)

type FunctionalSuite struct {
	testsuite.FunctionalSuite
}

func TestFunctional(t *testing.T) {
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	s := new(FunctionalSuite)
	s.FixtureDir = "tests/functional/testdata/fixtures"
	suite.Run(t, s)
}
