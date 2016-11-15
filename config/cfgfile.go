package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
)

// Command line flags
var configfile *string

func init() {

	configfile = flag.String("c", "etc/config.yml", "Configuration file")
	flag.Bool("d", false, "Enable debug logs")

}

func Read(out interface{}, path string) error {

	// default config
	err := internalRead(out, "etc/config.yml")
	if (err!=nil){
		return err
	}
	// user-override
	if (*configfile != flag.Lookup("c").DefValue) {
		err := internalRead(out, "")
		if (err!=nil){
			return err
		}
	}
	return nil
}

// Read reads the configuration from a yaml file into the given interface structure.
// In case path is not set this method reads from the default configuration file for the beat.
func internalRead(out interface{}, path string) error {

	if path == "" {
		path = *configfile
	}

	filecontent, err := ioutil.ReadFile(path)

	if err != nil {
		return fmt.Errorf("Failed to read %s: %v. Exiting.", path, err)
	}
	filecontent = expandEnv(filecontent)

	config, err := yaml.NewConfig(filecontent, ucfg.PathSep("."))
	if err != nil {
		return fmt.Errorf("YAML config parsing failed on %s: %v. Exiting.", path, err)
	}

	err = config.Unpack(out, ucfg.PathSep("."))
	if err != nil {
		return fmt.Errorf("Failed to apply config %s: %v. Exiting. ", path, err)
	}
	return nil
}


// expandEnv replaces ${var} or $var in config according to the values of the
// current environment variables. The replacement is case-sensitive. References
// to undefined variables are replaced by the empty string. A default value
// can be given by using the form ${var:default value}.
func expandEnv(config []byte) []byte {
	return []byte(os.Expand(string(config), func(key string) string {
		keyAndDefault := strings.SplitN(key, ":", 2)
		key = keyAndDefault[0]

		v := os.Getenv(key)
		if v == "" && len(keyAndDefault) == 2 {
			// Set value to the default.
			v = keyAndDefault[1]
			if (strings.HasPrefix(v, "$")) {
				v = os.Getenv(v[1:])
			}
		}

		return v
	}))
}


