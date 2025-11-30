package templates

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed *.html
var TemplateFS embed.FS

// LoadTemplates 加载嵌入的模板文件
func LoadTemplates() (*template.Template, error) {
	// 读取嵌入的文件系统
	tmpl := template.New("")
	
	// 遍历嵌入的文件系统中的所有HTML文件
	err := fs.WalkDir(TemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// 只处理HTML文件
		if !d.IsDir() && (path[len(path)-5:] == ".html") {
			// 读取文件内容
			content, err := TemplateFS.ReadFile(path)
			if err != nil {
				return err
			}
			
			// 创建模板
			_, err = tmpl.New(path).Parse(string(content))
			if err != nil {
				return err
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return tmpl, nil
}