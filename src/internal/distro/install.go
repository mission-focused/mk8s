package distro

import (
	"fmt"

	"github.com/brandtkeller/mk8s/src/types"
)

func Install(config types.MultiConfig) error {

	switch config.Distro {
	case "rke2":
		err := installMultiRKE2(config)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Distro %s not supported", config.Distro)
	}

	return nil

}
