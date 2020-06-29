// Helpers for very simple configuration and key loading.
// Will use json files for config.
package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Simple file-based config and secret store.
//
// Implements a ConfStore interface with List/Get/Set interface.
// TODO: Watch interface - using messages/pubsub !!!
//
// TODO: switch to yaml, support K8S style
//
type Conf struct {
	// Base directory. If not set, no config will be saved and read will fail.
	base string
	// Conf is configured from Android side with the config (settings)
	// ssid, pass, vpn_ext
	Conf map[string]string `json:"Conf,omitempty"`
}

// Returns a config stsore.
// Implements auth.ConfStore interface
func NewConf(base string) *Conf {
	// TODO: https for remote - possibly using local creds and K8S style or XDS
	return &Conf{base: base, Conf: map[string]string{}}
}

func (h2 *Conf) List(name string, tp string) ([]string, error) {
	return nil, nil
}

func Get(h2 *Conf, name string, to interface{}) error {
	raw, err := h2.Get("gate.json")
	if err != nil {
		log.Println("gate.json:", err)
		raw = []byte("{}")
		//return nil
	}
	if err := json.Unmarshal(raw, to); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Secrets - pem, acl
// From config dir, fallback to .sshterraform, .lego and /etc/certs
//
// "name" may be a hostname
func (h2 *Conf) Get(name string) ([]byte, error) {
	envName := strings.ReplaceAll(name, ".", "_")
	envName = strings.ReplaceAll(envName, "/", "_")
	envd := os.Getenv(envName)
	if envd != "" {
		return []byte(envd), nil
	}
	l := h2.base + name

	if _, err := os.Stat(l); err == nil { // || !os.IsNotExist(err)
		res, err := ioutil.ReadFile(l)
		if err == nil {
			return res, nil
		}
	}
	if _, err := os.Stat(l + ".json"); err == nil { // || !os.IsNotExist(err)
		res, err := ioutil.ReadFile(l + ".json")
		if err == nil {
			return res, nil
		}
	}

	// name may be a hostname - use it to load ACME certificate for the host.

	return nil, nil
}

func (h2 *Conf) Set(conf string, data []byte) error {
	if h2 == nil || h2.base == "" {
		return nil
	}
	err := ioutil.WriteFile(h2.base+conf, data, 0700)
	if err != nil {
		log.Println("Error saving ", err, h2.base, conf)
	}
	return err
}