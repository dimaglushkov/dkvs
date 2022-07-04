package internal

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dimaglushkov/dkvs/api/storagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	maxKeyLen = 50
	maxValLen = 100
)

func init() {
	rand.Seed(time.Now().Unix())
}

type KeyVal struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

type storageClient struct {
	storagepb.StorageClient
	active bool
}

type handler struct {
	storages      []storageClient
	numOfStorages int
	replicaFactor int
}

func NewHandler(storageLocs []string, rf int) *handler {
	h := handler{}

	// initializing grpc connections
	h.storages = make([]storageClient, 0, len(storageLocs))
	for _, sl := range storageLocs {

		//todo: add proper authentication instead of insecure.NewCredentials()
		conn, err := grpc.Dial(sl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("error while establishing gRPC connection with %s: %s", sl, err)
		}
		client := storagepb.NewStorageClient(conn)

		// checking if this connection actually works
		_, err = client.Get(context.Background(), &storagepb.Key{Key: "1"})
		if err != nil {
			log.Printf("error while initially checking connection storage at %s: %s", sl, err)
			continue
		}
		h.storages = append(h.storages, storageClient{client, true})
	}

	h.replicaFactor = rf
	h.numOfStorages = len(h.storages)

	// simple workaround to avoid redundant duplication
	if h.numOfStorages < h.replicaFactor {
		h.replicaFactor = h.numOfStorages
	}

	return &h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGetRequest(w, r)
	case "PUT":
		h.handlePutRequest(w, r)
	case "POST":
		h.handlePutRequest(w, r)
	case "DELETE":
		h.handleDeleteRequest(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *handler) hashKey(k string) int {
	var s int
	for _, b := range k {
		s += int(b)
	}
	return s % h.numOfStorages
}

func (h *handler) findTargetStorages(key string) []int {
	targetIds := make([]int, 0, h.replicaFactor)
	ts := h.hashKey(key)
	var tmp int

	for i := 0; i < h.replicaFactor; i++ {
		if ts+i < h.numOfStorages {
			tmp = ts + i
		} else {
			tmp = i
		}

		if h.storages[tmp].active {
			targetIds = append(targetIds, tmp)
		}
	}

	// this case means that more than RF nodes went down and all of those nodes were replicating same data
	// which means this data becomes unavailable. This can be solved by either introducing recovering mechanisms
	// for the nodes that went down or by introducing shard layer for the application and making sure every shard
	// is replicated enough times
	if len(targetIds) == 0 {
		log.Fatalf("too many nodes went down")
	}

	return targetIds
}

func (h *handler) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	key := r.RequestURI[1:]
	if key == "" || len(key) > maxKeyLen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ts := h.findTargetStorages(key)
	storageId := ts[rand.Intn(len(ts))]

	ctx := context.Background()

	resp, err := h.storages[storageId].Get(ctx, &storagepb.Key{Key: key})
	if err != nil {
		log.Printf("error while getting %s from %d storage: %s", key, storageId, err)
		h.storages[storageId].active = false
		h.handleGetRequest(w, r)
		return
	}

	if !resp.Success {
		_, err = w.Write([]byte(resp.Value))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	jsonResponse, err := json.Marshal(KeyVal{Key: key, Val: resp.Value})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) handlePutRequest(w http.ResponseWriter, r *http.Request) {
	kv := new(KeyVal)
	err := json.NewDecoder(r.Body).Decode(kv)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if kv.Key == "" || len(kv.Key) > maxKeyLen || len(kv.Val) > maxValLen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ts := h.findTargetStorages(kv.Key)
	var status bool

	for _, storageId := range ts {
		ctx := context.Background()

		resp, err := h.storages[storageId].Put(ctx, &storagepb.KeyValue{Key: kv.Key, Value: kv.Val})
		if err != nil {
			log.Printf("error while putting {%s:%s} to %d storage: %s", kv.Key, kv.Val, storageId, err)
			h.storages[storageId].active = false
		}

		if resp.Success {
			status = true
		}
	}

	if status {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (h *handler) handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
	key := r.RequestURI[1:]
	if key == "" || len(key) > maxKeyLen {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ts := h.findTargetStorages(key)
	var status bool

	for _, storageId := range ts {
		ctx := context.Background()

		resp, err := h.storages[storageId].Delete(ctx, &storagepb.Key{Key: key})
		if err != nil {
			log.Printf("error while deleting %s from %d storage: %s", key, storageId, err)
			h.storages[storageId].active = false
		}
		if resp.Success {
			status = true
		}
	}

	if status {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
