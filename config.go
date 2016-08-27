package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Podcast struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

type MixcloudUser struct {
	Username string `yaml:"username"`
	Title    string `yaml:"title"`
}

type Station struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Id       int    `yaml:"id"`
	Playlist string `yaml:"playlist"`
}

type Config struct {
	Podcasts []Podcast      `yaml:"podcasts"`
	Mixcloud []MixcloudUser `yaml:"mixcloud"`
	Station  Station        `yaml:"station"`
}

func getConfig() Config {

	if !FileExists("conf/config.yaml") {
		log.Fatal("ERROR: valid conf/config.yaml is required")
	}

	yamlFile, err := ioutil.ReadFile("conf/config.yaml")

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}