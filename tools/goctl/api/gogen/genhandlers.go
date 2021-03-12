package gogen

import (
	"fmt"
	"path"
	"strings"

	"github.com/sjclijie/go-zero/tools/goctl/api/spec"
	"github.com/sjclijie/go-zero/tools/goctl/config"
	"github.com/sjclijie/go-zero/tools/goctl/util"
	"github.com/sjclijie/go-zero/tools/goctl/util/format"
	"github.com/sjclijie/go-zero/tools/goctl/vars"
)

const handlerTemplate = `package handler

import (
	"net/http"

	{{.ImportPackages}}
)

{{ range $index, $element := .Handlers }}
{{ with $element }}
func {{.HandlerName}}(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}{{end}}

		l := logic.New{{.LogicType}}(r.Context(), ctx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err != nil {
			httpx.Error(w, err)
		} else {
			{{if .HasResp}}httpx.OkJson(w, resp){{else}}httpx.Ok(w){{end}}
		}
	}
}


{{ end }}
{{ end }}
`

type Handler struct {
	ImportPackages string
	HandlerName    string
	RequestType    string
	LogicType      string
	Call           string
	HasResp        bool
	HasRequest     bool
}

type RouterGroup struct {
	ImportPackages string
	Handlers       []Handler
}

type Imports []string

func (arr Imports) Contains(s string) bool {
	for _, i := range arr {
		if i == s {
			return true
			break
		}
	}
	return false
}

func (arr Imports) ToString() string {
	return strings.Join(arr, "\r\t")
}

func genHandler(dir string, cfg *config.Config, group spec.Group, route spec.Route) error {

	handler := getHandlerName(route)
	if getHandlerFolderPath(group, route) != handlerDir {
		handler = strings.Title(handler)
	}
	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	return doGenToFile(dir, handler, cfg, group, route, Handler{
		ImportPackages: strings.Join(genHandlerImports(group, route, parentPkg), "\n\t"),
		HandlerName:    handler,
		RequestType:    util.Title(route.RequestTypeName()),
		LogicType:      strings.Title(getLogicName(route)),
		Call:           strings.Title(strings.TrimSuffix(handler, "Handler")),
		HasResp:        len(route.ResponseTypeName()) > 0,
		HasRequest:     len(route.RequestTypeName()) > 0,
	})
}

func genHandlerGroup(dir string, cfg *config.Config, group spec.Group) error {

	routerGroup := RouterGroup{
		Handlers: make([]Handler, 0),
	}

	imports := Imports{}

	for _, route := range group.Routes {

		handler := getHandlerName(route)
		if getHandlerFolderPath(group, route) != handlerDir {
			handler = strings.Title(handler)
		}
		parentPkg, err := getParentPackage(dir)
		if err != nil {
			return err
		}

		importArr := genHandlerImports(group, route, parentPkg)
		for _, importStr := range importArr {
			if !imports.Contains(importStr) {
				imports = append(imports, importStr)
			}
		}

		routerGroup.Handlers = append(routerGroup.Handlers, Handler{
			HandlerName: handler,
			RequestType: util.Title(route.RequestTypeName()),
			LogicType:   strings.Title(getLogicName(route)),
			Call:        strings.Title(strings.TrimSuffix(handler, "Handler")),
			HasResp:     len(route.ResponseTypeName()) > 0,
			HasRequest:  len(route.RequestTypeName()) > 0,
		})
	}

	routerGroup.ImportPackages = strings.Join(imports, "\n\t")

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getGroupHandlerFolderPath(group),
		filename:        fmt.Sprintf("%s%s.go", group.GetAnnotation(groupProperty), "Handler"),
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data:            routerGroup,
	})
}

func doGenToFile(dir, handler string, cfg *config.Config, group spec.Group,
	route spec.Route, handleObj Handler) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getHandlerFolderPath(group, route),
		filename:        filename + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data:            handleObj,
	})
}

func genHandlers(dir string, cfg *config.Config, api *spec.ApiSpec) error {

	/*
		for _, group := range api.Service.Groups {
			for _, route := range group.Routes {
				if err := genHandler(dir, cfg, group, route); err != nil {
					return err
				}
			}
		}
	*/

	for _, group := range api.Service.Groups {
		genHandlerGroup(dir, cfg, group)
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) []string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"",
		util.JoinPackages(parentPkg, getLogicFolderPath(group, route))))
	imports = append(imports, fmt.Sprintf("\"%s\"", util.JoinPackages(parentPkg, contextDir)))
	if len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", util.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/rest/httpx\"", vars.ProjectOpenSourceUrl))

	return imports
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")
	return handler, nil
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return handlerDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(handlerDir, folder)
}

func getGroupHandlerFolderPath(group spec.Group) string {
	folder := group.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		return handlerDir
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(handlerDir, folder)
}

func getHandlerName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Handler"
}

func getLogicName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Logic"
}
