package cmd

import (
    "github.com/spf13/cobra"
    "github.com/levigross/grequests"
    "fmt"
    "os"
    "strings"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "github.com/doun/terminal/color"
)


var (
    APPID  = "20201027000599716"
    SECRET = "qQ7qefdSXnlfYRGTv_Rq"
    SALT   = "hello"
)


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
    From string  `json:"from"`
    To string  `json:"to"`
    TransResult []Result `json:"trans_result"`
}


func GuessLanguange(Q string) string {
    for _, c := range Q {
        if c > 1000 {
            return "en"
        }
    }
    return "zh"
}


var rootCmd = &cobra.Command{
    Use: "bb",
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
                "q": Q,
                "from": "auto",
                "to": GuessLanguange(Q),
                "appid": APPID,
                "salt": SALT,
                "sign": MD5(APPID + Q + SALT + SECRET),
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
                fmt.Println("")
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
