package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/cors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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
	Password  string   `json:"password"`
	Databases []string `json:"databases"`
	Storage   string   `json:"storage"`
}

// API Response struct
type CreateClusterResponse struct {
	ClusterName string   `json:"clusterName"`
	User        string   `json:"user"`
	Databases   []string `json:"databases"`
	Storage     string   `json:"storage"`
	Password    string   `json:"password"`
	NodePort    int32    `json:"nodePort"`
}

type GetClusterResponse struct {
	ClusterName string   `json:"clusterName"`
	User        string   `json:"user"`
	Databases   []string `json:"databases"`
	Storage     string   `json:"storage"`
	Password    string   `json:"password"`
	NodePort    int32    `json:"nodePort"`
}

func getDynamicClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func getKubeClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func createPostgresCluster(w http.ResponseWriter, r *http.Request) {
	var req PostgresClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.User == "" || req.Password == "" || len(req.Databases) == 0 || req.Storage == "" {
		http.Error(w, "name, user, password, databases, and storage are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	clientset, err := getKubeClientset()
	if err != nil {
		http.Error(w, "k8s clientset error: "+err.Error(), 500)
		return
	}

	log.Printf("start creating postgresCluster")
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
						"name":     "instance1",
						"replicas": int32(2),
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
		http.Error(w, "k8s dynamic client error: "+err.Error(), 500)
		return
	}

	_, err = client.Resource(gvr).Namespace(pgoNamespace).Create(ctx, cluster, metav1.CreateOptions{})
	if err != nil {
		http.Error(w, "k8s create error: "+err.Error(), 500)
		return
	}

	log.Printf("PostgresCluster %s created successfully", req.Name)

	serviceNames := []string{
		req.Name + "-ha",
		req.Name + "-primary",
	}

	var nodePort int32
	found := false
	for i := 0; i < 20 && !found; i++ {
		for _, serviceName := range serviceNames {
			svc, err := clientset.CoreV1().Services(pgoNamespace).Get(ctx, serviceName, metav1.GetOptions{})
			if err == nil && len(svc.Spec.Ports) > 0 {
				for _, port := range svc.Spec.Ports {
					if port.NodePort != 0 {
						nodePort = port.NodePort
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}

	// Wait for the pguser secret to be created by the operator (up to 30s)
	userSecretName := req.Name + "-pguser-" + req.User
	var userSecret *corev1.Secret
	secretFound := false
	for i := 0; i < 30; i++ {
		userSecret, err = clientset.CoreV1().Secrets(pgoNamespace).Get(ctx, userSecretName, metav1.GetOptions{})
		if err == nil {
			secretFound = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	if !secretFound {
		http.Error(w, "user secret not created by operator in time", 500)
		return
	}

	if userSecret.StringData == nil {
		userSecret.StringData = map[string]string{}
	}
	userSecret.StringData["password"] = req.Password
	_, err = clientset.CoreV1().Secrets(pgoNamespace).Update(ctx, userSecret, metav1.UpdateOptions{})
	if err != nil {
		http.Error(w, "failed to update user password secret: "+err.Error(), 500)
		return
	}
	log.Printf("Updated secret %s with new password", userSecretName)

	resp := CreateClusterResponse{
		ClusterName: req.Name,
		User:        req.User,
		Databases:   req.Databases,
		Storage:     req.Storage,
		Password:    req.Password,
		NodePort:    nodePort,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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

	// Parse fields from the returned object
	spec, _ := obj.Object["spec"].(map[string]interface{})
	users, _ := spec["users"].([]interface{})
	var user string
	var databases []string
	if len(users) > 0 {
		userObj, _ := users[0].(map[string]interface{})
		user, _ = userObj["name"].(string)
		if dbs, ok := userObj["databases"].([]interface{}); ok {
			for _, db := range dbs {
				if dbStr, ok := db.(string); ok {
					databases = append(databases, dbStr)
				}
			}
		}
	}
	storage := ""
	instances, _ := spec["instances"].([]interface{})
	if len(instances) > 0 {
		inst, _ := instances[0].(map[string]interface{})
		if dvc, ok := inst["dataVolumeClaimSpec"].(map[string]interface{}); ok {
			if resources, ok := dvc["resources"].(map[string]interface{}); ok {
				if requests, ok := resources["requests"].(map[string]interface{}); ok {
					storage, _ = requests["storage"].(string)
				}
			}
		}
	}

	// Find NodePort (same logic as in createPostgresCluster)
	clientset, err := getKubeClientset()
	if err != nil {
		http.Error(w, "k8s clientset error: "+err.Error(), 500)
		return
	}
	serviceNames := []string{
		name + "-ha",
		name + "-primary",
	}
	var nodePort int32
	found := false
	for i := 0; i < 5 && !found; i++ {
		for _, serviceName := range serviceNames {
			svc, err := clientset.CoreV1().Services(pgoNamespace).Get(ctx, serviceName, metav1.GetOptions{})
			if err == nil && len(svc.Spec.Ports) > 0 {
				for _, port := range svc.Spec.Ports {
					if port.NodePort != 0 {
						nodePort = port.NodePort
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	resp := GetClusterResponse{
		ClusterName: name,
		User:        user,
		Databases:   databases,
		Storage:     storage,
		NodePort:    nodePort,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func deletePostgresCluster(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name param", 400)
		return
	}
	ctx := context.Background()
	client, err := getDynamicClient()
	if err != nil {
		http.Error(w, "k8s client error: "+err.Error(), 500)
		return
	}
	clientset, err := getKubeClientset()
	if err != nil {
		http.Error(w, "k8s clientset error: "+err.Error(), 500)
		return
	}

	// 1. Delete PostgresCluster CR
	err = client.Resource(gvr).Namespace(pgoNamespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		http.Error(w, "delete error: "+err.Error(), 500)
		return
	}

	// 2. Wait for CR to actually disappear (finalizers, etc), up to 30 seconds
	for i := 0; i < 30; i++ {
		_, err := client.Resource(gvr).Namespace(pgoNamespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 3. Clean up leftover k8s resources by name pattern (if any still exist)
	prefix := name

	// Services
	svcList, _ := clientset.CoreV1().Services(pgoNamespace).List(ctx, metav1.ListOptions{})
	for _, svc := range svcList.Items {
		if strings.HasPrefix(svc.Name, prefix) {
			_ = clientset.CoreV1().Services(pgoNamespace).Delete(ctx, svc.Name, metav1.DeleteOptions{})
		}
	}

	// StatefulSets
	ssList, _ := clientset.AppsV1().StatefulSets(pgoNamespace).List(ctx, metav1.ListOptions{})
	for _, ss := range ssList.Items {
		if strings.HasPrefix(ss.Name, prefix) {
			_ = clientset.AppsV1().StatefulSets(pgoNamespace).Delete(ctx, ss.Name, metav1.DeleteOptions{})
		}
	}

	// Jobs (backups, etc)
	jobList, _ := clientset.BatchV1().Jobs(pgoNamespace).List(ctx, metav1.ListOptions{})
	for _, job := range jobList.Items {
		if strings.HasPrefix(job.Name, prefix) {
			_ = clientset.BatchV1().Jobs(pgoNamespace).Delete(ctx, job.Name, metav1.DeleteOptions{})
		}
	}

	// Pods (should be cleaned up with above, but in case some orphaned)
	podList, _ := clientset.CoreV1().Pods(pgoNamespace).List(ctx, metav1.ListOptions{})
	for _, pod := range podList.Items {
		if strings.HasPrefix(pod.Name, prefix) {
			_ = clientset.CoreV1().Pods(pgoNamespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
		}
	}

	// PVCs (optional, if you want to forcibly remove persistent data)
	pvcList, _ := clientset.CoreV1().PersistentVolumeClaims(pgoNamespace).List(ctx, metav1.ListOptions{})
	for _, pvc := range pvcList.Items {
		if strings.HasPrefix(pvc.Name, prefix) {
			_ = clientset.CoreV1().PersistentVolumeClaims(pgoNamespace).Delete(ctx, pvc.Name, metav1.DeleteOptions{})
		}
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
	mux.HandleFunc("/create", createPostgresCluster)
	mux.HandleFunc("/get", getPostgresCluster)
	mux.HandleFunc("/delete", deletePostgresCluster)
	mux.HandleFunc("/list", listPostgresClusters)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	handler := cors.AllowAll().Handler(mux)
	log.Println("pgo-go-backend listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
