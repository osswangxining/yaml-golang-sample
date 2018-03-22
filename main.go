package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigMap struct fields must be public in order for unmarshal to
// correctly populate the data.
type ConfigMap struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	}
	Data struct {
		CloudantDBName                 string `yaml:"cloudantDBName"`
		CloudantUserName               string `yaml:"cloudantUserName"`
		ObjectStorageEndPoint          string `yaml:"objectStorageEndPoint"`
		ObjectStorageIBMAuthEndPoint   string `yaml:"objectStorageIBMAuthEndPoint"`
		ObjectStorageServiceInstanceID string `yaml:"objectStorageServiceInstanceID"`
		WdsEndPoint                    string `yaml:"wdsEndPoint"`
		WdsEnvID                       string `yaml:"wdsEnvID"`
		WdsURL                         string `yaml:"wdsURL"`
		WdsUserName                    string `yaml:"wdsUserName"`
		WdsVersionDate                 string `yaml:"wdsVersionDate"`
		ServiceNameAPIPrivate          string `yaml:"serviceNameAPIPrivate"`
		ServicePortAPIPrivate          string `yaml:"servicePortAPIPrivate"`
		ServicePrctolAPIPrivate        string `yaml:"servicePrctolAPIPrivate"`
		AppIDAdminTenantID             string `yaml:"appIDAdminTenantID"`
		AppIDAdminClientID             string `yaml:"appIDAdminClientID"`
		AppIDAdminOAuthServerURL       string `yaml:"appIDAdminOAuthServerURL"`
		AppIDAdminRedirectURL          string `yaml:"appIDAdminRedirectURL"`
		AppIDUserTenantID              string `yaml:"appIDUserTenantID"`
		AppIDUserClientID              string `yaml:"appIDUserClientID"`
		AppIDUserOAuthServerURL        string `yaml:"appIDUserOAuthServerURL"`
		AppIDUserRedirectURL           string `yaml:"appIDUserRedirectURL"`
		APIPrivateContextPath          string `yaml:"apiPrivateContextPath"`
		APIPublicContextPath           string `yaml:"apiPublicContextPath"`
		UIAdminContextPath             string `yaml:"uiAdminContextPath"`
		UIEndUserContextPath           string `yaml:"uiEndUserContextPath"`
		WdsCrawlerContextPath          string `yaml:"wdsCrawlerContextPath"`
		JWTExpires                     string `yaml:"jwtExpires"`
	}
}

// Secret struct fields must be public in order for unmarshal to
// correctly populate the data.
type Secret struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	}
	Type string `yaml:"type"`
	Data struct {
		CloudantPassword    string `yaml:"cloudantPassword"`
		ObjectStorageAPIKey string `yaml:"objectStorageAPIKey"`
		WdsPassword         string `yaml:"wdsPassword"`
		AppIDAdminSecret    string `yaml:"appIDAdminSecret"`
		AppIDUserSecret     string `yaml:"appIDUserSecret"`
		JWTTokenSecret      string `yaml:"jwtTokenSecret"`
	}
}

// Ingress struct fields must be public in order for unmarshal to
// correctly populate the data.
type Ingress struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Annotations struct {
			RewritePath string `yaml:"ingress.bluemix.net/rewrite-path"`
		}
	}
	Spec struct {
		TLS []struct {
			Hosts      []string
			SecretName string `yaml:"secretName"`
		}
		Rules []struct {
			Host string
			HTTP struct {
				Paths []struct {
					Path    string `yaml:"path"`
					Backend struct {
						ServiceName string `yaml:"serviceName"`
						ServicePort int    `yaml:"servicePort"`
					}
				}
			}
		}
	}
}

