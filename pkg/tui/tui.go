package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"

	"github.com/vs4vijay/lazykubectl/pkg/k8s"
)

var (
	viewArr = []string{"Info", "Namespaces", "Pods", "Services"}
	active  = 0
)

func Start(kubeConfig k8s.KubeConfig) {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.InputEsc = true
	g.Cursor = true
	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln("SetKeybinding --- ", err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("Pods", gocui.KeyCtrlP, gocui.ModNone, quit); err != nil {
		log.Panicln("SetKeybinding --- ", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln("MainLoop --- ", err)
	}

	// g.SetViewOnTop("Namespaces")
	// g.SetCurrentView("Namespaces")
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	pad := 1
	gridX, gridY := maxX/4, maxY/3

	if v, err := g.SetView("Info", 0, 0, gridX - pad, gridY - pad); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Info"
		v.Editable = true
		v.Wrap = true

		fmt.Fprintf(v, "Context: %v\n", "")
		fmt.Fprintf(v, "Cluster: %v\n", "")
		fmt.Fprintf(v, "User: %v\n", "")
		fmt.Fprintf(v, "Nodes: %v\n", "")
	}

	if v, err := g.SetView("Namespaces", 0, gridY, gridX - pad, (gridY * 2) - pad) ; err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Namespaces"
		v.Editable = true
		v.Wrap = true
		fmt.Fprintf(v, "Namespaces: %v\n", "")
	}

	if v, err := g.SetView("Pods", gridX, 0, (gridX * 4), (gridY * 2) - pad) ; err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Pods"
		v.Editable = true
		v.Autoscroll = true
		fmt.Fprintf(v, "Pods: %v\n", "")
	}

	if v, err := g.SetView("Services", 0, gridY * 2, (gridX * 4), gridY * 3) ; err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Services"
		v.Editable = true
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

	// out, err := g.View("Pods")
	// if err != nil {
	// 	return err
	// }
	// fmt.Fprintf(out, "Going from view %v to %v\n", v.Name(), name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

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
