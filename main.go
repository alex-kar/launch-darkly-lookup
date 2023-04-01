package main

import (
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"strconv"
	"gopkg.in/yaml.v3"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v6"
)

type config struct {
	SdkKey string `yaml:"sdkKey"`
	Environments []string `yaml:"environment"`
	Companies []string `yaml:"company"`
	Flags []string `yaml:"flag"`
}

func main() {
	configPath := configPath(os.Args)
	config := loadConfig(configPath) 
	ldClient := initLDClient(config.SdkKey)
	for _, env := range config.Environments {
			fmt.Printf("===== %v\n", env)
		for _, company := range config.Companies {
			fmt.Printf("== %v\n", company)
			for _, flag := range config.Flags {
				value := evaluateFlag(ldClient, company, env, flag)
				fmt.Printf("%v: %v\n", flag, value)
			}

		}
	}
}

func configPath(args []string) string {
	if (len(args) != 2) {
		fmt.Println("configuration file is missing")
		os.Exit(1)
	}
	return args[1]
}

func loadConfig(path string) *config {
	buf, err := ioutil.ReadFile(path)
	if (err != nil) {
		fmt.Println("File not found: ", path)
		os.Exit(1)
	}
	result := &config{}
	err = yaml.Unmarshal(buf, result)
	if (err != nil) {
		fmt.Println("Failed to parse configuration file")
		os.Exit(1)
	}
	return result
}

func initLDClient(sdkKey string) *ld.LDClient {
	ldClient, _ := ld.MakeClient(sdkKey, 5*time.Second)
	if (!ldClient.Initialized()) {
		fmt.Println("LD failed to connect: ")
		os.Exit(1)
	}
	return ldClient
}

func evaluateFlag(ldClient *ld.LDClient, company string, env string, flag string) string {
	context := ldcontext.NewBuilder(company).Build()
	flagValue, err := ldClient.BoolVariation(flag, context, false)
	if (err != nil) {
		return "unknown"
	}
	return strconv.FormatBool(flagValue)
}
