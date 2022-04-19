package generate

// ExecutionEnv : Struct Data object to be applied to execution docker-compose script environment template
type ExecutionEnv struct {
	DataDir string
}

// ConsensusEnv : Struct Data object to be applied to consensus docker-compose script environment template
type ConsensusEnv struct {
	ExecutionNodeURL string
	DataDir          string
}

// ValidatorEnv : Struct Data object to be applied to validator docker-compose script environment template
type ValidatorEnv struct {
	ConsensusNodeURL    string
	ExecutionEngineName string
	DataDir             string
	KeystoreDir         string
}
