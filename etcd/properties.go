package etcd

import (
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

type Properties struct {
	Etcd      *PropertiesEtcd           `yaml:"etcd,omitempty"`
	Consul    *consul.PropertiesConsul  `yaml:"consul,omitempty"`
	WardenCPI *iaas.PropertiesWardenCPI `yaml:"warden_cpi,omitempty"`
	AWS       *iaas.PropertiesAWS       `yaml:"aws,omitempty"`
	Registry  *core.PropertiesRegistry  `yaml:"registry,omitempty"`
	Blobstore *core.PropertiesBlobstore `yaml:"blobstore,omitempty"`
	Agent     *core.PropertiesAgent     `yaml:"agent,omitempty"`
}

type PropertiesEtcd struct {
	Cluster                         []PropertiesEtcdCluster `yaml:"cluster"`
	Machines                        []string                `yaml:"machines"`
	PeerRequireSSL                  bool                    `yaml:"peer_require_ssl"`
	RequireSSL                      bool                    `yaml:"require_ssl"`
	HeartbeatIntervalInMilliseconds int                     `yaml:"heartbeat_interval_in_milliseconds"`
	AdvertiseURLsDNSSuffix          string                  `yaml:"advertise_urls_dns_suffix"`
	CACert                          string                  `yaml:"ca_cert"`
	ClientCert                      string                  `yaml:"client_cert"`
	ClientKey                       string                  `yaml:"client_key"`
	PeerCACert                      string                  `yaml:"peer_ca_cert"`
	PeerCert                        string                  `yaml:"peer_cert"`
	PeerKey                         string                  `yaml:"peer_key"`
	ServerCert                      string                  `yaml:"server_cert"`
	ServerKey                       string                  `yaml:"server_key"`
}

type PropertiesEtcdCluster struct {
	Instances int    `yaml:"instances"`
	Name      string `yaml:"name"`
}