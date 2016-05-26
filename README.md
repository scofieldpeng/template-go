# Template-go

a simple and convinent go template parser base on html/template

assign the directory path and package will look up the all tpl file your define(.tpl,.html etc),the default html/template package need to define the `{{define "xxx"}}{{end}}` at the top/end of the template file, now u are not need to add these boring innovation now,
the template package will add these texts for you automatically when be loaded to the memory,and the usage is so easy,just initialize() and render(),it's easy and happy:-)


