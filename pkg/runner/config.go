package runner

import (
	"Yi/pkg/logging"
	"Yi/pkg/utils"
	"bytes"
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/**
  @author: yhy
  @since: 2022/12/12
  @desc: //TODO
**/

type QLFile struct {
	PythonQL []string `mapstructure:"python_ql"`
}

var QLFiles *QLFile

// HotConf uses viper to hot load configurations.
func HotConf() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logging.Logger.Fatalf("cmd.HotConf, fail to get current path: %v", err)
	}
	// Configuration file path Current folder + config.yaml
	configFile := path.Join(dir, ConfigFileName)
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	// watch monitors configuration file changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Callback functions that are called after a change in the configuration file
		logging.Logger.Infoln("config file changed: ", e.Name)
		oldQls := QLFiles
		ReadYamlConfig(configFile)
		newQls := QLFiles
		// When the rule is updated, get the item from the database and run it through with the new rule
		NewRules(oldQls, newQls)
	})
}

// Init Load Configuration
func Init() {
	//Configuration file path Current folder + config.yaml
	configFile := path.Join(Pwd, ConfigFileName)

	// Detecting the existence of a configuration file
	if !utils.Exists(configFile) {
		WriteYamlConfig(configFile)
		logging.Logger.Infof("%s not find, Generate profile.", configFile)
	} else {
		logging.Logger.Infoln("Load profile ", configFile)
	}
	ReadYamlConfig(configFile)

}

func ReadYamlConfig(configFile string) {
	// load config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to read 'config.yaml': %+v", err)
	}
	err = viper.Unmarshal(&QLFiles)
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to parse 'config.yaml', check format: %v", err)
	}
}

func WriteYamlConfig(configFile string) {
	// Generate default config
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultYamlByte))
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to read default config bytes: %v", err)
	}
	// write a file
	err = viper.SafeWriteConfigAs(configFile)
	if err != nil {
		logging.Logger.Fatalf("setting.Setup, fail to write 'config.yaml': %v", err)
	}
}
