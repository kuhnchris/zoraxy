package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	plugin "github.com/kuhnchris/zoraxy_docker_host_ip_resolve/mod/zoraxy_plugin"
)

const (
	PLUGIN_ID = "eu.kuhnchris.zoraxy_docker_host_ip_resolve"
	UI_PATH   = "/plugins/dockerhostipresolve"
	WEB_ROOT  = "/www"
)

//go:embed www/*
var content embed.FS
var pluginCfg plugin.ConfigureSpec

func callAPIEndpoint(cfg *plugin.ConfigureSpec, apiURL string) (*http.Response, error) {
	// Make an API call to the permitted endpoint
	client := &http.Client{}
	//apiCallURL := fmt.Sprintf("http://127.0.0.1:%d/plugin/"+apiURL, cfg.ZoraxyPort)
	apiCallURL := fmt.Sprintf("http://127.0.0.1:%d/plugin/"+apiURL, 8443)
	req, err := http.NewRequest(http.MethodGet, apiCallURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	// Make sure to set the Authorization header
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey) // Use the API key from the runtime config
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making API call: %v", err)
	}
	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {

		return nil, fmt.Errorf("error dumping response: %v", err)
	}

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status %d - %s", resp.StatusCode, string(respDump))
	}

	return resp, nil
}

func checkAllUpstreams() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("Calling endpoint...\n")
		resp, err := callAPIEndpoint(&pluginCfg, "api/proxy/list?type=host")
		if err != nil {
			fmt.Printf("error: %s", err)
			//panic(err)
			continue
		}

		//jdec := json.NewDecoder(resp.Body)
		// Unmarshal into a map for dynamic handling
		//var data []map[string]interface{}

		retResp, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON: %v", err)
			continue

		}
		/*
			err = json.Unmarshal(retResp, &data)
			if err != nil {
				fmt.Printf("Error unmarshalling JSON: %v", err)
				continue
			}*/

		var pConfigs ProxyConfigs
		err = json.Unmarshal(retResp, &pConfigs)
		if err != nil {
			fmt.Printf("cannot unmarshal proxy config: %s\n", err)
			continue
		}

		// Iterate over the map to process dynamic keys and values
		/*for key, value := range data {
			fmt.Printf("Key: %s, Value: %v\n", key, value)
		}*/
		for _, entry := range pConfigs {
			resp, err := callAPIEndpoint(&pluginCfg, "api/proxy/upstream/list?ep="+*entry.RootOrMatchingDomain)
			if err != nil {
				fmt.Printf("error: %s", err)
				//panic(err)
				continue
			}
			fmt.Printf("Upstream reply: %s", resp)

			//entry["RootOrMatchingDomain"]
		}
	}
}

func main() {
	// Serve the plugin intro spect
	// This will print the plugin intro spect and exit if the -introspect flag is provided
	runtimeCfg, err := plugin.ServeAndRecvSpec(&plugin.IntroSpect{
		ID:            PLUGIN_ID,
		Name:          "Docker Hostname-to-IP resolver",
		Author:        "KuhnChris",
		AuthorContact: "kuhnchris@users.github.com",
		Description:   "A plugin that tries to resolve the given hostname to a docker currently running docker container IP",
		URL:           "https://github.com/kuhnchris",
		Type:          plugin.PluginType_Utilities,
		VersionMajor:  1,
		VersionMinor:  0,
		VersionPatch:  0,

		UIPath: UI_PATH,
		/* API Access Control */
		PermittedAPIEndpoints: []plugin.PermittedAPIEndpoint{
			{
				Method:   http.MethodGet,
				Endpoint: "/plugin/api/proxy/list",
				Reason:   "Used to display all configured Access Rules",
			},
		},
	})
	if err != nil {
		//Terminate or enter standalone mode here
		panic(err)
	}
	pluginCfg = *runtimeCfg

	embedWebRouter := plugin.NewPluginEmbedUIRouter(PLUGIN_ID, &content, WEB_ROOT, UI_PATH)
	embedWebRouter.RegisterTerminateHandler(func() {
		// Do cleanup here if needed
		fmt.Println("Docker H2IP handler terminated")
	}, nil)

	fmt.Println("We should use apiKey '" + runtimeCfg.APIKey + "' as key against localhost:" + strconv.Itoa(runtimeCfg.ZoraxyPort))

	go checkAllUpstreams()

	// Serve the hello world page in the www folder
	http.Handle(UI_PATH, embedWebRouter.Handler())
	fmt.Println("Docker IP resolve watcher started at http://127.0.0.1:" + strconv.Itoa(runtimeCfg.Port))
	err = http.ListenAndServe("127.0.0.1:"+strconv.Itoa(runtimeCfg.Port), nil)
	if err != nil {
		panic(err)
	}

}
