package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/vs4vijay/lazykubectl/pkg/k8s"
)

var (
	app     *App
	viewArr = []string{"Info", "Namespaces", "Main", "Services"}
	active  = 0
	state   = map[string]string{}
)

type App struct {
	kubeapi *k8s.KubeAPI
}

func NewApp(kubeapi *k8s.KubeAPI) (*App, error) {
	app = &App{
		kubeapi: kubeapi,
	}
	return app, nil
}

func (app *App) Start() {
	// app.kubeapi.Clientset, _ = k8s.Getapp.kubeapi.Clientset(kubeConfig)

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

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

	// Select Panel
	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, selectPanel); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("Namespaces", gocui.MouseLeft, gocui.ModNone, onSelectNamespace); err != nil {
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

	if v, err := g.SetView("Info", 0, 0, gridX-pad, gridY-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Info"
		fmt.Fprintf(v, "Context: %v\n", "")
		fmt.Fprintf(v, "Cluster: %v\n", "")
		fmt.Fprintf(v, "User: %v\n", "")
		fmt.Fprintf(v, "Nodes: %v\n", "")
	}

	if v, err := g.SetView("Namespaces", 0, gridY, gridX-pad, (gridY*2)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Namespaces"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		namespaces, _ := app.kubeapi.SearchNamespaces()
		for _, item := range namespaces {
			fmt.Fprintln(v, item.GetName())
		}
	}

	if v, err := g.SetView("Main", gridX, 0, gridX*4, (gridY*2)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Pods"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Autoscroll = true
	}

	if v, err := g.SetView("Services", 0, gridY*2, gridX*4, gridY*3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Services"
		v.Wrap = true
		fmt.Fprintf(v, "Services: %v\n", "")
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
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	active = nextIndex
	return nil
}

func selectPanel(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(v.Name())
	return err
}

func onSelectNamespace(g *gocui.Gui, v *gocui.View) error {
	namespaceName := getSelectedText(v)
	state["namespace"] = namespaceName

	g.Update(func(g *gocui.Gui) error {
		podsView, _ := g.View("Main")
		podsView.Clear()

		pods, _ := app.kubeapi.SearchPods(namespaceName)
		podsView.Title = fmt.Sprintf("Pods(%v) - %v", len(pods), namespaceName)
		fmt.Fprintf(podsView, "%-20s %-15s\n", "POD NAME", "POD STATUS")
		for _, item := range pods {
			fmt.Fprintf(podsView, "%-20s %-15s\n", item.GetName(), item.Status.Phase)
		}

		servicesView, _ := g.View("Services")
		servicesView.Clear()

		services, _ := app.kubeapi.SearchServices(namespaceName)
		servicesView.Title = fmt.Sprintf("Services(%s)", namespaceName)
		for _, item := range services {
			fmt.Fprintln(servicesView, item.GetName())
		}
		return nil
	})

	return nil
}

func onSelectMain(g *gocui.Gui, view *gocui.View) error {
	selectedData := getSelectedText(view)

	if selectedData == "" {
		return nil
	}
	selectedData = strings.Fields(selectedData)[0]

	if strings.HasPrefix(view.Title, "Pods") {
		// Handling for Pods View
		state["pod"] = selectedData

		g.Update(func(g *gocui.Gui) error {
			view.Clear()
			view.SelBgColor = gocui.ColorBlue

			containers, _ := app.kubeapi.GetContainers(state["namespace"], state["pod"])
			fmt.Fprintf(view, "%-20s\n", "CONTAINER NAME")
			view.Title = fmt.Sprintf("Containers(%v) - %v", len(containers), state["pod"])
			for _, item := range containers {
				fmt.Fprintf(view, "%-20s\n", item.Name)
			}
			return nil
		})
	} else if strings.HasPrefix(view.Title, "Containers") {
		// Handling for Pods View
		state["container"] = selectedData

		g.Update(func(g *gocui.Gui) error {
			view.Clear()
			view.Title = fmt.Sprintf("Logs - %v", state["container"])
			view.Highlight = false
			// view.Editable = true
			// view.Wrap = true
			// view.FgColor = gocui.ColorWhite
			// view.SelBgColor = gocui.ColorBlue

			app.kubeapi.GetContainerLogs(state["namespace"], state["pod"], state["container"], view)
			return nil
		})
	}

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
