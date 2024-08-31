package config

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

type JwtConfig struct {
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Expire    time.Duration
}

func ParseExpireTime(expireStr string) (time.Duration, error) {
	// 正则表达式，匹配格式为数字+单位 (可选)
	re := regexp.MustCompile(`^(\d+)([smhd]?)$`)

	// 提取匹配结果
	matches := re.FindStringSubmatch(expireStr)
	if len(matches) != 3 {
		return 0, errors.New("invalid expire format")
	}

	// 解析数字部分
	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	// 根据单位返回相应的时间段
	unit := matches[2]
	switch unit {
	case "s", "": // 秒 (默认)
		return time.Duration(value) * time.Second, nil
	case "m": // 分钟
		return time.Duration(value) * time.Minute, nil
	case "h": // 小时
		return time.Duration(value) * time.Hour, nil
	case "d": // 天
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, errors.New("invalid time unit")
	}
}

func (c *Config) GetJwtConfig() *JwtConfig {
	secretKey := c.v.GetString("jwt.secret_key")
	expire := c.v.GetDuration("jwt.expire")
	return &JwtConfig{
		SecretKey: secretKey,
		Expire:    expire,
	}
}
