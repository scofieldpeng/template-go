// template template包基于官方的html/template,为了偷懒和便利性而做了此包,此包基本为echo框架而开发,方便在使用echo框架进行web开发时的操作
// TODO: 1. 启动时的公共配置参数
package template

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// tplT 定义template模型
type tplT struct {
	lastTemplate *template.Template     // 上次缓存
	template     *template.Template     // 本次模板
	tplParsedMap map[string]string      // 模板内容缓存,key为模板的文件名,value为解析后的模板文件内容
	fns          map[string]interface{} // 公共函数
	vars         map[string]interface{} // 公共值
	delimeter    struct {
		left  string // 左侧模板分隔符
		right string // 右侧模板分隔符
	} // 模板的分隔符
	tplFileSuffix string // 模板文件后缀
	rootPath      string // 模板根目录
}

var (
	Tpl tplT
)

func new() tplT {
	return tplT{
		delimeter: struct {
			left  string
			right string
		}{
			left:  DefaultLeftDelimeter,
			right: DefaultRightDelimeter,
		},
		tplFileSuffix: DefaultTplSuffix,
		tplParsedMap:  make(map[string]string),
		fns:           make(map[string]interface{}),
		vars:          make(map[string]interface{}),
		lastTemplate:  &template.Template{},
		template:      &template.Template{},
	}
}

// SetFns 设置运行函数
func (t *tplT) SetFns(name string, fns interface{}) {
	t.fns[name] = fns
}

// SetDelimeter 设置delimeter
func (t *tplT) SetDelimeter(left, right string) {
	if left != "" {
		t.delimeter.left = left
	}
	if right != "" {
		t.delimeter.right = right
	}
}

// SetTplSuffix 设置模板的后缀(含.,例如.html)
func (t *tplT) SetTplSuffix(suffix string) {
	if suffix != "" && strings.Index(suffix, ".") == 0 {
		t.tplFileSuffix = suffix
	}
}

// loopPath 从文件夹中遍历所有符合条件的模板文件路径
func (t *tplT) loopDir(root string) ([]string, error) {
	filePath := make([]string, 0)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, t.tplFileSuffix) {
			filePath = append(filePath, path)
		}
		return nil
	}); err != nil {
		return []string{}, err
	}
	return filePath, nil
}

// addDefine 给每个模板文件头和尾添加{{define xxx}}...{{end}}的字符串
func (t *tplT) addDefine(tplName, content string) string {
	return fmt.Sprintf("%sdefine \"%s\"%s\n", t.delimeter.left, tplName, t.delimeter.right) + content + fmt.Sprintf("\n%send%s", t.delimeter.left, t.delimeter.right)
}

// transferTemplate 将模板中的{{template xxx}}默认转化为{{template xxx .}}
func (t *tplT) transferTemplate(content string) string {
	reg := regexp.MustCompile(fmt.Sprintf(`%stemplate\s(\S*)%s`, t.delimeter.left, t.delimeter.right))
	return reg.ReplaceAllString(content, fmt.Sprintf(`%stemplate "$1" .%s`, t.delimeter.left, t.delimeter.right))
}

// New 初始化模板
func (t *tplT) New(dirPath string) error {
	// 给文件夹的末尾添加上slash
	if string([]byte(dirPath)[len(dirPath)-1]) != string(os.PathSeparator) {
		dirPath = dirPath + string(os.PathSeparator)
	}

	t.rootPath = dirPath

	// 遍历获取到所有的tpl绝对路径
	tplpaths, err := t.loopDir(dirPath)
	if err != nil {
		return err
	}

	var tmpTemplate *template.Template
	if tmpTemplate,err = t.parseFiles(tmpTemplate, tplpaths...); err != nil {
		return err
	}

	// 更新当前最新模板
	t.lastTemplate = t.template
	t.template = tmpTemplate

	return nil
}

// parseFiles 依次解析相应的文件夹并
func (t *tplT) parseFiles(tmpTemplate *template.Template, tplpaths ...string) (*template.Template,error) {
	for _, tplpath := range tplpaths {
		fileObj, err := os.Open(tplpath)
		if err != nil {
			return tmpTemplate,err
		}

		// 根据该模板的路径解析出该模板文件的template调用名称
		fileBytes, err := ioutil.ReadAll(fileObj)
		file := strings.Split(tplpath, t.rootPath)[1]
		filename := strings.Split(file, t.tplFileSuffix)[0]

		// 先转换模板,然后在头和尾加上{{define}}
		t.tplParsedMap[tplpath] = strings.TrimSpace(t.addDefine(filename, t.transferTemplate(string(fileBytes))))

		// 解析模板
		var tmpl *template.Template
		if tmpTemplate == nil {
			tmpTemplate = template.New(filename)
		}
		if filename == tmpTemplate.Name() {
			tmpl = tmpTemplate
		} else {
			tmpl = tmpTemplate.New(filename)
		}
		_, err = tmpl.Funcs(template.FuncMap(t.fns)).Parse(t.tplParsedMap[tplpath])
		if err != nil {
			return tmpTemplate,err
		}

		fileObj.Close()
	}

	return tmpTemplate,nil
}

// Render 渲染函数
func (t *tplT) Render(w io.Writer, tplName string, data interface{}) error {
	return t.template.ExecuteTemplate(w, tplName, data)
}

func init() {
	Tpl = new()
}
