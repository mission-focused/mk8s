package types

type MultiConfig struct {
	Distro    string              `yaml:"distro" json:"distro"`
	Version   string              `yaml:"version" json:"version"`
	Arch      string              `yaml:"arch" json:"arch"`
	Airgap    bool                `yaml:"airgap" json:"airgap"`
	Nodes     []NodeConfig        `yaml:"nodes" json:"nodes"`
	Artifacts map[string]Artifact `yaml:"artifacts" json:"artifacts"`
}

type NodeConfig struct {
	Address    string `yaml:"address" json:"address"`
	User       string `yaml:"user" json:"user"`
	Role       string `yaml:"role" json:"role"`
	Primary    bool   `yaml:"primary" json:"primary"`
	Local      bool   `yaml:"local" json:"local"`
	SshKeyPath string `yaml:"ssh_key_path" json:"ssh_key_path" `
	Config     string `yaml:"config" json:"config"`
}

type Artifact struct {
	Name     string `yaml:"name" json:"name"`
	URL      string `yaml:"url" json:"url"`
	Checksum string `yaml:"checksum" json:"checksum"`
}
