package main

import (
	"strings"
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"strconv"
	"gopkg.in/yaml.v3"
	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v6"
    _ "github.com/launchdarkly/go-server-sdk/v6/ldfiledata"
)

type config struct {
	SdkKey string `yaml:"sdkKey"`
	Environments []string `yaml:"environments"`
	Flags []flag `yaml:"flags"`
}

type flag struct {
	Key string `yaml:"key"`
	Desc string `yaml:"desc"`
}

func main() {
	configPath := configPath(os.Args)
	config := loadConfig(configPath) 
	ldClient := initLDClient(config.SdkKey)
	companies := companyIds(os.Args)
	for _, company := range companies {
		fmt.Printf("==== Company ID:  %v\n", company)
		for _, env := range config.Environments {
			fmt.Printf("== env: %v\n", env)
			for _, flag := range config.Flags {
				value := evaluateFlag(ldClient, company, env, flag.Key)
				fmt.Println(flagName(flag), value)
			}
		}

	}
}

func configPath(args []string) string {
	if (len(args) < 2) {
		fmt.Println("configuration file is missing")
		os.Exit(1)
	}
	if (len(args) < 3) {
		fmt.Println("At least one copmany id should be provided")
		os.Exit(1)
	}
	return args[1]
}

func companyIds(args []string) []string {
	if (len(args) < 3) {
		fmt.Println("At least one copmany id should be provided")
		os.Exit(1)
	}
	return os.Args[2:]
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
	var config ld.Config
	//config.DataSource = ldfiledata.DataSource().FilePaths("example.json")
	ldClient, _ := ld.MakeCustomClient(sdkKey, config, 5*time.Second)
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

func flagName(flag flag) string {
	var sb strings.Builder
	sb.WriteString(flag.Key)
	if (flag.Desc != "") {
		sb.WriteString("(")
		sb.WriteString(flag.Desc)
		sb.WriteString("): ")
	}
	return sb.String()
}
