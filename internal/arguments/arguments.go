package arguments

import (
	"flag"
	"fmt"
	"os"
)

const (
	ENV_DEFAULT_STORAGE = "DEFAULT_STORAGE"
)

type Arguments struct {
	metricsAddr          string
	probeAddr            string
	enableLeaderElection bool
	defaultStorage       string
}

func NewWithArguments(Args []string) (*Arguments, error) {
	var metricsAddr string
	var probeAddr string
	var enableLeaderElection bool
	var defaultStorage string

	var Flag = flag.NewFlagSet(Args[0], flag.PanicOnError)

	Flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	Flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the metric endpoint binds to.")
	Flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	Flag.StringVar(&defaultStorage,
		"default-storage",
		os.Getenv(ENV_DEFAULT_STORAGE),
		"Default storage. An empty string here make sure the persistent volume"+
			"claim section is not empty. "+
			"It can be 'ephemeral' or 'hostpath'. The default value is retrieved from "+
			"the DEFAULT_STORAGE environment variable.")
	Flag.Usage = func() {
		fmt.Printf("Controller usage:\n")
		fmt.Printf("%s --metrics-bind-address :8080 --health-probe-bind-address :8081 --leader-elect --default-storage {ephemeral,hostPath}\n", os.Args[0])
	}

	Flag.Parse(Args[1:])

	if defaultStorage != "" && defaultStorage != "ephemeral" && defaultStorage != "hostpath" {
		Flag.Usage()
		fmt.Printf("ERROR:'--default-storage' must be empty, 'ephemeral' or 'hostpath' values\n")
		fmt.Printf("ERROR:When not specified it defaults to DEFAULT_STORAGE environment variable\n")
		panic(func() {})
	}

	return &Arguments{
		metricsAddr:          metricsAddr,
		probeAddr:            probeAddr,
		enableLeaderElection: enableLeaderElection,
		defaultStorage:       defaultStorage,
	}, nil
}

func New() (*Arguments, error) {
	return NewWithArguments(os.Args)
}

func (a *Arguments) GetMetricsAddr() string {
	return a.metricsAddr
}

func (a *Arguments) GetProbeAddr() string {
	return a.probeAddr
}

func (a *Arguments) GetEnableLeaderElection() bool {
	return a.enableLeaderElection
}

func (a *Arguments) GetDefaultStorage() string {
	return a.defaultStorage
}
