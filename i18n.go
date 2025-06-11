package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var I18nBundle *i18n.Bundle

func InitI18n() error {
	I18nBundle = i18n.NewBundle(language.Chinese)
	I18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	langDir := "locales"
	if _, err := os.Stat(langDir); os.IsNotExist(err) {
		if err := os.MkdirAll(langDir, 0755); err != nil {
			return err
		}
		CreateDefaultLanguageFile(filepath.Join(langDir, "zh.toml"))
		CreateDefaultLanguageFile(filepath.Join(langDir, "en.toml"))
	}

	files, err := os.ReadDir(langDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".toml" {
			_, err := I18nBundle.LoadMessageFile(filepath.Join(langDir, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateDefaultLanguageFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	isEnglish := filepath.Base(path) == "en.toml"

	var content string

	if isEnglish {
		content = `
# 欢迎消息
[welcome]
other = "Welcome"
description = "Welcome message"

# 错误消息
[error]
other = "Error"
description = "Error message"

# 命令帮助
[cmd_help]
other = "Usage: gits [options] --cmd \"git command\""
description = "Command help message"

# 选项帮助
[options_help]
other = "Options:"
description = "Options help message"
`
	} else {
		content = `
# 欢迎消息
[welcome]
other = "欢迎使用"
description = "欢迎消息"

# 错误消息
[error]
other = "错误"
description = "错误消息"

# 命令帮助
[cmd_help]
other = "使用方法: gits [选项] --cmd \"git命令\""
description = "命令帮助信息"

# 选项帮助
[options_help]
other = "选项:"
description = "选项帮助信息"
`
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func i(id string, args ...interface{}) string {
	localizer := i18n.NewLocalizer(I18nBundle, language.Chinese.String())

	if len(args) > 0 {
		templateData := make(map[string]interface{})
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				if key, ok := args[i].(string); ok {
					templateData[key] = args[i+1]
				}
			}
		}

		msg, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    id,
			TemplateData: templateData,
		})

		if err != nil {
			return id
		}
		return msg
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})

	if err != nil {
		return id
	}
	return msg
}
