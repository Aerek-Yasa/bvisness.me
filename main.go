package main

import (
	"fmt"
	"html/template"
	"os"
	"runtime/debug"
	"time"

	"github.com/bvisness/bvisness.me/bhp2"
	"github.com/bvisness/bvisness.me/pkg/images"
)

var hash string = fmt.Sprintf("%d", time.Now().Unix())

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to read build info")
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			hash = setting.Value
		}
	}

	time.Local = time.UTC
}

type Bvisness struct {
	BaseData // shut up, errors

	Desmos DesmosData
}

type BaseData struct {
	Title          string
	Description    string
	OpenGraphImage string // Relative URL within site folder
	Banner         string // Relative URL within site folder
	BannerScale    int    // e.g. 2 for a 2x resolution source image
	LightOnly      bool
}

type Article struct {
	BaseData
	Date time.Time
	Slug string
	Url  string
}

type DesmosData struct {
	NextThreegraphID int
	NextDesmosID     int
}

type Threegraph struct {
	ID int
	JS template.JS
}

type Desmos struct {
	ID   int
	Opts template.JS
	JS   template.JS
}

var bvisnessIncludes = bhp2.FSSearcher{
	FS: os.DirFS("include"),
}

func main() {
	b := bhp2.Instance{
		SrcDir:      "site",
		FSSearchers: []bhp2.FSSearcher{bvisnessIncludes},
		StaticPaths: []string{"apps/"},
		Middleware:  bhp2.ChainMiddleware(images.Middleware),
		Libs: map[string]bhp2.GoLibLoader{
			"images": images.LoadLib,
		},
	}
	b.Run()

	// bhp2.Options[Bvisness]{
	// TODO: Dunno if this is necessary any more.
	// StaticPaths: []string{"apps/"},
	// Funcs: func(b bhp.Instance[Bvisness], r bhp.Request[Bvisness]) template.FuncMap {
	// 	return bhp.MergeFuncMaps(
	// 		images.TemplateFuncs(b, r),
	// 		markdown.TemplateFuncs,
	// 		template.FuncMap{
	// 			// Desmos article
	// 			"threegraph": func(js string) template.HTML {
	// 				result := template.HTML(bhp.Eval(r.T, "desmos/threegraph.html", Threegraph{
	// 					ID: r.User.Desmos.NextThreegraphID,
	// 					JS: template.JS(js),
	// 				}))
	// 				r.User.Desmos.NextThreegraphID++
	// 				return result
	// 			},
	// 			"desmos": func(opts template.JS, js string) template.HTML {
	// 				result := template.HTML(bhp.Eval(r.T, "desmos/desmos.html", Desmos{
	// 					ID:   r.User.Desmos.NextDesmosID,
	// 					Opts: opts,
	// 					JS:   template.JS(js),
	// 				}))
	// 				r.User.Desmos.NextDesmosID++
	// 				return result
	// 			},
	// 		},
	// 	)
	// },
	// 	Middleware: bhp2.ChainMiddleware(images.Middleware[Bvisness]),
	// },
}
