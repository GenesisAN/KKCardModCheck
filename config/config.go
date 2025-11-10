package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	GamePath             string `json:"game_path"`
	CardPath             string `json:"card_path"`
	ModPath              string `json:"mod_path"`
	InitConfig           bool   `json:"init_config"`
	ModInfoBuild         bool   `json:"info_build"`              // 是否构筑了Mod信息
	ModInfoBuildTime     string `json:"info_build_time"`         // Mod信息构筑时间
	HideModeInfoBuild    bool   `json:"hide_mode_info_build"`    // 隐藏Mod信息构筑页面
	AutoDetectModChanges bool   `json:"auto_detect_mod_changes"` // 启动时自动检测 Mod 目录变化
	ModInfoPath          string `json:"mod_info_path"`           // Mod 信息文件路径
	configPath           string `json:"-"`
}

var ModInfoChage bool = false // 标记 Mod 信息是否有变更

// 固定的应用名称，用于生成与可执行文件同级的 Data 目录，避免根据 exe 名称改变目录名
const AppName = "KKCardModCheck"

// 只读配置文件路径

var Instance Config

func Load() {
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("json")   // 配置文件类型
	// 配置文件放到可执行文件同级的 <AppName>_Data 目录下
	viper.AddConfigPath(GetDataDir())

	// 设置默认值
	viper.SetDefault("game_path", "")
	viper.SetDefault("card_path", "")
	viper.SetDefault("mod_path", "")
	viper.SetDefault("init_config", false)
	viper.SetDefault("info_build", false)
	viper.SetDefault("info_build_time", "")
	viper.SetDefault("auto_detect_mod_changes", false)
	viper.SetDefault("hide_mode_info_build", false)
	viper.SetDefault("mod_info_path", "")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Creating a new one with default values.")
			Save() // 如果配置文件不存在，创建一个新的配置文件
		} else {
			log.Printf("Error reading config file: %v. Using default values.", err)
		}
	}

	// 将配置值加载到结构体
	if err := viper.Unmarshal(&Instance); err != nil {
		log.Printf("Error unmarshaling config: %v", err)
	}
	Instance.configPath = viper.ConfigFileUsed()
	// 将配置值写入结构体
	Instance.GamePath = viper.GetString("game_path")
	Instance.CardPath = viper.GetString("card_path")
	Instance.ModPath = viper.GetString("mod_path")
	Instance.InitConfig = viper.GetBool("init_config")
	Instance.ModInfoBuild = viper.GetBool("info_build")
	Instance.ModInfoBuildTime = viper.GetString("info_build_time")
	Instance.AutoDetectModChanges = viper.GetBool("auto_detect_mod_changes")
	Instance.HideModeInfoBuild = viper.GetBool("hide_mode_info_build")
	Instance.ModInfoPath = viper.GetString("mod_info_path")
}

// 重置配置为默认值
func ResetToDefault() {
	Instance = Config{
		GamePath:             "",
		CardPath:             "",
		ModPath:              "",
		InitConfig:           false,
		ModInfoBuild:         false,
		ModInfoPath:          "",
		ModInfoBuildTime:     "",
		AutoDetectModChanges: false,
	}
	Save()
}

func Save() {
	viper.Set("game_path", Instance.GamePath)
	viper.Set("card_path", Instance.CardPath)
	viper.Set("mod_path", Instance.ModPath)
	viper.Set("init_config", Instance.InitConfig)
	viper.Set("info_build", Instance.ModInfoBuild)
	viper.Set("info_build_time", Instance.ModInfoBuildTime)
	viper.Set("auto_detect_mod_changes", Instance.AutoDetectModChanges)
	viper.Set("hide_mode_info_build", Instance.HideModeInfoBuild)
	viper.Set("mod_info_path", Instance.ModInfoPath)

	// 尝试写入当前已知配置文件路径；如果不存在则在 Data 目录中创建
	if err := viper.WriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Creating a new one.")
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

// GetDataDir 返回可执行文件同级的 <AppName>_Data 目录路径（如果不存在会创建）
func GetDataDir() string {
	exePath, err := os.Executable()
	if err != nil {
		// 回退到当前目录下的 <AppName>_Data 文件夹（避免根据 exe 名称变化）
		fallback := "./" + AppName + "_Data"
		_ = os.MkdirAll(fallback, 0755)
		return fallback
	}
	// 使用固定的 AppName 而不是可执行文件名，这能防止用户重命名 exe 导致 Data 目录变化
	appName := AppName
	dataDir := filepath.Join(filepath.Dir(exePath), appName+"_Data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Unable to create data dir %s: %v", dataDir, err)
	}
	return dataDir
}
