package main

import (
	"io/ioutil"

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
}

type Config struct {
	Podcasts []Podcast      `yaml:"podcasts"`
	Mixcloud []MixcloudUser `yaml:"mixcloud"`
	Station  Station        `yaml:"station"`
}

func getConfig() Config {
	yamlFile, err := ioutil.ReadFile("conf/config.yaml")

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// func getPodcasts() []Podcast {
// 	podcastMap := getConfig().Podcasts
//
// 	// keys := []int{}
// 	// for k := range podcastMap {
// 	// 	keys = append(keys, k)
// 	// }
// 	return keys
// }
