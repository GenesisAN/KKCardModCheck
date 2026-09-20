package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	UnusedModOutputPath  string `json:"unused_mod_output_path"`
	GamePath             string `json:"game_path"`
	CardPath             string `json:"card_path"`
	ModPath              string `json:"mod_path"`
	InitConfig           bool   `json:"init_config"`
	ModInfoBuild         bool   `json:"info_build"`
	ModInfoBuildTime     string `json:"info_build_time"`
	HideModeInfoBuild    bool   `json:"hide_mode_info_build"`
	AutoDetectModChanges bool   `json:"auto_detect_mod_changes"`
	ModInfoPath          string `json:"mod_info_path"`
	configPath           string `json:"-"`
}

var ModInfoChage bool

const AppName = "KKCardModCheck"
const DataDirEnvironment = "KKCARDMODCHECK_DATA_DIR"

var Instance Config

func Load() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(GetDataDir())
	viper.SetDefault("game_path", "")
	viper.SetDefault("card_path", "")
	viper.SetDefault("mod_path", "")
	viper.SetDefault("unused_mod_output_path", "")
	viper.SetDefault("init_config", false)
	viper.SetDefault("info_build", false)
	viper.SetDefault("info_build_time", "")
	viper.SetDefault("auto_detect_mod_changes", false)
	viper.SetDefault("hide_mode_info_build", false)
	viper.SetDefault("mod_info_path", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Creating a new one with default values.")
			Save()
		} else {
			log.Printf("Error reading config file: %v. Using default values.", err)
		}
	}
	if err := viper.Unmarshal(&Instance); err != nil {
		log.Printf("Error unmarshaling config file: %v", err)
	}
	Instance.configPath = viper.ConfigFileUsed()
	Instance.GamePath = viper.GetString("game_path")
	Instance.CardPath = viper.GetString("card_path")
	Instance.ModPath = viper.GetString("mod_path")
	Instance.UnusedModOutputPath = viper.GetString("unused_mod_output_path")
	Instance.InitConfig = viper.GetBool("init_config")
	Instance.ModInfoBuild = viper.GetBool("info_build")
	Instance.ModInfoBuildTime = viper.GetString("info_build_time")
	Instance.AutoDetectModChanges = viper.GetBool("auto_detect_mod_changes")
	Instance.HideModeInfoBuild = viper.GetBool("hide_mode_info_build")
	Instance.ModInfoPath = viper.GetString("mod_info_path")
}

func ResetToDefault() {
	Instance = Config{}
	Save()
}

func Save() {
	viper.Set("game_path", Instance.GamePath)
	viper.Set("card_path", Instance.CardPath)
	viper.Set("mod_path", Instance.ModPath)
	viper.Set("unused_mod_output_path", Instance.UnusedModOutputPath)
	viper.Set("init_config", Instance.InitConfig)
	viper.Set("info_build", Instance.ModInfoBuild)
	viper.Set("info_build_time", Instance.ModInfoBuildTime)
	viper.Set("auto_detect_mod_changes", Instance.AutoDetectModChanges)
	viper.Set("hide_mode_info_build", Instance.HideModeInfoBuild)
	viper.Set("mod_info_path", Instance.ModInfoPath)
	if err := viper.WriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfgPath := filepath.Join(GetDataDir(), "config.json")
			if err := viper.WriteConfigAs(cfgPath); err != nil {
				log.Printf("Error creating config file: %v", err)
			}
		} else {
			log.Printf("Error writing config file: %v", err)
		}
	}
}

func GetConfigPath() string {
	return filepath.Join(GetDataDir(), "config.json")
}

func GetDataDir() string {
	if override := strings.TrimSpace(os.Getenv(DataDirEnvironment)); override != "" {
		if absolute, err := filepath.Abs(override); err == nil {
			override = absolute
		}
		if err := os.MkdirAll(override, 0755); err != nil {
			log.Printf("Unable to create overridden data dir %s: %v", override, err)
		}
		return override
	}

	exePath, err := os.Executable()
	if err != nil {
		fallback := "./" + AppName + "_Data"
		_ = os.MkdirAll(fallback, 0755)
		return fallback
	}
	dataDir := filepath.Join(filepath.Dir(exePath), AppName+"_Data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Unable to create data dir %s: %v", dataDir, err)
	}
	return dataDir
}
