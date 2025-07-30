package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	pgoNamespace = "postgres-operator"
)

var (
	gvr = schema.GroupVersionResource{
		Group:    "postgres-operator.crunchydata.com",
		Version:  "v1beta1",
		Resource: "postgresclusters",
	}
)

// API Request struct
type PostgresClusterRequest struct {
	Name      string   `json:"name"`
	User      string   `json:"user"`
	Databases []string `json:"databases"`
	Storage   string   `json:"storage"` // e.g. "200Mi", "1Gi", "1.5Gi"
}

// Create a dynamic Kubernetes client
func getDynamicClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func createPostgresCluster(w http.ResponseWriter, r *http.Request) {
	var req PostgresClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.User == "" || len(req.Databases) == 0 || req.Storage == "" {
		http.Error(w, "name, user, databases, and storage are required", http.StatusBadRequest)
		return
	}

	// Build PostgresCluster object (unstructured)
	cluster := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "postgres-operator.crunchydata.com/v1beta1",
			"kind":       "PostgresCluster",
			"metadata": map[string]interface{}{
				"name": req.Name,
				"annotations": map[string]interface{}{
					"postgres-operator.crunchydata.com/autoCreateUserSchema": "true",
				},
			},
			"spec": map[string]interface{}{
				"postgresVersion": 17,
				"users": []interface{}{
					map[string]interface{}{
						"name":      req.User,
						"databases": req.Databases,
					},
				},
				"instances": []interface{}{
					map[string]interface{}{
						"name":    "instance1",
						"replica": 2,
						"dataVolumeClaimSpec": map[string]interface{}{
							"accessModes": []interface{}{"ReadWriteOnce"},
							"resources": map[string]interface{}{
								"requests": map[string]interface{}{
									"storage": req.Storage,
								},
							},
						},
					},
				},
				"backups": map[string]interface{}{
					"pgbackrest": map[string]interface{}{
						"repos": []interface{}{
							map[string]interface{}{
								"name": "repo1",
								"volume": map[string]interface{}{
									"volumeClaimSpec": map[string]interface{}{
										"accessModes": []interface{}{"ReadWriteOnce"},
										"resources": map[string]interface{}{
											"requests": map[string]interface{}{
												"storage": req.Storage,
											},
										},
									},
								},
							},
						},
					},
				},
				"service": map[string]interface{}{
					"type": "NodePort",
				},
			},
		},
	}

	client, err := getDynamicClient()
	if err != nil {
		http.Error(w, "k8s client error: "+err.Error(), 500)
		return
	}
	ctx := context.Background()

	// Create
	created, err := client.Resource(gvr).Namespace(pgoNamespace).Create(ctx, cluster, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "k8s create error: "+err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created.Object)
}

func getPostgresCluster(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name param", 400)
		return
	}
	client, err := getDynamicClient()
	if err != nil {
		http.Error(w, "k8s client error: "+err.Error(), 500)
		return
	}
	ctx := context.Background()
	obj, err := client.Resource(gvr).Namespace(pgoNamespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		http.Error(w, "not found: "+err.Error(), 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj.Object)
}

func deletePostgresCluster(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name param", 400)
		return
	}
	client, err := getDynamicClient()
	if err != nil {
		http.Error(w, "k8s client error: "+err.Error(), 500)
		return
	}
	ctx := context.Background()
	err = client.Resource(gvr).Namespace(pgoNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "delete error: "+err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func listPostgresClusters(w http.ResponseWriter, r *http.Request) {
	client, err := getDynamicClient()
	if err != nil {
		http.Error(w, "k8s client error: "+err.Error(), 500)
		return
	}
	ctx := context.Background()
	list, err := client.Resource(gvr).Namespace(pgoNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		http.Error(w, "list error: "+err.Error(), 500)
		return
	}

	// Extract the "items" array from UnstructuredList
	items := list.Items
	clusterNames := []string{}
	for _, item := range items {
		clusterNames = append(clusterNames, item.GetName())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clusterNames)
}

func main() {
	mux := http.NewServeMux()
	// http.HandleFunc("/create", createPostgresCluster)
	// http.HandleFunc("/get", getPostgresCluster)
	// http.HandleFunc("/delete", deletePostgresCluster)
	// http.HandleFunc("/list", listPostgresClusters)

	mux.HandleFunc("/create", createPostgresCluster)
	mux.HandleFunc("/get", getPostgresCluster)
	mux.HandleFunc("/delete", deletePostgresCluster)
	mux.HandleFunc("/list", listPostgresClusters)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handler := cors.AllowAll().Handler(mux) // allow all origins for simplicity

	log.Println("pgo-go-backend listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
