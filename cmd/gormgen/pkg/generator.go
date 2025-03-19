package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
)

// The Generator is the one responsible for generating the code, adding the imports, formating, and writing it to the file.
type Generator struct {
	buf           map[string]map[string]*bytes.Buffer
	inputFile     string
	config        config
	structConfigs []structConfig
}

// NewGenerator function creates an instance of the generator given the name of the output file as an argument.
func NewGenerator(outputFile string) *Generator {
	return &Generator{
		buf:       map[string]map[string]*bytes.Buffer{},
		inputFile: outputFile,
	}
}

// ParserAST parse by go file
func (g *Generator) ParserAST(p *Parser, structs []string) (ret *Generator) {
	for _, v := range structs {
		g.buf[gorm.ToDBName(v)] = make(map[string]*bytes.Buffer)
	}
	g.structConfigs = p.Parse()
	g.config.PkgName = p.pkg.Name
	g.config.Helpers = structHelpers{
		Titelize: strings.Title,
	}
	g.config.QueryBuilderName = SQLColumnToHumpStyle(p.pkg.Name) + "QueryBuilder"
	return g
}

func (g *Generator) checkConfig() (err error) {
	if len(g.config.PkgName) == 0 {
		err = errors.New("package name dose'n set")
		return
	}
	for i := 0; i < len(g.structConfigs); i++ {
		g.structConfigs[i].config = g.config
	}
	return
}

// Generate executes the template and store it in an skitii buffer.
func (g *Generator) Generate() *Generator {
	if err := g.checkConfig(); err != nil {
		panic(err)
	}

	for i, v := range g.structConfigs {
		if _, ok := g.buf[gorm.ToDBName(v.StructName)]; !ok {
			continue
		}
		dir, err2 := os.Getwd()
		if err2 != nil {
			panic(err2)
		}
		index := strings.Index(dir, "go-gin-api")
		if index == -1 {
			panic("not found go-gin-api")
		}
		rootDir := dir[:index] + "go-gin-api"
		fmt.Println("rootDir: ", rootDir)

		// 生成 entity
		err := g.generateByTemplateFilePath(i, v, rootDir+"/cmd/gormgen/template/entity.txt")
		if err != nil {
			panic(err)
		}
		// 生成 repo
		err = g.generateByTemplateFilePath(i, v, rootDir+"/cmd/gormgen/template/repo.txt")
		if err != nil {
			panic(err)
		}
		// 生成default_repo
		err = g.generateByTemplateFilePath(i, v, rootDir+"/cmd/gormgen/template/default_repo.txt")
		if err != nil {
			panic(err)
		}
	}

	return g
}

func (g *Generator) generateByTemplateFilePath(i int, v structConfig, filePath string) error {
	templateBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// filePath 中取文件名作为模板名称
	templateName := strings.Split(filePath, "/")[len(strings.Split(filePath, "/"))-1]
	// 去除文件后缀
	templateName = strings.Split(templateName, ".")[0]
	g.structConfigs[i].TemplateName = templateName
	outputTemplate = parseTemplateOrPanic(string(templateBytes))
	g.buf[gorm.ToDBName(v.StructName)][templateName] = new(bytes.Buffer)
	fmt.Println("templateName: ", templateName)
	if err := outputTemplate.Execute(g.buf[gorm.ToDBName(v.StructName)][templateName], v); err != nil {
		panic(err)
	}
	return err
}

// Format function formats the output of the generation.
func (g *Generator) Format() *Generator {
	for k := range g.buf {
		for k2, v := range g.buf[k] {
			formattedOutput, err := format.Source(v.Bytes())
			if err != nil {
				panic(err)
			}
			g.buf[k][k2] = bytes.NewBuffer(formattedOutput)
		}
	}
	return g
}

// Flush function writes the output to the output file.
func (g *Generator) Flush() error {
	for k := range g.buf {
		for k2, v := range g.buf[k] {
			filename := "/gen_" + strings.ToLower(k) + "_" + k2 + ".go"
			filepath := g.inputFile + filename
			if err := ioutil.WriteFile(filepath, v.Bytes(), 0777); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("  └── file : ", fmt.Sprintf("%s/%s", strings.ToLower(k), filename))
		}
	}
	return nil
}
