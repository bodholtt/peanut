package pcfg

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	}
	Images struct {
		RootLocation     string `yaml:"root_location"`
		DeleteImageFiles bool   `yaml:"delete_image_files"`
	}
	Database struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"database_name"`
		Params       string `yaml:"params"`
	}
	Client struct {
		Host string `yaml:"host"`
	}
	Permissions struct {
		DefaultRank       int `yaml:"default_rank"`
		SignUp            int `yaml:"sign_up"`
		CreateUsers       int `yaml:"create_users"`
		DeleteUsers       int `yaml:"delete_users"`
		EditUsers         int `yaml:"edit_users"`
		ViewPosts         int `yaml:"view_posts"`
		CreatePosts       int `yaml:"create_posts"`
		DeleteOwnPosts    int `yaml:"delete_own_posts"`
		DeleteOthersPosts int `yaml:"delete_others_posts"`
		EditOthersPosts   int `yaml:"edit_others_posts"`
		CreateTags        int `yaml:"create_tags"`
		EditTags          int `yaml:"edit_tags"`
		DeleteTags        int `yaml:"delete_tags"`
	}
}

// Cfg - Global config object
var Cfg = Config{}

// InitConfig - initialize the pcfg
func InitConfig(configPath string) error {

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&Cfg); err != nil {
		return err
	}

	return nil
}
