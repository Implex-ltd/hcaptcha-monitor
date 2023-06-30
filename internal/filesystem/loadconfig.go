package filesystem

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func LoadConfigFile() (*ConfigStruct, error) {
	fmt.Printf(`  _  _    __  __          _ _           
 | || |__|  \/  |___ _ _ (_) |_ ___ _ _ 
 | __ / _| |\/| / _ \ ' \| |  _/ _ \ '_|
 |_||_\__|_|  |_\___/_||_|_|\__\___/_|
\n`)

	var config ConfigStruct
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
