package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type commands struct {
	l *Logic
}

func NewCommands(logic *Logic) *commands {
	return &commands{
		l: logic,
	}
}

//映射相应的命令
func (c *commands) Handlers() map[string]func(args []string) int {
	return map[string]func(args []string) int{
		"0": c.CustomDir,
		//"1":     c.MarkDown,
		//"2":     c.GenerateEntry,
		"3": c.GenerateCURD,
		//"4":     c.CustomFormat,
		"5":     c.ShowTableList,
		"7":     c.Clean,
		"clear": c.Clean,
		"c":     c.Clean,
		"8":     c.Help,
		"h":     c.Help,
		"help":  c.Help,
		"ll":    c.Help,
		"ls":    c.Help,
		"quit":  c.Quit,
		"q":     c.Quit,
		"exit":  c.Quit,
	}
}

//生成数据库表的markdown文档
func (c *commands) MarkDown(args []string) int {
	fmt.Println("Preparing to generate the markdown document...")
	//检查目录是否存在
	CreateDir(c.l.Path)
	err := c.l.CreateMarkdown()
	if err != nil {
		log.Println("MarkDown>>", err)
	}
	return 0
}

//help list
func (c *commands) Help(args []string) int {
	for _, row := range CmdHelp {
		s := fmt.Sprintf("%s %s\n", "NO:"+row.No, row.Msg)
		fmt.Print(s)
	}
	return 0
}

//生成golang表对应的结构实体
func (c *commands) GenerateEntry(args []string) int {
	fmt.Print("Do you need to set the format of the structure?(Yes|No)>")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	switch strings.ToLower(string(line)) {
	case "yes", "y":
		fmt.Print("set formats string start ======== ")
		formats = c._setFormat()
		fmt.Println(formats)
		fmt.Print("set formats string end ======== ")

	}
	err := c.l.CreateEntity(formats)
	if err != nil {
		log.Println("GenerateEntry>>", err.Error())
	}
	go Gofmt(GetExeRootDir())
	return 0
}

//还可以自定义结构体解析实体,如json,gorm,xml
func (c *commands) CustomFormat(args []string) int {
	formats = c._setFormat()
	return 0
}

// 生成相关的curd 相关的参数
func (c *commands) GenerateCURDReq(args []string) (p PackageReq) {

	var tableName, fileName, packageName, structName string

	for tableName == "" {
		tableName = c.SetTableName(args)
	}

	for fileName == "" {
		fileName = c.SetFileName(args)
	}

	packageName = c.SetPackageName(args)

	if packageName == "" {
		packageName = fileName
	}

	structName = c.SetStructName(args)

	if structName == "" {
		structName = fileName
	}

	p.TableName = tableName
	p.FileName = fileName
	p.PackageName = packageName
	p.StructName = structName

	return
}

// 设置packgeName
func (c *commands) SetPackageName(args []string) string {
	fmt.Print("Please set the package name default:file name>")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	return string(line)
}

// 设置tableName
func (c *commands) SetTableName(args []string) string {
	fmt.Print("Please set the table name>")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	return string(line)
}

// 设置FileName
func (c *commands) SetFileName(args []string) string {
	fmt.Print("Please set the file name >")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	return string(line)
}

// 设置structName
func (c *commands) SetStructName(args []string) string {
	fmt.Print("Please set the struct name default:file name >")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	return string(line)
}

//生成golang操作mysql的CRUD增删改查语句
func (c *commands) GenerateCURD(args []string) int {
	req := c.GenerateCURDReq(args)
	fmt.Printf("%+v\n", req)
	err := c.l.CreateCURD(formats, req)
	if err != nil {
		log.Println("GenerateCURD>>", err.Error())
	}
	go Gofmt(c.l.GetFilePath(req))
	return 0
}

//自定义生成目录
func (c *commands) CustomDir(args []string) int {
	fmt.Print("Please set the build directory>")
	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	if string(line) != "" {
		path, err := c.l.T.GenerateDir(string(line))
		if err == nil {
			c.l.Path = path
			fmt.Println("Directory success:", path)
		} else {
			log.Println("Set directory failed>>", err)
		}
	}
	return 0
}

//显示所有的表名
func (c *commands) ShowTableList(args []string) int {
	if len(c.l.DB.Tables) == 0 {
		fmt.Println("Whoops, Nothing at all!!!")
		return 0
	}
	c._showTableList(c.l.DB.Tables)
	//fmt.Print("Select the table sequence number you need?(By default all, comma separated,all represents all)>")
	//line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	//if !strings.EqualFold(string(line), "") {
	//	c.l.DB.DoTables = c._filterTables(string(line), c.l.DB.Tables)
	//}
	return 0
}

//清屏
func (c *commands) Clean(args []string) int {
	Clean()
	return 0
}

//退出
func (c *commands) Quit(args []string) int {
	return 1
}

//过滤表名
func (c *commands) _filterTables(ids string, tables []TableNameAndComment) []TableNameAndComment {
	lst := strings.Split(ids, ",")
	result := make([]TableNameAndComment, 0)
	if strings.ToLower(ids) == "all" {
		return tables
	}
	for _, id := range lst {
		id = strings.TrimSpace(id)
		for _, t := range tables {
			if strconv.Itoa(t.Index) == id || id == t.Name {
				result = append(result, t)
			}
		}
	}
	return result
}

//显示所有名视图
func (c *commands) _showTableList(NameAndComment []TableNameAndComment) {
	for idx, table := range NameAndComment {
		idx++
		info := fmt.Sprintf("%s:%s", strconv.Itoa(idx), table.Name)
		if table.Comment != "" {
			info += fmt.Sprintf("(%s)", table.Comment)
		}
		fmt.Println(info)
	}
	fmt.Println("Total " + strconv.Itoa(len(NameAndComment)) + " tables\n")
}

//set struct format
func (c *commands) _setFormat() []string {
	fmt.Print("Set the mapping name of the structure, separated by a comma (example :json,gorm)>")
	input, _, _ := bufio.NewReader(os.Stdin).ReadLine()
	if string(input) != "" {
		formatList := CheckCharDoSpecialArr(string(input), ',', `[\w\,\-]+`)
		if len(formatList) > 0 {
			fmt.Printf("Set format success: %v\n", formatList)
			return formatList
		}
	}
	fmt.Println("Set failed")
	return nil
}
