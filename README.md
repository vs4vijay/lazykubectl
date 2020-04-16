# lazykubectl

---

## ToDo
- [x] Auth
- [x] Custom Error Handler, for custom errors
- [ ] Proper Logger
- [x] CORS
- [x] Validator
- [ ] Version
- [ ] ENV

---

### Development Notes

```go

    
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
```