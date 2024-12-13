package e2esuite

type MulberryConfig struct {
	JackalConfig struct {
		RPC      string `yaml:"rpc"`
		GRPC     string `yaml:"grpc"`
		SeedFile string `yaml:"seed_file"`
		Contract string `yaml:"contract"`
	} `yaml:"jackal_config"`
	NetworksConfig []struct {
		Name     string `yaml:"name"`
		RPC      string `yaml:"rpc"`
		Contract string `yaml:"contract"`
		ChainID  int    `yaml:"chain_id"`
		Finality int    `yaml:"finality"`
	} `yaml:"networks_config"`
}
