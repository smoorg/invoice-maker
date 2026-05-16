package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"
)

const APP_CONFIG_FOLDER_NAME = "invoice-maker"

// Returns default config directory appending `APP_CONFIG_FOLDER_NAME` to it.
// If not found it falls back to /etc/invoice-maker
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join("/etc", APP_CONFIG_FOLDER_NAME), nil
	}

	return filepath.Join(configDir, APP_CONFIG_FOLDER_NAME), nil
}

func getConfigFile() string {
	dir, _ := GetConfigDir()
	return filepath.Join(dir, "config.yaml")
}

// Picks invoice directory user configured or ~/.config/invoice-maker/[year]/[month]
func (c *Config) GetInvoiceDirectory() (string, error) {
	if c.Config.InvoiceDirectory == "" {
		return "", errors.New("No defined invoice directory to save the file.")
	}
	path := c.Config.InvoiceDirectory
	invoicePath := filepath.Join(path, time.Now().Format("2006"), time.Now().Format("01"))

	if err := os.MkdirAll(invoicePath, 0744); err != nil {
		return "", err
	}

	return invoicePath, nil
}

func GetInvoicePath(invoiceDirectory string) (string, error) {
	if invoiceDirectory == "" {
		return "", errors.New("No defined invoice directory to save the file.")
	}
	path := invoiceDirectory
	invoicePath := filepath.Join(path, time.Now().Format("2006"), time.Now().Format("01"))

	if err := os.MkdirAll(invoicePath, 0744); err != nil {
		return "", err
	}

	return invoicePath, nil
}

func IsValidInvoiceDirectory(dir string) bool {
	testDir := filepath.Join(dir, fmt.Sprint(uuid.New()))

	if err := os.MkdirAll(testDir, 0744); err != nil {
		return false
	}

	os.Remove(testDir)

	if isAbs := filepath.IsAbs(dir); isAbs == false {
		return false
	}

	return true
}

func GetConfig(cfg *Config) (*Config, error) {
	config := &Config{}
	if cfg != nil {
		config = cfg
	}

	// viper
	appname := "invoice-maker"
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/" + appname)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("issuer", &config.Issuer); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("receivers", &config.Receivers); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("invoices", &config.Invoices); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("invoiceDirectory", &config.Config.InvoiceDirectory); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("font", &config.Config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) WriteConfig() error {
	cfg, marshalErr := yaml.Marshal(c)
	if marshalErr != nil {
		return errors.New("Unable to marshal config")
	}

	err := os.WriteFile(getConfigFile(), cfg, 0744)
	if err != nil {
		return errors.New("Unable to save config file")
	}

	conf, err := GetConfig(nil)
	if err != nil {
		return errors.New("Unable to get config file")
	}
	c = conf
	return nil
}

func (c *Config) WriteReceiver(receiver Company, row int) {
	c.Receivers[row] = receiver
	c.WriteConfig()
}

func (c *Config) WriteInvoiceItem(item InvoiceItem, invoiceRow int, itemRow int) {
	c.Invoices[invoiceRow].Items[itemRow] = item
	c.WriteConfig()

}
