# lazykubectl

---

## References
- https://pkg.go.dev/k8s.io/client-go/kubernetes?tab=doc
- 

## ToDo
- [x] Auth
- [x] Custom Error Handler, for custom errors
- [ ] Proper Logger
- [x] CORS
- [x] Validator
- [ ] Version
- [ ] ENV
- Handle Up / Down Arrow
- Stream Logs
- Events

---

### Development Notes

```go
    
Docker:

import (
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
)

cli, err := client.NewEnvClient()
cli.Info(context.Background())
cli.DiskUsage(context.Background())
cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})

Client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
stream, err := c.Client.ContainerStats(context.Background(), container.ID, true)
images, err := c.Client.ImageList(context.Background(), types.ImageListOptions{})
result, err := c.Client.VolumeList(context.Background(), filters.Args{})



Kubernetes:

import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)



    
    watch, _ := api.Services("").Watch(metav1.ListOptions{})
    for event := range watch.ResultChan() {
        fmt.Printf("Type: %v\n", event.Type)
        p, ok := event.Object.(*v1.Pod)
        if !ok {
            fmt.Errorf("unexpected type")
        }
        fmt.Println(p.Status.ContainerStatuses)
        fmt.Println(p.Status.Phase)
    }

	deploymentsClient := clientset.ExtensionsV1beta1().Deployments("namespace-ffledgling")

	// List existing deployments in namespace
	deployments, err := deploymentsClient.List(metav1.ListOptions{})


	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// Take required information from error and context and send it to a service like New Relic
		fmt.Println(c.Path(), c.QueryParams(), err.Error())

		switch err.(type) {
		case orchestrator.CustomError:
			fmt.Println("custom")
		default:
			fmt.Println("normal") // here v has type interface{}
		}

		// Call the default handler to return the HTTP response
		e.DefaultHTTPErrorHandler(err, c)
	}

    https://github.com/alitari/kubexp

    https://github.com/JulienBreux/pody



	if v, err := g.SetView("help", maxX-25, 0, maxX-1, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "Space: New View")
		fmt.Fprintln(v, "Tab: Next View")
		fmt.Fprintln(v, "← ↑ → ↓: Move View")
		fmt.Fprintln(v, "Backspace: Delete View")
		fmt.Fprintln(v, "t: Set view on top")
		fmt.Fprintln(v, "b: Set view on bottom")
		fmt.Fprintln(v, "^C: Exit")
	}


func Loader() string {
	characters := "|/-\\"
	now := time.Now()
	nanos := now.UnixNano()
	index := nanos / 50000000 % int64(len(characters))
	return characters[index : index+1]
}

https://stackoverflow.com/questions/40975307/how-to-watch-events-on-a-kubernetes-service-using-its-go-client

cache.NewInformer
NewSharedIndexInformer




    CPU
    MEM
    View Logs
    Execute Shell
    ? Events
```