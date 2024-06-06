package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	yaml "gopkg.in/yaml.v3"
)

const APP_CONFIG_FOLDER_NAME = "invoice-maker"

// Returns default config directory appending `APP_CONFIG_FOLDER_NAME` to it.
// If not found it falls back to /etc/invoice-maker
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		msg := fmt.Sprintf(
			"user config directory at '%s' missing, falling back to '%s' default package config",
			configDir,
			"/etc/"+APP_CONFIG_FOLDER_NAME,
		)

		return filepath.Join("/etc", APP_CONFIG_FOLDER_NAME), errors.New(msg)
	}

	return filepath.Join(configDir, APP_CONFIG_FOLDER_NAME), nil
}

func getConfigFile() string {
	dir, _ := GetConfigDir()
	return filepath.Join(dir, "config.yaml")
}

// Picks invoice directory user configured or ~/.config/invoice-maker/[year]/[month]
func (c *Config) GetInvoiceDirectory() (string, error) {
	if c.InvoiceDirectory == "" {
		return "", errors.New("No defined invoice directory to save the file.")
	}
	path := c.InvoiceDirectory
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

func GetConfig() (*Config, error) {
	cfg := &Config{}
	dir, _ := GetConfigDir()
	os.MkdirAll(dir, 0744)

	configYaml, err := os.ReadFile(getConfigFile())
	if err != nil {
		file, fileCreateErr := os.Create(getConfigFile())
		defer file.Close()
		if fileCreateErr != nil {
			return nil, errors.New("unable to create yaml config file")
		}

		_, err = file.Read(configYaml)
	}

	marshalErr := yaml.Unmarshal(configYaml, &cfg)
	if marshalErr != nil {
		return nil, errors.New("Unable to parse config file")
	}

	return cfg, nil
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

	if conf, err := GetConfig(); err == nil {
		c = conf
	}
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
