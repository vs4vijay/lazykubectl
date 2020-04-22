package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/vs4vijay/lazykubectl/pkg/k8s"
	v1 "k8s.io/api/core/v1"
)

var (
	ViewInfo       = "Info"
	ViewNamespaces = "Namespaces"
	ViewServices = "Services"
	ViewMain       = "Main"
	ViewLogs       = "Logs"
)

var (
	app             *App // TODO: Not needed here once we move trigger methods to App struct
	viewSequence    = []string{ViewInfo, ViewNamespaces, ViewMain, ViewLogs}
	activeViewIndex = 0
	state           = map[string]string{} // TODO: Move this to App struct
)

type App struct {
	kubeapi *k8s.KubeAPI
	g       *gocui.Gui
}

func NewApp(kubeapi *k8s.KubeAPI) (*App, error) {
	app = &App{
		kubeapi: kubeapi,
	}
	return app, nil
}

func (app *App) Start() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	app.g = g

	g.InputEsc = true
	g.Highlight = true
	g.Mouse = true
	// g.Cursor = true
	// g.SelBgColor = gocui.ColorCyan
	g.SelFgColor = gocui.ColorBlue

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Rotating Views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	// Refresh Views
	if err := g.SetKeybinding("", 'r', gocui.ModNone, refreshViews); err != nil {
		log.Panicln(err)
	}

	// Select Panel
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, selectView); err != nil {
		log.Panicln(err)
	}

	// Namespaces Handler
	if err := g.SetKeybinding("Namespaces", gocui.MouseLeft, gocui.ModNone, onSelectNamespace); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("Namespaces", gocui.KeyEnter, gocui.ModNone, onSelectNamespace); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("Namespaces", 'd', gocui.ModNone, deleteNamespace); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("Main", gocui.KeyEnter, gocui.ModNone, onSelectMain); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	pad := 1
	gridX, gridY := maxX/4, maxY/3

	if v, err := g.SetView(ViewInfo, 0, 0, gridX-pad, (gridY/2)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Info"

		info, _ := app.kubeapi.GetInfo()
		nodes, _ := app.kubeapi.GetNodes()
		fmt.Fprintf(v, "Cluster: %-10v\n", info["cluster"])
		fmt.Fprintf(v, "Context: %-10v\n", info["context"])
		fmt.Fprintf(v, "User: %-10v\n", info["context"])
		fmt.Fprintf(v, "Nodes: %-10v\n", len(nodes))
	}

	if v, err := g.SetView(ViewNamespaces, 0, gridY/2, gridX-pad, (gridY)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Namespaces"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		renderNamespaces(v.Name())
	}

	if v, err := g.SetView(ViewServices, 0, gridY, gridX-pad, (gridY*2)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Services"
		v.Highlight = true
		v.SelBgColor = gocui.ColorRed
		v.SelFgColor = gocui.ColorBlack

		renderServices(v.Name(), "")
	}

	if v, err := g.SetView(ViewMain, gridX, 0, gridX*4, (gridY*2)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Pods"
		v.Highlight = true
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		renderPods(ViewMain, "")

		// app.kubeapi.GetContainerLogs("kube-system", "kube-apiserver-kind-control-plane", "kube-apiserver", v)
	}

	if v, err := g.SetView(ViewLogs, 0, gridY*2, gridX*4, gridY*3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Logs"
		v.Highlight = true
		v.Wrap = true
		v.Autoscroll = true

		eventWatch, err := app.kubeapi.WatchEvents()
		if err != nil {
			return err
		}
		// TODO: Stop event at some trigger
		// eventWatch.Stop()
		go func() {
			for event := range eventWatch.ResultChan() {
				e, _ := event.Object.(*v1.Event)
				l := fmt.Sprintf("%v (%v) - %v/%v : %v", event.Type, e.Kind, e.Namespace, e.Name, e.Message)
				renderData(v.Name(), l+"\n", false, "")
			}
		}()

	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (activeViewIndex + 1) % len(viewSequence)
	name := viewSequence[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	activeViewIndex = nextIndex
	return nil
}

func selectView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(v.Name())
	return err
}

