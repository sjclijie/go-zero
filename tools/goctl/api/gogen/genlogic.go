package gogen

import (
	"fmt"
	"path"
	"strings"

	"github.com/sjclijie/go-zero/tools/goctl/api/spec"
	"github.com/sjclijie/go-zero/tools/goctl/config"
	ctlutil "github.com/sjclijie/go-zero/tools/goctl/util"
	"github.com/sjclijie/go-zero/tools/goctl/util/format"
	"github.com/sjclijie/go-zero/tools/goctl/vars"
)

const logicTemplate = `package logic

import (
	{{.Imports}}
)

type {{.LogicName}} struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{.LogicName}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.LogicName}} {
	return &{{.LogicName}}{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
{{ range $index, $element := .Logics }}
func (l *{{$.LogicName}}) {{$element.Function}}({{$element.Request}}) {{$element.ResponseType}} {
	// todo: add your logic here and delete this line

	{{$element.ReturnString}}
}
{{ end }}
`

type LogicData struct {
	Function     string
	ResponseType string
	ReturnString string
	Request      string
}

type genLogicByGroupConfig struct {
	LogicName string
	Imports   string
	Logics    []LogicData
}

func genLogic(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	/*
		for _, g := range api.Service.Groups {
			for _, r := range g.Routes {
				err := genLogicByRoute(dir, cfg, g, r)
				if err != nil {
					return err
				}
			}
		}
	*/

	for _, group := range api.Service.Groups {
		if err := genLogicByGroup(dir, cfg, group); err != nil {
			return err
		}
	}

	return nil
}

func genLogicByRoute(dir string, cfg *config.Config, group spec.Group, route spec.Route) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	imports := genLogicImports(route, parentPkg)
	var responseString string
	var returnString string
	var requestString string
	if len(route.ResponseTypeName()) > 0 {
		resp := responseGoTypeName(route, typesPacket)
		responseString = "(" + resp + ", error)"
		if strings.HasPrefix(resp, "*") {
			returnString = fmt.Sprintf("return &%s{}, nil", strings.TrimPrefix(resp, "*"))
		} else {
			returnString = fmt.Sprintf("return %s{}, nil", resp)
		}
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestTypeName()) > 0 {
		requestString = "req " + requestGoTypeName(route, typesPacket)
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getLogicFolderPath(group, route),
		filename:        goFile + ".go",
		templateName:    "logicTemplate",
		category:        category,
		templateFile:    logicTemplateFile,
		builtinTemplate: logicTemplate,
		data: map[string]string{
			"imports":      imports,
			"logic":        strings.Title(logic),
			"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
			"responseType": responseString,
			"returnString": returnString,
			"request":      requestString,
		},
	})
}

func genLogicByGroup(dir string, cfg *config.Config, group spec.Group) error {
	groupName := group.GetAnnotation(groupProperty)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, fmt.Sprintf("%sLogic", groupName))
	if err != nil {
		return err
	}

	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	config := genLogicByGroupConfig{LogicName: strings.Title(goFile)}
	imports := make(Imports, 0)
	config.Logics = make([]LogicData, 0)

	for _, route := range group.Routes {
		importArray := genLogicImportArray(route, parentPkg)
		for _, item := range importArray {
			if !imports.Contains(item) {
				imports = append(imports, item)
			}
		}

		var responseString string
		var returnString string
		var requestString string
		var logicFunction string

		logicFunction = getLogicName(route)

		if len(route.ResponseTypeName()) > 0 {
			resp := responseGoTypeName(route, typesPacket)
			responseString = "(" + resp + ", error)"
			if strings.HasPrefix(resp, "*") {
				returnString = fmt.Sprintf("return &%s{}, nil", strings.TrimPrefix(resp, "*"))
			} else {
				returnString = fmt.Sprintf("return %s{}, nil", resp)
			}
		} else {
			responseString = "error"
			returnString = "return nil"
		}
		if len(route.RequestTypeName()) > 0 {
			requestString = "req " + requestGoTypeName(route, typesPacket)
		}

		config.Logics = append(config.Logics, LogicData{
			Function:     logicFunction,
			ResponseType: responseString,
			ReturnString: returnString,
			Request:      requestString,
		})
	}

	config.Imports = imports.ToString()

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getLogicGroupFolderPath(group),
		filename:        goFile + ".go",
		templateName:    "logicTemplate",
		category:        category,
		templateFile:    logicTemplateFile,
		builtinTemplate: logicTemplate,
		data:            config,
	})
}

func getLogicFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return logicDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDir, folder)
}

func getLogicGroupFolderPath(group spec.Group) string {
	folder := group.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		return logicDir
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDir, folder)
}

func genLogicImports(route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, `"context"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s\"", ctlutil.JoinPackages(parentPkg, contextDir)))
	if len(route.ResponseTypeName()) > 0 || len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", ctlutil.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceUrl))
	return strings.Join(imports, "\n\t")
}

func genLogicImportArray(route spec.Route, parentPkg string) []string {
	var imports []string
	imports = append(imports, `"context"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s\"", ctlutil.JoinPackages(parentPkg, contextDir)))
	if len(route.ResponseTypeName()) > 0 || len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", ctlutil.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceUrl))
	return imports
}
