package main

import (
	"os"
	"fmt"
	"github.com/logmatic/beats-forwarder/forwarder"
	cfg "github.com/logmatic/beats-forwarder/config"
)

var config = cfg.Config{}


func main() {

	// read the configuration
	err := cfg.Read(&config, "config-logmatic.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// launch the forwarder
	err = forwarder.Run(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}



}
