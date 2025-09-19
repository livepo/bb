package cmd

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doun/terminal/color"
	"github.com/levigross/grequests"
	"github.com/spf13/cobra"
)

var (
	APPID  string
	SECRET string
	SALT   string
)

// Version will be set at build time using -ldflags "-X 'bb/cmd.Version=...'"
var Version = "dev"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("无法获取用户主目录:", err)
		os.Exit(1)
	}
	configPath := filepath.Join(homeDir, ".bb")
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("无法打开配置文件:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "APPID":
			APPID = val
		case "SECRET":
			SECRET = val
		case "SALT":
			SALT = val
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取配置文件出错:", err)
		os.Exit(1)
	}
	if APPID == "" || SECRET == "" || SALT == "" {
		fmt.Println("配置文件缺少必要项(APPID, SECRET, SALT)")
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	// attach version command
	rootCmd.AddCommand(versionCmd)
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type Result struct {
	Src string `json:"src"`
	Dst string `Json:"dst"`
}

type RespStruct struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	TransResult []Result `json:"trans_result"`
}

func GuessLang(Q string) string {
	for _, c := range Q {
		if c > 1000 {
			return "en"
		}
	}
	return "zh"
}

var rootCmd = &cobra.Command{
	Use:   "bb",
	Short: "bb is a terminal translate app",
	Long: `bb uses baidu translate api, support English to Chinese or Chinese to English,
           e.g:
                bb 中国
                bb Hello
          `,
	Run: func(cmd *cobra.Command, args []string) {
		Q := strings.Join(args, " ")
		ro := &grequests.RequestOptions{
			Params: map[string]string{
				"q":     Q,
				"from":  "auto",
				"to":    GuessLang(Q),
				"appid": APPID,
				"salt":  SALT,
				"sign":  MD5(APPID + Q + SALT + SECRET),
			},
		}
		resp, err := grequests.Get("https://fanyi-api.baidu.com/api/trans/vip/translate", ro)
		if err != nil {
			fmt.Println(err)
		}
		var respStruct RespStruct
		if err := json.Unmarshal([]byte(resp.String()), &respStruct); err != nil {
			fmt.Println(err)
		} else {
			for _, result := range respStruct.TransResult {
				color.Println("@{wK}" + result.Src)
				color.Println("@{wK}-------------")
				color.Println("@{bK}" + result.Dst)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
