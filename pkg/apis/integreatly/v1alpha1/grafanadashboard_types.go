package v1alpha1

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const GrafanaDashboardKind = "GrafanaDashboard"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GrafanaDashboardSpec defines the desired state of GrafanaDashboard
type GrafanaDashboardSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Json         string                       `json:"json"`
	Jsonnet      string                       `json:"jsonnet"`
	Name         string                       `json:"name"`
	Plugins      PluginList                   `json:"plugins,omitempty"`
	Url          string                       `json:"url,omitempty"`
	ConfigMapRef *corev1.ConfigMapKeySelector `json:"configMapRef,omitempty"`
	Datasources  []GrafanaDashboardDatasource `json:"datasources,omitempty"`
}

type GrafanaDashboardDatasource struct {
	InputName      string `json:"inputName"`
	DatasourceName string `json:"datasourceName"`
}

// Used to keep a dashboard reference without having access to the dashboard
// struct itself
type GrafanaDashboardRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
	Hash      string `json:"hash"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaDashboard is the Schema for the grafanadashboards API
// +k8s:openapi-gen=true
type GrafanaDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GrafanaDashboardSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaDashboardList contains a list of GrafanaDashboard
type GrafanaDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GrafanaDashboard `json:"items"`
}

type GrafanaDashboardStatusMessage struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func init() {
	SchemeBuilder.Register(&GrafanaDashboard{}, &GrafanaDashboardList{})
}

func (d *GrafanaDashboard) Hash() string {
	var datasources strings.Builder
	for _, input := range d.Spec.Datasources {
		datasources.WriteString(input.DatasourceName)
		datasources.WriteString(input.InputName)
	}

	hash := sha256.New()
	io.WriteString(hash, d.Spec.Json)
	io.WriteString(hash, d.Spec.Url)
	io.WriteString(hash, d.Spec.Jsonnet)

	if d.Spec.ConfigMapRef != nil {
		io.WriteString(hash, d.Spec.ConfigMapRef.Name)
		io.WriteString(hash, d.Spec.ConfigMapRef.Key)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (d *GrafanaDashboard) UID() string {
	// Use md5 here because Grafana UIDs can only have a maximum length of
	// 40 characters.
	return fmt.Sprintf("%x", sha1.Sum([]byte(d.Namespace+d.Name)))
}
