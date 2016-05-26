package template

import (
    "testing"
    "strings"
    "os"
)

func Test_TplNew(t *testing.T) {
    Tpl.SetFns("title",strings.Title)
    gopath := os.Getenv("GOPATH")
    tplDir := gopath + string(os.PathSeparator) + "src/github.com/scofieldpeng/template-go/test_tpl/"
    t.Log("tpldir:",tplDir)
    if err := Tpl.New(tplDir);err != nil {
        t.Fatal("初始化tpl失败!失败原因:",err)
    }

    if err := Tpl.Render(os.Stdout,"index/index",map[string]interface{}{
        "title":"test template",
        "name":"Scofield",
        "age":"24",
        "address":"China",
        "copyright":"hello world",
    });err != nil {
        t.Fatal("render fail!error:",err)
    }
}