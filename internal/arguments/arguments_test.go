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
			"--metrics-bind-address",
			":8080",
			"--health-probe-bind-address",
			":8081",
			"--leader-elect",
			"--default-storage",
			"ephemeral",
		}
		// WHEN
		It("initializes with no errors", func() {
			result, err := arguments.NewWithArguments(args)
			// EXPECT
			Expect(err).Should(BeNil())
			Expect(result).ShouldNot(BeNil())
			Expect(result.GetMetricsAddr()).Should(Equal(":8080"))
			Expect(result.GetProbeAddr()).Should(Equal(":8081"))
			Expect(result.GetEnableLeaderElection()).Should(BeTrue())
			Expect(result.GetDefaultStorage()).Should(Equal("ephemeral"))
		})
	})

	Context("NewWithArguments failure", func() {
		// GIVEN
		var args []string = []string{
			"./test",
			"--metrics-bind-address-bad", // Wrong argument
			":8080",
			"--health-probe-bind-address",
			":8081",
			"--leader-elect",
		}
		// WHEN
		It("panics for invalid arguments", func() {
			// EXPECT
			Expect(func() {
				arguments.NewWithArguments(args)
			}).To(Panic())
		})
	})

	Context("NewWithArguments failure", func() {
		// GIVEN
		var args []string = []string{
			"./test",
			"--metrics-bind-address",
			":8080",
			"--health-probe-bind-address-bad", // Wrong argument
			":8081",
			"--leader-elect",
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
			"--metrics-addr",
			":8080",
			"--enable-leader-election",
			"--health-probe-bind-address",
			":8081",
			"--default-storage",
			"unknown", // Wrong argument
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
