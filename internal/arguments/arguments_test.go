package arguments_test

import (
	arguments "github.com/freeipa/freeipa-operator/internal/arguments"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Arguments package", func() {
	Context("NewWithArguments success", func() {
		// GIVEN
		var args []string = []string{
			"./test",
			"--metrics-addr",
			":8081",
			"--enable-leader-election",
			"--default-storage",
			"ephemeral",
		}
		// WHEN
		It("initializes with no errors", func() {
			result, err := arguments.NewWithArguments(args)
			// EXPECT
			Expect(err).Should(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.GetMetricsAddr()).Should(Equal(":8081"))
			Expect(result.GetEnableLeaderElection()).Should(BeTrue())
			Expect(result.GetDefaultStorage()).Should(Equal("ephemeral"))
		})
	})

	Context("NewWithArguments failure", func() {
		// GIVEN
		var args []string = []string{
			"./test",
			"--metrics-addr2", // Wrong argument
			":8080",
			"--enable-leader-election",
		}
		// WHEN
		It("panics for invalid arguments", func() {
			// EXPECT
			Expect(func() {
				arguments.NewWithArguments(args)
			}).To(Panic())
		})
	})

	Context("NewWithArguments failure for defaultStorage", func() {
		// GIVEN
		var args []string = []string{
			"./test",
			"--metrics-addr", // Wrong argument
			":8080",
			"--enable-leader-election",
			"--default-storage",
			"unknown",
		}
		// WHEN
		It("panics for invalid arguments", func() {
			// EXPECT
			Expect(func() {
				arguments.NewWithArguments(args)
			}).To(Panic())
		})

	})
})
