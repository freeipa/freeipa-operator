package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFreeipaOperator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FreeipaOperator Suite")
}
