package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AdminUsername        string
	AdminPassword        string
	ServerAddress        string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func ReadConfig() (Config, error) {
	var config Config

	file, err := os.Open("../.env")
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		lineSplit := strings.Split(line, " ")
		switch lineSplit[0] {
		case "admin_username:":
			if lineSplit[1] == "CHANGE" {
				return config, fmt.Errorf("admin_username in .env unchanged")
			}
			config.AdminUsername = lineSplit[1]
		case "admin_password:":
			if lineSplit[1] == "CHANGE" {
				return config, fmt.Errorf("admin_password in .env unchanged")
			}
			config.AdminPassword = lineSplit[1]
		case "server_address:":
			config.ServerAddress = lineSplit[1]
		case "access_token_duration:":
			sliceLimit := len(lineSplit[1]) - 1
			timeUnitChar := lineSplit[1][sliceLimit:]
			timeUnit := charToTimeValue(timeUnitChar)

			timeValueChar := lineSplit[1][:sliceLimit]
			timeValue, err := strconv.Atoi(timeValueChar)
			if err != nil {
				return config, err
			}
			config.AccessTokenDuration = time.Duration(timeValue) * timeUnit
		case "refresh_token_duration:":
			sliceLimit := len(lineSplit[1]) - 1
			timeUnitChar := lineSplit[1][sliceLimit:]
			timeUnit := charToTimeValue(timeUnitChar)

			timeValueChar := lineSplit[1][:sliceLimit]
			timeValue, err := strconv.Atoi(timeValueChar)
			if err != nil {
				return config, err
			}
			config.RefreshTokenDuration = time.Duration(timeValue) * timeUnit
		}
	}

	if err := scanner.Err(); err != nil {
		return config, fmt.Errorf("Error scanning file:%w", err)
	}

	return config, nil
}

func charToTimeValue(timeChar string) time.Duration {
	switch timeChar {
	case "s":
		return time.Second
	case "m":
		return time.Minute
	case "h":
		return time.Hour
	default:
		return time.Minute
	}
}
