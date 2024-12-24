package proxy

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var conf = `base {
        log_debug = off;
        log_info = off;

        log = "syslog:daemon";

        daemon = on;
        user = redsocks;
        group = redsocks;
        redirector = iptables;
}

redsocks {
        local_ip = 127.0.0.1;
        local_port = 31338;


        ip = PROXYIP;
        port = PROXYPORT;


        type = PROXYTYPE;

        // login = "foobar";
        // password = "baz";
}

redudp {
        local_ip = 127.0.0.1;
        local_port = 10053;


        ip = 127.0.0.1;
        port = 4711;

        dest_ip = 8.8.8.8;
        dest_port = 53;

        udp_timeout = 30;
        udp_timeout_stream = 180;
}

dnstc {
        local_ip = 127.0.0.1;
        local_port = 5300;
}
`

func logMessage(message string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	logFile, err := os.OpenFile("proxy.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	writer := bufio.NewWriter(logFile)
	writer.WriteString(fmt.Sprintf("[*] %s : %s\n", currentTime, message))
	writer.Flush()
}
func changeSocks(proxyType, proxyIp, proxyPort, username, password string) error {
	if proxyType == "http" {
		conf = strings.Replace(conf, "PROXYTYPE", "http-connect", 1)
	} else if proxyType == "socks5" {
		conf = strings.Replace(conf, "PROXYTYPE", "socks5", 1)
	} else if proxyType == "socks4" {
		conf = strings.Replace(conf, "PROXYTYPE", "socks4", 1)
	}
	if len(username) > 0 {
		conf = strings.Replace(conf, "// login = \"foobar\";", "login = \""+username+"\";", 1)
		conf = strings.Replace(conf, "// password = \"baz\";", "password = \""+password+"\";", 1)
	}
	conf = strings.Replace(conf, "PROXYIP", proxyIp, 1)
	conf = strings.Replace(conf, "PROXYPORT", proxyPort, 1)

	err := os.WriteFile("redsocks.conf", []byte(conf), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	// Restart redsocks
	cmd := exec.Command("pkill", "redsocks")
	cmd.Run()

	cmd = exec.Command("/usr/bin/redsocks", "-c", "redsocks.conf")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start redsocks: %v", err)
	}

	logMessage("Run redsocks....")
	return nil
}

func installIptables(proxyRules string) error {
	//cmd := exec.Command("iptables", "-t", "nat", "-N", "REDSOCKS")
	//if err := cmd.Run(); err != nil {
	//	return fmt.Errorf("failed to flush OUTPUT chain: %v", err)
	//}
	//
	//cmd = exec.Command("iptables", "-t", "nat", "-F", "REDSOCKS")
	//if err := cmd.Run(); err != nil {
	//	return fmt.Errorf("failed to flush PREROUTING chain: %v", err)
	//}

	//isRedsocks := checkChainExists("REDSOCKS")
	//if !isRedsocks {
	//	cmd = exec.Command("iptables", "-t", "nat", "-N", "REDSOCKS")
	//	if err := cmd.Run(); err != nil {
	//		return fmt.Errorf("failed to create REDSOCKS chain: %v", err)
	//	}
	//}
	//
	//cmd = exec.Command("iptables", "-t", "nat", "-F", "REDSOCKS")
	//if err := cmd.Run(); err != nil {
	//	return fmt.Errorf("failed to flush REDSOCKS chain: %v", err)
	//}
	proxyRule := strings.Split(proxyRules, "\n")

	rules := [][]string{
		{"-N", "REDSOCKS"},
		{"-F", "REDSOCKS"},
		{"-A", "PREROUTING", "-p", "tcp", "-j", "REDSOCKS"},
		//{"-A", "REDSOCKS", "-p", "tcp", "--dport", 22, "-j", "RETURN"},
		//{"-A", "REDSOCKS", "-d", "<SOCKS API Server>", "-j", "RETURN"},
	}
	for _, line := range proxyRule {
		if len(line) > 0 {
			rules = append(rules, []string{"-A", "REDSOCKS", "-p", "tcp", "-d", line, "-j", "REDIRECT", "--to-port", "31338"})
		}

	}

	for _, rule := range rules {
		cmd := exec.Command("iptables", append([]string{"-t", "nat"}, rule...)...)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add iptables rule: %v", err)
		}
	}

	logMessage("Install Success!")
	return nil
}

func uninstallIptables() error {
	cmd := exec.Command("iptables", "-t", "nat", "-F", "OUTPUT")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to flush OUTPUT chain: %v", err)
	}

	cmd = exec.Command("iptables", "-t", "nat", "-F", "PREROUTING")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to flush PREROUTING chain: %v", err)
	}

	logMessage("Uninstall iptables ...")

	isRedsocks := checkChainExists("REDSOCKS")
	if isRedsocks {
		cmd = exec.Command("iptables", "-t", "nat", "-F", "REDSOCKS")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to flush REDSOCKS chain: %v", err)
		}

		cmd = exec.Command("iptables", "-t", "nat", "-X", "REDSOCKS")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to delete REDSOCKS chain: %v", err)
		}
	}

	return nil
}

func setIptables() error {

	// 构建iptables命令
	cmd := exec.Command("iptables", "-t", "nat", "-A", "OUTPUT", "-p", "tcp", "-m", "owner", "!", "--uid-owner", "redsocks", "-j", "REDSOCKS")
	
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to add OUTPUT rule: %v", err)
	}
	return nil
}

func unsetIptables() error {
	output, err := exec.Command("iptables", "-t", "nat", "-nL", "OUTPUT", "--line-number").Output()
	if err != nil {
		return fmt.Errorf("failed to list OUTPUT chain: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var ids []int

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "REDSOCKS") {
			idStr := strings.Fields(line)[0]
			if id, err := strconv.Atoi(idStr); err == nil {
				ids = append(ids, id)
			}
		}
	}

	logMessage(fmt.Sprintf("REDSOCKS OUTPUT Chain ID: %v", ids))

	for _, id := range ids {
		cmd := exec.Command("iptables", "-t", "nat", "-D", "OUTPUT", fmt.Sprintf("%d", id))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to delete OUTPUT rule: %v", err)
		}
	}

	if len(ids) == 0 {
		logMessage("No Set Iptables ...")
	}

	return nil
}

func checkChainExists(chain string) bool {
	output, err := exec.Command("iptables", "-t", "nat", "-nvL", chain).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

func StartProxy(proxyType, proxyIp, proxyPort, username, password, proxyRules string) error {
	err := uninstallIptables()
	if err != nil {
		return err
	}
	err = installIptables(proxyRules)
	if err != nil {
		return err
	}
	err = unsetIptables()
	if err != nil {
		return err
	}
	err = setIptables()
	if err != nil {
		return err
	}
	if err = changeSocks(proxyType, proxyIp, proxyPort, username, password); err != nil {
		return err
	}

	return nil
}
func CleanProxy() error {
	err := uninstallIptables()
	if err != nil {
		return err
	}
	return nil
}
func getUserID(username string) (string, error) {
	cmd := exec.Command("id", "-u", username)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get UID for user %s: %v", username, err)
	}
	return strings.TrimSpace(string(output)), nil
}
