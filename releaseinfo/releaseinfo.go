package releaseinfo

// Version is the current version of service.
const Version string = "v1"

// BuildInformation holds current build information.
var BuildInformation string

package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	Database Database `mapstructure:"database"`
	Aws      Aws      `mapstructure:"aws"`
}

type Database struct {
	Name                 string `mapstructure:"name"`
	Host                 string `mapstructure:"host"`
	Pass                 string `mapstructure:"pass"`
	User                 string `mapstructure:"user"`
	Port                 string `mapstructure:"port"`
	ProductCollection    string `mapstructure:"products"`
	ObjectInfoCollection string `mapstructure:"objectinfo"`
}

type Aws struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	S3        []S3   `mapstructure:"S3"`
}

type S3 struct {
	BucketName string `mapstructure:"BucketName"`
	ObjectKey  string `mapstructure:"ObjectKey"`
}

func (c *Config) LoadDatabase() error {

	//db.Name = os.Getenv("DB_NAME")
	//db.Host = os.Getenv("DB_HOST")
	//db.Pass = os.Getenv("DB_PASS")
	//db.User = os.Getenv("DB_USER")
	//db.Port = os.Getenv("DB_PORT")
	//for _, env := range []string{"DB_NAME", "DB_HOST", "DB_PASS", "DB_USER", "DB_PORT"} {
	//	if os.Getenv(env) == "" {
	//		return db, errors.New(env + " is required")
	//	}
	//}

	c.Database.Name = "cimri-go"
	c.Database.Host = "localhost"
	c.Database.Pass = ""
	c.Database.User = ""
	c.Database.Port = "27017"
	c.Database.ObjectInfoCollection = "object-info"
	c.Database.ProductCollection = "products"
	return nil
}

func (c *Config) LoadAws() error {
	//aws.Region = os.Getenv("AWS_REGION")
	//aws.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	//aws.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	//for _, env := range []string{"AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"} {
	//	if os.Getenv(env) == "" {
	//		return aws, errors.New(env + " is required")
	//	}
	//}
	c.Aws.Region = "eu-central-1"
	c.Aws.AccessKey = "AKIA5OUQ2DJ2KKE2JV7P"
	c.Aws.SecretKey = "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"
	return nil
}

func (c *Config) LoadS3Objects() error {
	var s3s []S3
	//if err := viper.UnmarshalKey("aws.S3", &s3s); err != nil {
	//	return nil, err
	//}
	s3s = []S3{
		{
			BucketName: "cimri-casestudy",
			ObjectKey:  "products-1.jsonl",
		},
		{
			BucketName: "cimri-casestudy",
			ObjectKey:  "products-2.jsonl",
		},
		{
			BucketName: "cimri-casestudy",
			ObjectKey:  "products-3.jsonl",
		},
		{
			BucketName: "cimri-casestudy",
			ObjectKey:  "products-4.jsonl",
		},
	}
	c.Aws.S3 = s3s
	return nil
}

// LoadConfig loads configuration from file.
// It sets initial values for database and aws configurations.
func LoadConfig() (*Config, error) {
	var (
		cfg Config
		err error
	)
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigName("s3-objects")
	viper.SetConfigType("yml")
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pwd = strings.TrimRight(pwd, "/cmd")
	viper.AddConfigPath(pwd)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg.Aws); err != nil {
		return nil, err
	}
	if err = cfg.LoadDatabase(); err != nil {
		return nil, err
	}
	if err = cfg.LoadAws(); err != nil {
		return nil, err
	}
	if err = cfg.LoadS3Objects(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
