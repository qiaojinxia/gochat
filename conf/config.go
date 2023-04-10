package conf

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

var GlobalConfig map[string]map[string]string

var NoConfigError = errors.New("no module config error")

func GetConfig(module string, configName string) string {
	if moduleConfig, exists := GlobalConfig[module]; exists {
		if config, exists := moduleConfig[configName]; exists {
			return config
		}
	}
	panic(NoConfigError)
}

func GetConfigInt(module string, configName string) int {
	if moduleConfig, exists := GlobalConfig[module]; exists {
		if config, exists := moduleConfig[configName]; exists {
			iConfig, err := strconv.Atoi(config)
			if err != nil {
				panic(err)
			}
			return iConfig
		}
	}
	panic(NoConfigError)
}

func InitGlobalConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var section string
	config := make(map[string]map[string]string)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore blank lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Check for section header
		if line[0] == '[' && line[len(line)-1] == ']' {
			section = line[1 : len(line)-1]
			config[section] = make(map[string]string)
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Println("Invalid line:", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[section][key] = value
	}
	GlobalConfig = config
	log.Println("read Config success!")
	for moduleName, v := range config {
		log.Printf("[%s]\n", moduleName)
		for configKey, configValue := range v {
			log.Printf("	 [%s = %s]\n", configKey, configValue)
		}
		log.Println()
	}
	// Print the config

}