func main() {
	filename := os.Args[1]
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(source, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- kind:\n%v\n\n", m["kind"])

	// d, err = yaml.Marshal(&m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	//fmt.Printf("--- m dump:\n%s\n\n", string(d))
	if m["kind"] == "ConfigMap" {
		d := initConfigMap(source)
		outputfilename := os.Args[2]
		ioutil.WriteFile(outputfilename, d, 0755)
	} else if m["kind"] == "Secret" {
		d := initSecret(source)
		outputfilename := os.Args[2]
		ioutil.WriteFile(outputfilename, d, 0755)
	} else if m["kind"] == "Ingress" {
		d := initIngress(source)
		outputfilename := os.Args[2]
		ioutil.WriteFile(outputfilename, d, 0755)
	} else {
		fmt.Printf("--- the kind should be one of ConfigMap | Secret | Ingress:\n%v\n\n", m["kind"])
	}
}

func initIngress(source []byte) []byte {
	t := Ingress{}
	err := yaml.Unmarshal(source, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if os.Getenv("ingressHost") != "" {
		t.Spec.TLS[0].Hosts[0] = os.Getenv("ingressHost")
		t.Spec.Rules[0].Host = os.Getenv("ingressHost")
	}
	if os.Getenv("secretName") != "" {
		t.Spec.TLS[0].SecretName = os.Getenv("secretName")
	}

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- t dump:\n%s\n\n", string(d))

	return d
}

func initSecret(source []byte) []byte {
	t := Secret{}
	err := yaml.Unmarshal(source, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if os.Getenv("secretMetadataName") != "" {
		t.Metadata.Name = os.Getenv("secretMetadataName")
	}
	if os.Getenv("secretMetadataNamespace") != "" {
		t.Metadata.Namespace = os.Getenv("secretMetadataNamespace")
	}
	if os.Getenv("secretType") != "" {
		t.Type = os.Getenv("secretType")
	}

	if os.Getenv("cloudantPassword") != "" {
		t.Data.CloudantPassword = os.Getenv("cloudantPassword")
	}
	if os.Getenv("objectStorageAPIKey") != "" {
		t.Data.ObjectStorageAPIKey = os.Getenv("objectStorageAPIKey")
	}
	if os.Getenv("wdsPassword") != "" {
		t.Data.WdsPassword = os.Getenv("wdsPassword")
	}
	if os.Getenv("jwtTokenSecret") != "" {
		t.Data.JWTTokenSecret = os.Getenv("jwtTokenSecret")
	}
	if os.Getenv("appIDAdminSecret") != "" {
		t.Data.AppIDAdminSecret = os.Getenv("appIDAdminSecret")
	}
	if os.Getenv("appIDUserSecret") != "" {
		t.Data.AppIDUserSecret = os.Getenv("appIDUserSecret")
	}
	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- t dump:\n%s\n\n", string(d))

	return d
}

func initConfigMap(source []byte) []byte {
	t := ConfigMap{}
	err := yaml.Unmarshal(source, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- t:\n%+v\n\n", t)

	if os.Getenv("configmapMetadataName") != "" {
		t.Metadata.Name = "wea-config"
	}
	if os.Getenv("configmapMetadataNamespace") != "" {
		t.Metadata.Namespace = "default"
	}

	if os.Getenv("cloudantDBName") != "" {
		t.Data.CloudantDBName = os.Getenv("cloudantDBName")
	}
	if os.Getenv("cloudantUserName") != "" {
		t.Data.CloudantUserName = os.Getenv("cloudantUserName")
	}
	if os.Getenv("objectStorageEndPoint") != "" {
		t.Data.ObjectStorageEndPoint = os.Getenv("objectStorageEndPoint")
	}
	if os.Getenv("objectStorageIBMAuthEndPoint") != "" {
		t.Data.ObjectStorageIBMAuthEndPoint = os.Getenv("objectStorageIBMAuthEndPoint")
	}
	if os.Getenv("objectStorageServiceInstanceID") != "" {
		t.Data.ObjectStorageServiceInstanceID = os.Getenv("objectStorageServiceInstanceID")
	}

	if os.Getenv("wdsEndPoint") != "" {
		t.Data.WdsEndPoint = os.Getenv("wdsEndPoint")
	}
	if os.Getenv("wdsEnvID") != "" {
		t.Data.WdsEnvID = os.Getenv("wdsEnvID")
	}
	if os.Getenv("wdsURL") != "" {
		t.Data.WdsURL = os.Getenv("wdsURL")
	}
	if os.Getenv("wdsUserName") != "" {
		t.Data.WdsUserName = os.Getenv("wdsUserName")
	}
	if os.Getenv("wdsVersionDate") != "" {
		t.Data.WdsVersionDate = os.Getenv("wdsVersionDate")
	}

	if os.Getenv("serviceNameAPIPrivate") != "" {
		t.Data.ServiceNameAPIPrivate = os.Getenv("serviceNameAPIPrivate")
	}
	if os.Getenv("servicePortAPIPrivate") != "" {
		t.Data.ServicePortAPIPrivate = os.Getenv("servicePortAPIPrivate")
	}
	if os.Getenv("servicePrctolAPIPrivate") != "" {
		t.Data.ServicePrctolAPIPrivate = os.Getenv("servicePrctolAPIPrivate")
	}

	if os.Getenv("appIDAdminTenantID") != "" {
		t.Data.AppIDAdminTenantID = os.Getenv("appIDAdminTenantID")
	}
	if os.Getenv("appIDAdminClientID") != "" {
		t.Data.AppIDAdminClientID = os.Getenv("appIDAdminClientID")
	}

	if os.Getenv("appIDAdminOAuthServerURL") != "" {
		t.Data.AppIDAdminOAuthServerURL = os.Getenv("appIDAdminOAuthServerURL")
	}
	if os.Getenv("appIDAdminRedirectURL") != "" {
		t.Data.AppIDAdminRedirectURL = os.Getenv("appIDAdminRedirectURL")
	}

	if os.Getenv("appIDUserTenantID") != "" {
		t.Data.AppIDUserTenantID = os.Getenv("appIDUserTenantID")
	}
	if os.Getenv("appIDUserClientID") != "" {
		t.Data.AppIDUserClientID = os.Getenv("appIDUserClientID")
	}

	if os.Getenv("appIDUserOAuthServerURL") != "" {
		t.Data.AppIDUserOAuthServerURL = os.Getenv("appIDUserOAuthServerURL")
	}
	if os.Getenv("appIDUserRedirectURL") != "" {
		t.Data.AppIDUserRedirectURL = os.Getenv("appIDUserRedirectURL")
	}

	if os.Getenv("apiPrivateContextPath") != "" {
		t.Data.APIPrivateContextPath = os.Getenv("apiPrivateContextPath")
	}
	if os.Getenv("apiPublicContextPath") != "" {
		t.Data.APIPublicContextPath = os.Getenv("apiPublicContextPath")
	}
	if os.Getenv("uiAdminContextPath") != "" {
		t.Data.UIAdminContextPath = os.Getenv("uiAdminContextPath")
	}
	if os.Getenv("uiEndUserContextPath") != "" {
		t.Data.UIEndUserContextPath = os.Getenv("uiEndUserContextPath")
	}
	if os.Getenv("wdsCrawlerContextPath") != "" {
		t.Data.WdsCrawlerContextPath = os.Getenv("wdsCrawlerContextPath")
	}
	if os.Getenv("jwtExpires") != "" {
		t.Data.JWTExpires = os.Getenv("jwtExpires")
	}

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- t dump:\n%s\n\n", string(d))

	return d
}
