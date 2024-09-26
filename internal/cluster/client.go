package cluster

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

type ClusterClient struct {
	client *http.Client
}

func NewClusterClient() *ClusterClient {
	httpClient := http.Client{}

	httpClient.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &ClusterClient{
		client: &httpClient,
	}
}

func (c *ClusterClient) PingNode(node *NodeInfo) error {

	pingEndpoint := "http://" + node.Host + ":" + node.Port + "/ping"

	req, _ := http.NewRequest(http.MethodGet, pingEndpoint, nil)

	res, err := c.client.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (c *ClusterClient) ConnectToNode(thisNode *NodeInfo, masterNode *NodeInfo) error {
	connectEndpoint := "http://" + masterNode.Host + ":" + masterNode.Port + "/connect"

	message := getAddToClusterMessage(thisNode, masterNode, "")
	messageBytes, _ := json.Marshal(message)
	body := bytes.NewBuffer(messageBytes)

	req, _ := http.NewRequest(http.MethodPost, connectEndpoint, body)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func (c *ClusterClient) SyncStores(node *NodeInfo, data map[string]int) error {
	syncEndpoint := "http://" + node.Host + ":" + node.Port + "/store"

	messageBytes, _ := json.Marshal(data)
	body := bytes.NewBuffer(messageBytes)

	request, _ := http.NewRequest(http.MethodPut, syncEndpoint, body)
	request.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(request)
	if err != nil || res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
