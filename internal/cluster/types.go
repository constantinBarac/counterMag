package cluster

import "strconv"

type NodeStatus string

const (
	StatusReady NodeStatus = "READY"
	StatusDown  NodeStatus = "DOWN"
)

type Cluster struct {
	Master NodeInfo
	Slaves map[string]*NodeInfo
}

type NodeInfo struct {
	Id     int        `json:"nodeId"`
	Host   string     `json:"nodeIpAddr"`
	Port   string     `json:"port"`
	Status NodeStatus `json:"status"`
}

type AddToClusterRequest struct {
	Source  NodeInfo `json:"source"`
	Dest    NodeInfo `json:"dest"`
	Message string   `json:"message"`
}

func (node NodeInfo) String() string {
	return "NodeInfo:{ nodeId:" + strconv.Itoa(node.Id) + ", nodeIpAddr:" + node.Host + ", port:" + node.Port + " }"
}

func (req AddToClusterRequest) String() string {
	return "AddToClusterMessage:{\n  source:" + req.Source.String() + ",\n  dest: " + req.Dest.String() + ",\n  message:" + req.Message + " }"
}
