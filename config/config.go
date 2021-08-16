package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	FinnhubConf  FinnhubConfigParam `yaml:"finnhub-conf"`
	DatabaseConf DBConfigParam      `yaml:"pg-database-conf"`
}

type FinnhubConfigParam struct {
	FinnhubApiKey string `yaml:"finnhub-api-key"`
}

type DBConfigParam struct {
	PgHost     string `yaml:"pg-host"`
	PgPort     string `yaml:"pg-port"`
	PgSchema   string `yaml:"pg-schema"`
	PgDatabase string `yaml:"pg-database"`
	PgUser     string `yaml:"pg-user"`
	PgPassword string `yaml:"pg-password"`
}

func (c *Conf) GetConf() *Conf {
	ymlFile := readFileConfig()
	err := yaml.Unmarshal(ymlFile, c)
	if err != nil {
		log.Fatalf("Unmarshall: %v", err)
	}
	return c
}

func readFileConfig() []byte {
	ymlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Printf("yamlFile.Get error:   #%v ", err)
		return nil
	}
	return ymlFile
}
