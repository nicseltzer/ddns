package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	APIToken        string
	Timeout         time.Duration
	ZoneID          string
	PublicHostname  string
	PrivateHostname string
	Port            int
}

type config struct {
	ApiToken        string `json:"api_token"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	ZoneID          string `json:"zone_id"`
	PublicHostname  string `json:"public_hostname"`
	PrivateHostname string `json:"private_hostname"`
	Port            int    `json:"service_port"`
}

func NewConfig() *Config {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalln("no config file found at 'config.json'")
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file", err)
		return nil
	}

	var c config
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Println("error unmarshalling body to struct", err)
		return nil
	}

	return &Config{
		APIToken:        c.ApiToken,
		Timeout:         time.Duration(c.TimeoutSeconds) * time.Second,
		ZoneID:          c.ZoneID,
		PublicHostname:  c.PublicHostname,
		PrivateHostname: c.PrivateHostname,
		Port:            c.Port,
	}
}
