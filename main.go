package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/kargo"
)

var (
	hostname string
	httpAddr string
)

func main() {
	flag.StringVar(&httpAddr, "http", "0.0.0.0:80", "HTTP service address")
	flag.Parse()

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Starting hello-universe...")
	errChan := make(chan error, 10)

	var dm *kargo.DeploymentManager
	if kargo.EnableKubernetes {
		link, err := kargo.Upload(kargo.UploadConfig{
			ProjectID:  "cloud-patterns",
			BucketName: "hello-universe-jlandure",
			ObjectName: "hello-universe",
		})
		fmt.Println("Upload done!")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		env := make(map[string]string)
		env["HELLO_UNIVERSE_TOKEN"] = os.Getenv("HELLO_UNIVERSE_TOKEN")
		fmt.Println("Now deploy...")
		dm = kargo.New()
		err = dm.Create(kargo.DeploymentConfig{
			Args:      []string{"-http=0.0.0.0:80"},
			Env:       env,
			Name:      "hello-universe",
			BinaryURL: link,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = dm.Logs(os.Stdout)
		if err != nil {
			fmt.Println("Local logging has been disabled.")
		}
	} else {
		http.HandleFunc("/", httpHandler)

		go func() {
			errChan <- http.ListenAndServe(httpAddr, nil)
		}()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-errChan:
			if err != nil {
				fmt.Printf("%s - %s\n", hostname, err)
				os.Exit(1)
			}
		case <-signalChan:
			fmt.Printf("%s - Shutdown signal received, exiting...\n", hostname)
			if kargo.EnableKubernetes {
				err := dm.Delete()
				if err != nil {
					fmt.Printf("%s - %s\n", hostname, err)
					os.Exit(1)
				}
			}
			os.Exit(0)
		}
	}
}
