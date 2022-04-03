package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func ReadAndUpdateConfig(config *viper.Viper, name string, version string) {
	created := false
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("%s configuration file not found: %v\n", name, err)
			config.SetDefault("version", version)
			if err := config.SafeWriteConfig(); err != nil {
				panic(fmt.Errorf("failed to write %s configuration file: %v", name, err))
			}
			fmt.Printf("default %s configuration file created.\n", name)
			created = true
		} else {
			panic(fmt.Errorf("fatal error reading %s configuration: %v", name, err))
		}
	}
	if created {
		return
	}
	fileVersion := config.GetString("version")
	config.SetDefault("version", version)
	if cmp := VersionCompare(version, fileVersion); cmp != 0 {
		fmt.Printf("%s configuration file version mismatch. expect version: %s, but got: %s.\n", name, version, config.GetString("version"))
		if cmp < 0 {
			fmt.Printf("you may need to update your app vesion.\n")
		}
		if cmp > 0 {
			fmt.Printf("updating your %s configuration file. conflict entries will not be updated.\n", name)
			config.Set("version", version)
			if err := config.WriteConfig(); err != nil {
				panic(fmt.Errorf("failed to write %s configuration file: %v", name, err))
			}
			fmt.Printf("%s configuration file updated to version %s.\n", name, version)
		}
	}
}

func VersionCompare(a, b string) int {
	a = strings.SplitN(a, " ", 2)[0]
	a = strings.SplitN(a, "-", 2)[0]
	aslice := strings.Split(a, ".")
	b = strings.SplitN(b, " ", 2)[0]
	b = strings.SplitN(b, "-", 2)[0]
	bslice := strings.Split(b, ".")
	for i := 0; i < len(aslice) && i < len(bslice); i++ {
		ai, _ := strconv.Atoi(aslice[i])
		bi, _ := strconv.Atoi(bslice[i])
		if ai > bi {
			return 1
		}
		if ai < bi {
			return -1
		}
	}
	if len(aslice) > len(bslice) {
		return 1
	}
	if len(aslice) < len(bslice) {
		return -1
	}
	return 0
}
