package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
)

func eventCreated(obj interface{}) {
	ev := obj.(*api.Event)
	fmt.Printf("%v %q\n", ev.LastTimestamp, ev.Message)
}

func watchEvents(client *client.Client) cache.Store {

	watchlist := cache.NewListWatchFromClient(client, "events", api.NamespaceAll, fields.Everything())
	resyncPeriod := 0 * time.Minute

	//Setup an informer to call functions when the watchlist changes
	store, controller := framework.NewInformer(
		watchlist,
		&api.Event{},
		resyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc: eventCreated,
		},
	)
	//Run the controller as a goroutine
	go controller.Run(wait.NeverStop)
	return store
}

/*
Hacked together listener using Kubernetes restclient.
*/
func restVersion() {
	config := &restclient.Config{
		Host:     "http://localhost:8080",
		Insecure: true,
	}
	//Create a new client to interact with cluster and freak if it doesn't work
	kubeClient, err := client.New(config)

	if err != nil {
		log.Fatalln("Client not created sucessfully:", err)
	}

	//This actually returns a cache, but we don't really need it...
	watchEvents(kubeClient)
}

func simpleHTTP() {
	resp, err := http.Get("http://localhost:8080/api/v1/events?watch=true")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get events from API server: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err.Error() != "EOF" {
			fmt.Fprintf(os.Stdout, "Error reading event: %v\n", err)
		}
		if len(line) > 0 {
			fmt.Println(line)
		}
	}

}

func main() {
	restVersion()

	for {
		time.Sleep(time.Second)
	}
}
