package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/types"
)

func main() {

	nodeName := ""
	details := flag.Bool("d", false, "details")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		nodeName = args[0]
	}
	
	inspect, err := fetchNode(nodeName, *details)
	if err != nil {
		fmt.Printf("Failed due to %v", err)
		os.Exit(1)
	}
	if len(inspect.Error) > 0 {
		fmt.Println(inspect.Error)
		os.Exit(1)
	}
	if len(inspect.Nodes) == 0 {
		fmt.Println("no node in inspct")
		os.Exit(1)
	}
	
	if *details {
		displayDetails(inspect)
		os.Exit(0)
	}
	
	displaySummary(inspect.Nodes)
}

func fetchNode(node string, detail bool) (*types.InspectResult, error) {
	url := "http://127.0.0.1:32743/gputopology-scheduler/inspect"
	if len(node) != 0 {
		url += "/" + node
	}
	if detail {
		url += "?detail=true"
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexepected status code %d", resp.StatusCode)
	}
	rawData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var inspectResult types.InspectResult
	err = json.Unmarshal(rawData, &inspectResult)

	if err != nil {
		return nil, err
	}

	return &inspectResult, nil
}