func onSelectNamespace(g *gocui.Gui, v *gocui.View) error {
	namespaceName := getSelectedText(v)
	state["namespace"] = namespaceName

	renderPods(ViewMain, namespaceName)
	renderServices(ViewServices, namespaceName)

	return nil
}

func onSelectMain(g *gocui.Gui, view *gocui.View) error {
	selectedData := getSelectedText(view)

	if selectedData == "" {
		return nil
	}
	selectedData = strings.Fields(selectedData)[0]

	if strings.HasPrefix(view.Title, "Pods") {
		// Building Containers View
		state["pod"] = selectedData

		renderContainers(ViewMain, state["namespace"], state["pod"])

	} else if strings.HasPrefix(view.Title, "Containers") {
		// Building Logs View
		state["container"] = selectedData

		g.Update(func(g *gocui.Gui) error {
			view.Clear()
			view.Title = fmt.Sprintf("Logs - %v", state["container"])
			view.Highlight = false
			view.Autoscroll = true
			// view.Editable = true
			// view.Wrap = true
			// view.FgColor = gocui.ColorWhite
			// view.SelBgColor = gocui.ColorBlue

			app.kubeapi.GetContainerLogs(state["namespace"], state["pod"], state["container"], view)

			// logWatch, err := app.kubeapi.WatchPodLogs(state["namespace"], state["pod"])
			// if err != nil {
			// 	return err
			// }
			// go func() {
			// 	for event := range logWatch.ResultChan() {
			// 		// e, _ := event.Object.(*v1.Namespace)
			// 		fmt.Printf("%v \n", event)
			// 		// fmt.Printf("%v : %v \n", e.GetName(), e.Message)
			// 		renderData(ViewLogs, "string(event.Type)" + "\n", false)
			// 	}
			// }()

			return nil
		})
	}

	return nil
}

func deleteNamespace(g *gocui.Gui, v *gocui.View) error {
	namespaceName := getSelectedText(v)
	err := app.kubeapi.DeleteNamespaces(namespaceName)
	if err != nil {
		return err
	}
	renderNamespaces(v.Name())
	return nil
}

func getSelectedText(view *gocui.View) string {
	_, cy := view.Cursor()
	line, err := view.Line(cy)
	if err != nil {
		return ""
	}
	return line
}

func renderData(viewName string, data string, clear bool, title string) {
	app.g.Update(func(g *gocui.Gui) error {
		view, _ := g.View(viewName)
		if clear {
			view.Clear()
		}
		if title != "" {
			view.Title = title
		}
 		fmt.Fprintf(view, data)
		return nil
	})
}

func renderNamespaces(viewName string) {
	namespaces, _ := app.kubeapi.GetNamespaces()
	var ns []string
	ns = append(ns, "\n")
	for _, item := range namespaces {
		ns = append(ns, item.GetName())
	}
	renderData(viewName, strings.Join(ns, "\n"), true, "")
}

func renderPods(viewName string, namespaceName string) {
	pods, _ := app.kubeapi.GetPods(namespaceName)
	title := fmt.Sprintf("Pods(%v)", len(pods))
	var pos []string
	for _, item := range pods {
		pos = append(pos, item.GetName())
	}
	renderData(viewName, strings.Join(pos, "\n"), true, title)
}

func renderContainers(viewName string, namespaceName string, podName string) {
	containers, _ := app.kubeapi.GetContainers(namespaceName, podName)
	title := fmt.Sprintf("Containers(%v) - %v", len(containers), podName)
	var cs []string
	for _, item := range containers {
		cs = append(cs, item.Name)
	}
	renderData(viewName, strings.Join(cs, "\n"), true, title)
}

func renderServices(viewName string, namespaceName string) {
	services, _ := app.kubeapi.GetServices(namespaceName)
	var svcs []string
	for _, item := range services {
		svcs = append(svcs, item.GetName())
	}
	renderData(viewName, strings.Join(svcs, "\n"), true, "Services")
}

func refreshViews(g *gocui.Gui, v *gocui.View) error {
	renderNamespaces(ViewNamespaces)
	return nil
}
