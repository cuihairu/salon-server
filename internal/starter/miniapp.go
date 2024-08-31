package starter

type WechatConfig struct {
	AppID     string `mapstructure:"app_id" yaml:"app_id" json:"app_id"`
	AppSecret string `mapstructure:"app_secret" yaml:"app_secret" json:"app_secret"`
}

type MiniappConfig struct {
	Wechat WechatConfig `mapstructure:"wechat" yaml:"wechat" json:"wechat"`
}
