package cluster

import (
	"context"
	"countermag/internal/database"
	"countermag/internal/http/middleware"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func getAddToClusterMessage(source *NodeInfo, dest *NodeInfo, message string) AddToClusterRequest {
	return AddToClusterRequest{
		Source: NodeInfo{
			Id:   source.Id,
			Host: source.Host,
			Port: source.Port,
		},
		Dest: NodeInfo{
			Id:   dest.Id,
			Host: dest.Host,
			Port: dest.Port,
		},
		Message: message,
	}
}

func syncSlaves(
	logger *slog.Logger,
	counterStore *database.Database,
	cluster *Cluster,
	clusterClient *ClusterClient,
) {
	for _, slave := range cluster.Slaves {
		err := clusterClient.PingNode(slave)
		if err != nil {
			logger.Warn("Slave failed healthcheck", "slaveId", slave.Id)
			slave.Status = StatusDown
			continue
		}

		slave.Status = StatusReady
		data := counterStore.Export()

		go func(node *NodeInfo) {
			err := clusterClient.SyncStores(node, data)

			if err != nil {
				logger.Warn("Slave sync failed", "slaveId", node.Id)
				return
			}

			logger.Info("Slave sync succeeded", "slaveId", node.Id)
		}(slave)
	}
}

func connectToCluster(
	logger *slog.Logger,
	clusterClient *ClusterClient,
	src NodeInfo,
	dest NodeInfo,
) bool {

	err := clusterClient.PingNode(&dest)
	if err != nil {
		logger.Debug("No active cluster found.", "nodeId", src.Id)
		return false
	}

	logger.Debug("Found cluster. Requesting connection.")

	err = clusterClient.ConnectToNode(&src, &dest)

	if err != nil {
		logger.Error("Connection to cluster failed.", "nodeId", src.Id)
		return false
	}
	logger.Debug("Connection to cluster succeeded.")

	return true
}

func handleConnectToCluster(
	logger *slog.Logger,
	cluster *Cluster,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var requestMessage AddToClusterRequest
		json.NewDecoder(r.Body).Decode(&requestMessage)

		logger.Info("Got request", "request", requestMessage.String())

		requestMessage.Source.Status = StatusReady
		cluster.Slaves[requestMessage.Source.Port] = &requestMessage.Source

		w.WriteHeader(http.StatusOK)
	})
}

func handlePing() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func handleGetCluster(cluster *Cluster) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cluster)
	})
}

func handleStoreUpdate(logger *slog.Logger, counterStore *database.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	
		var requestMessage map[string]int
		json.NewDecoder(r.Body).Decode(&requestMessage)

		counterStore.Import(requestMessage)
		logger.Info("Updated counter store", "request", requestMessage)

		w.WriteHeader(http.StatusOK)
	})
}

func newHandler(
	logger *slog.Logger,
	counterStore *database.Database,
	cluster *Cluster,
) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/cluster", handleGetCluster(cluster))
	mux.Handle("/store", handleStoreUpdate(logger, counterStore))
	mux.Handle("/connect", handleConnectToCluster(logger, cluster))
	mux.Handle("/ping", handlePing())

	var handler http.Handler = mux

	handler = middleware.AddLogging(logger, handler)

	return handler
}

func RunClusterServer(
	ctx context.Context,
	logger *slog.Logger,
	counterStore *database.Database,
	thisNode NodeInfo,
	master NodeInfo,
) error {
	logger = logger.With("server", "cluster")

	clusterClient := NewClusterClient()
	ableToConnect := connectToCluster(logger, clusterClient, thisNode, master)

	if !ableToConnect {
		logger.Info("No active cluster found. Starting as master")
	}

	cluster := &Cluster{
		Master: thisNode,
		Slaves: make(map[string]*NodeInfo),
	}

	handler := newHandler(
		logger,
		counterStore,
		cluster,
	)

	server := &http.Server{
		Addr:    fmt.Sprint(":" + thisNode.Port),
		Handler: handler,
	}

	go func() {
		logger.Info(fmt.Sprintf("Cluster server listening on %s", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving\n", "error", err)
		}
	}()

	quit := make(chan bool)
	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for {
			select {
			case <-ticker.C:
				syncSlaves(logger, counterStore, cluster, clusterClient)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-ctx.Done()
	quit <- true

	logger.Info("Shutting down cluster server...")

	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)

	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error shutting down cluster server: %s\n", "error", err)
	}

	return nil
}
