package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"k8s.io/client-go/kubernetes"

	"github.com/vs4vijay/lazykubectl/pkg/k8s"
)

var (
	viewArr   = []string{"Info", "Namespaces", "Pods", "Services"}
	active    = 0
	clientset *kubernetes.Clientset
)

func Start(kubeConfig k8s.KubeConfig) {
	clientset, _ = k8s.GetClientset(kubeConfig)

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

	// g.SetCurrentView("Info")

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, selectWidgets); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("Namespaces", gocui.MouseLeft, gocui.ModNone, selectNamespace); err != nil {
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
		namespaces, _ := k8s.SearchNamespaces(clientset)
		for _, item := range namespaces {
			fmt.Fprintln(v, item.GetName())
		}
	}

	if v, err := g.SetView("Pods", gridX, 0, (gridX*4), (gridY*2)-pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Pods"
		// v.Autoscroll = true
	}

	if v, err := g.SetView("Services", 0, gridY*2, (gridX*4), gridY*3); err != nil {
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
	// TODO: Try to use g.View()
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	active = nextIndex
	return nil
}

func selectWidgets(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(v.Name())
	return err
}

func selectNamespace(g *gocui.Gui, v *gocui.View) error {
	namespaceName := getSelectedText(v)

	g.Update(func(g *gocui.Gui) error {
		podsView, _ := g.View("Pods")
		podsView.Clear()

		pods, _ := k8s.SearchPods(clientset, namespaceName)
		for _, item := range pods {
			podsView.Title = fmt.Sprintf("Pods(%s)", namespaceName)
			fmt.Fprintln(podsView, item.GetName())
		}

		servicesView, _ := g.View("Services")
		servicesView.Clear()

		services, _ := k8s.SearchServices(clientset, namespaceName)
		for _, item := range services {
			servicesView.Title = fmt.Sprintf("Services(%s)", namespaceName)
			fmt.Fprintln(servicesView, item.GetName())
		}
		return nil
	})

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

// func renderList(list []interface{}, ) error {
// 	for _, item := range list {
// 		item.()
// 	}
// }

func Try(kubeConfig k8s.KubeConfig) {
	fmt.Println("Rendering")
	// fmt.Println("kubeConfig", kubeConfig)

	clientset, _ := k8s.GetClientset(kubeConfig)

	k8s.SearchNamespaces(clientset)

	k8s.SearchPods(clientset, "kube-system")

	k8s.GetContainers(clientset, "kube-system", "kube-apiserver-kind-control-plane")
	k8s.GetContainers(clientset, "kube-system", "kube-controller-manager-kind-control-plane")
	k8s.GetContainers(clientset, "kube-system", "kube-scheduler-kind-control-plane")

	err := k8s.GetContainerLogs(clientset, "kube-system", "kube-apiserver-kind-control-plane", "kube-apiserver", os.Stdout)

	fmt.Println(err)

	// kube-apiserver-kind-control-plane
}
