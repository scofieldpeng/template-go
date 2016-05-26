# Template-go

this package is base on `html/template` package and just for creating web app more convenient and easy.

## Install

```bash
go get github.com/scofieldpeng/template-go
```

## Usage

initialize template by:
 
```go
tplDirPath := `/path/to/tpldir`
if err := template.New(tplPath);err != nil {
    log.Fataln(err)
}
```

default,the template delimeters are same as the `html/template`,use `{{` and `}}`,if you want to use another one like `[[` and `]]`,you can use this function before invoke `template.New()`:
 
```go
template.SetDelimeter("[[","]]")
```

and the default template file is `.tpl`,you can change by invoke this function:

```go
template.SetTplSuffix(".html")
```

you can invoke the `template.Render()` function in you program

## About the template file

when invoke the `template.New()` function, it will :

1. replace your `{{template "common/header"}}` to `{{template "common/header" .}}`
2. walk your template dictory path and find the all matched template file,the and the `define` actions automatically, and the name is your template file relative path + filename(without file suffix), eg:`/home/www/tpl/index/index.tpl` is your tpl file and the template root path is `/home/www/tpl/`,so the template name is `index/index`
3. parse all your templates file

## Licence 

MIT Licence
