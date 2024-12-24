package main

import (
	"ProxifierForLinux/pkg/database"
	"ProxifierForLinux/pkg/proxy"
	"context"
	"fmt"
	"regexp"
)

func init() {
	database.ConnectDatabase()
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetRules() []string {
	var rules []database.Rules
	err := database.Engine.Find(&rules)
	if err != nil {
		return nil
	}
	var result []string
	for _, rule := range rules {
		result = append(result, rule.Type+"://"+rule.Ip+":"+rule.Port)
	}
	return result

}

func (a *App) GetRuleList() []database.Rules {
	var rules []database.Rules
	err := database.Engine.Find(&rules)
	if err != nil {
		return nil
	}
	return rules
}

type CurrentRule struct {
	Type     string `json:"type"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	NeedAuth bool   `json:"needAuth"`
	Username string `json:"username"` // 使用指针表示可选字段
	Password string `json:"password"` // 使用指针表示可选字段
}

func (a *App) SaveRules(currentRule CurrentRule) int {
	//判断该ip、port是否已经有代理了
	var rule database.Rules
	exists, err := database.Engine.Where("ip = ? AND port = ?", currentRule.IP, currentRule.Port).Get(&rule)
	if err != nil {
		return 5 //error
	}
	if exists {
		return 2 // already exists
	}
	var ruleTmp database.Rules
	ruleTmp.Type = currentRule.Type
	ruleTmp.Ip = currentRule.IP
	ruleTmp.Port = currentRule.Port
	ruleTmp.Username = currentRule.Username
	ruleTmp.Password = currentRule.Password
	//插入代理
	_, err = database.Engine.Insert(&ruleTmp)
	if err != nil {
		return 5 //error
	}
	return 1 // success

}
func (a *App) DeleteRule(rule database.Rules) int {
	_, err := database.Engine.Delete(&rule)
	if err != nil {
		return 2
	}
	return 1
}
func (a *App) StartProxy(rule string, proxyIps string) int {
	// 正则表达式提取 type、ip 和 port
	regex := regexp.MustCompile(`^(?P<type>[a-zA-Z0-9]+)://(?P<ip>[0-9\.]+):(?P<port>[0-9]+)$`)

	// 匹配字符串
	matches := regex.FindStringSubmatch(rule)
	if len(matches) == 0 {
		fmt.Println("not matches")
		return 2 // error
	}

	// 提取命名分组的结果
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	var ruleTmp database.Rules
	_, err := database.Engine.Where("ip = ? AND port = ?", result["ip"], result["port"]).Get(&ruleTmp)
	if err != nil {
		fmt.Println(err)
		return 2
	}
	//开始代理
	err = proxy.StartProxy(result["type"], result["ip"], result["port"], ruleTmp.Username, ruleTmp.Password, proxyIps)
	if err != nil {
		fmt.Println(err)
		return 2
	}
	return 1
}
func (a *App) CleanProxy() int {
	err := proxy.CleanProxy()
	if err != nil {
		return 2
	}
	return 1
}
