package utils

import (
	"errors"
	"strings"
)

var ErrInvalidEnvironment = errors.New("invalid environment specified, must be production, sandbox, development/staging, lab, or local")

type Environment string

const (
	// DevelopmentEnv ...
	DevelopmentEnv Environment = "development"
	// LabEnv ...
	LabEnv Environment = "lab"
	// SandboxEnv ...
	SanboxEnv Environment = "sandbox"
	// ProductionEnv ...
	ProductionEnv Environment = "production"
	// StagingEnv ...
	StagingEnv Environment = "staging"
	// LocalEnv ...
	LocalEnv Environment = "local"
)

func (e Environment) String() string {
	return string(e)
}

func ParseEnvironment(input string) (Environment, error) {
	var runEnv Environment

	switch strings.ToLower(input) {
	case DevelopmentEnv.String(), "dev":
		runEnv = DevelopmentEnv
	case LabEnv.String(), "lab":
		runEnv = LabEnv
	case SanboxEnv.String(), "sandbox":
		runEnv = SanboxEnv
	case ProductionEnv.String(), "prod":
		runEnv = ProductionEnv
	case StagingEnv.String(), "stage":
		runEnv = StagingEnv
	case LocalEnv.String(), "local":
		runEnv = LocalEnv
	default:
		return "", ErrInvalidEnvironment
	}

	return runEnv, nil
}
