package common

import (
	"fmt"
	operator "github.com/alexellis/k3sup/pkg/operator"
)

func ExecuteLocalCommand(command string) error {
	operator := operator.ExecOperator{}

	fmt.Printf("Executing: %s\n", command)

	res, err := operator.Execute(command)
	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		if len(res.StdErr) > 0 {
			fmt.Printf("stderr: %q", res.StdErr)
		}
	}

	if len(res.StdOut) > 0 {
		fmt.Printf("stdout: %q", res.StdOut)
	}

	return nil
}
