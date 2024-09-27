package main

import (
	"context"
	"countermag/internal/analysis"
	"countermag/internal/cluster"
	"countermag/internal/database"
	"countermag/internal/logging"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func getNodePair(clusterMasterAddr string, clusterPort int) (cluster.NodeInfo, cluster.NodeInfo) {
	myHost := "127.0.0.1"
	thisPort := strconv.Itoa(clusterPort)
	thisNode := cluster.NodeInfo{
		Id:     clusterPort,
		Host:   myHost,
		Port:   thisPort,
		Status: cluster.StatusReady,
	}

	masterHost := strings.Split(clusterMasterAddr, ":")[0]
	masterPort := strings.Split(clusterMasterAddr, ":")[1]
	clusterMaster := cluster.NodeInfo{
		Id:     -1,
		Host:   masterHost,
		Port:   masterPort,
		Status: cluster.StatusReady,
	}

	return thisNode, clusterMaster
}

func main() {
	clusterMasterAddr := flag.String("cluster", "127.0.0.1:8000", "cluster master node address")
	clusterPort := flag.Int("port", 8000, "ip address to run this node on. default is 8000")
	flag.Parse()

	appPort := *clusterPort + 100

	thisNode, clusterMaster := getNodePair(*clusterMasterAddr, *clusterPort)

	logger := logging.GetLogger("local")
	snapshotPath := fmt.Sprintf("/usr/local/countermag/counter-%d.log", thisNode.Id)
	logger.Debug("Snapshot path", "path", snapshotPath)

	ctx := context.Background()
	signalCtx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	
	counterStore := database.NewDatabase(signalCtx, logger, &database.FileSnapshotPersister{Path: snapshotPath})

	go func() {
		if err := cluster.RunClusterServer(
			signalCtx,
			logger,
			counterStore,
			thisNode,
			clusterMaster,
		); err != nil {
			fmt.Fprintf(os.Stderr, "Error running cluster server: %s\n", err)
			os.Exit(1)
		}
	}()

	analysis.RunAnalysisServer(
		signalCtx,
		logger,
		counterStore,
		appPort,
	)

	closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	counterStore.Close(closeCtx)
	defer cancel()
}
