package turbulence_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/turbulence"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("Manifest", func() {
	Describe("NewManifest", func() {
		It("generates a valid Turbulence AWS manifest", func() {
			manifest := turbulence.NewManifest(turbulence.Config{
				Name:         "turbulence",
				DirectorUUID: "some-director-uuid",
				IPRange:      "10.0.16.0/24",
				BOSH: turbulence.ConfigBOSH{
					Target:         "some-bosh-target",
					Username:       "some-bosh-username",
					Password:       "some-bosh-password",
					DirectorCACert: "some-ca-cert",
				},
			}, iaas.AWSConfig{
				AccessKeyID:           "some-access-key-id",
				SecretAccessKey:       "some-secret-access-key",
				DefaultKeyName:        "some-default-key-name",
				DefaultSecurityGroups: []string{"some-default-security-group1"},
				Region:                "some-region",
				Subnet:                "subnet-1234",
				RegistryHost:          "some-registry-host",
				RegistryPassword:      "some-registry-password",
				RegistryPort:          25777,
				RegistryUsername:      "some-registry-username",
			})

			Expect(manifest).To(Equal(turbulence.Manifest{
				DirectorUUID: "some-director-uuid",
				Name:         "turbulence",
				Releases: []core.Release{
					{
						Name:    "turbulence",
						Version: "latest",
					},
					{
						Name:    "bosh-aws-cpi",
						Version: "latest",
					},
				},
				ResourcePools: []core.ResourcePool{
					{
						Name:    "turbulence",
						Network: "turbulence",
						Stemcell: core.ResourcePoolStemcell{
							Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
							Version: "latest",
						},
						CloudProperties: core.ResourcePoolCloudProperties{
							InstanceType:     "m3.medium",
							AvailabilityZone: "us-east-1a",
							EphemeralDisk: &core.ResourcePoolCloudPropertiesEphemeralDisk{
								Size: 1024,
								Type: "gp2",
							},
						},
					},
				},
				Compilation: core.Compilation{
					Network:             "turbulence",
					ReuseCompilationVMs: true,
					Workers:             3,
					CloudProperties: core.CompilationCloudProperties{
						InstanceType:     "m3.medium",
						AvailabilityZone: "us-east-1a",
						EphemeralDisk: &core.CompilationCloudPropertiesEphemeralDisk{
							Size: 1024,
							Type: "gp2",
						},
					},
				},
				Update: core.Update{
					Canaries:        1,
					CanaryWatchTime: "1000-180000",
					MaxInFlight:     1,
					Serial:          true,
					UpdateWatchTime: "1000-180000",
				},
				Jobs: []core.Job{
					{
						Instances: 1,
						Name:      "api",
						Networks: []core.JobNetwork{
							{
								Name:      "turbulence",
								StaticIPs: []string{"10.0.16.12"},
							},
						},
						PersistentDisk: 1024,
						ResourcePool:   "turbulence",
						Templates: []core.JobTemplate{
							{
								Name:    "turbulence_api",
								Release: "turbulence",
							},
							{
								Name:    "aws_cpi",
								Release: "bosh-aws-cpi",
							},
						},
					},
				},
				Networks: []core.Network{
					{
						Name: "turbulence",
						Subnets: []core.NetworkSubnet{
							{
								CloudProperties: core.NetworkSubnetCloudProperties{
									Subnet: "subnet-1234",
								},
								Gateway: "10.0.16.1",
								Range:   "10.0.16.0/24",
								Reserved: []string{
									"10.0.16.2-10.0.16.11",
									"10.0.16.17-10.0.16.254",
								},
								Static: []string{
									"10.0.16.12",
									"10.0.16.13",
								},
							},
						},
						Type: "manual",
					},
				},
				Properties: turbulence.Properties{
					TurbulenceAPI: &turbulence.PropertiesTurbulenceAPI{
						Certificate: turbulence.APICertificate,
						CPIJobName:  "aws_cpi",
						Director: turbulence.PropertiesTurbulenceAPIDirector{
							CACert:   "some-ca-cert",
							Host:     "some-bosh-target",
							Password: "some-bosh-password",
							Username: "some-bosh-username",
						},
						Password:   "turbulence-password",
						PrivateKey: turbulence.APIPrivateKey,
					},
					AWS: &iaas.PropertiesAWS{
						AccessKeyID:           "some-access-key-id",
						DefaultKeyName:        "some-default-key-name",
						DefaultSecurityGroups: []string{"some-default-security-group1"},
						Region:                "some-region",
						SecretAccessKey:       "some-secret-access-key",
					},
					Registry: &core.PropertiesRegistry{
						Host:     "some-registry-host",
						Password: "some-registry-password",
						Port:     25777,
						Username: "some-registry-username",
					},
					Blobstore: &core.PropertiesBlobstore{
						Address: "10.0.16.12",
						Port:    2520,
						Agent: core.PropertiesBlobstoreAgent{
							User:     "agent",
							Password: "agent-password",
						},
					},
					Agent: &core.PropertiesAgent{
						Mbus: "nats://nats:password@10.0.16.12:4222",
					},
				},
			}))
		})

		It("generates a valid Turbulence BOSH-Lite manifest", func() {
			manifest := turbulence.NewManifest(turbulence.Config{
				DirectorUUID: "some-director-uuid",
				IPRange:      "10.244.4.0/24",
				BOSH: turbulence.ConfigBOSH{
					Target:   "some-bosh-target",
					Username: "some-bosh-username",
					Password: "some-bosh-password",
				},
				Name: "turbulence",
			}, iaas.NewWardenConfig())

			Expect(manifest).To(Equal(turbulence.Manifest{
				DirectorUUID: "some-director-uuid",
				Name:         "turbulence",
				Releases: []core.Release{
					{
						Name:    "turbulence",
						Version: "latest",
					},
					{
						Name:    "bosh-warden-cpi",
						Version: "latest",
					},
				},
				ResourcePools: []core.ResourcePool{
					{
						Name:    "turbulence",
						Network: "turbulence",
						Stemcell: core.ResourcePoolStemcell{
							Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
							Version: "latest",
						},
					},
				},
				Compilation: core.Compilation{
					Network:             "turbulence",
					ReuseCompilationVMs: true,
					Workers:             3,
				},
				Update: core.Update{
					Canaries:        1,
					CanaryWatchTime: "1000-180000",
					MaxInFlight:     1,
					Serial:          true,
					UpdateWatchTime: "1000-180000",
				},
				Jobs: []core.Job{
					{
						Instances: 1,
						Name:      "api",
						Networks: []core.JobNetwork{
							{
								Name:      "turbulence",
								StaticIPs: []string{"10.244.4.12"},
							},
						},
						PersistentDisk: 1024,
						ResourcePool:   "turbulence",
						Templates: []core.JobTemplate{
							{
								Name:    "turbulence_api",
								Release: "turbulence",
							},
							{
								Name:    "warden_cpi",
								Release: "bosh-warden-cpi",
							},
						},
					},
				},
				Networks: []core.Network{
					{
						Name: "turbulence",
						Subnets: []core.NetworkSubnet{
							{
								CloudProperties: core.NetworkSubnetCloudProperties{
									Name: "random",
								},
								Gateway: "10.244.4.1",
								Range:   "10.244.4.0/24",
								Reserved: []string{
									"10.244.4.2-10.244.4.11",
									"10.244.4.17-10.244.4.254",
								},
								Static: []string{
									"10.244.4.12",
									"10.244.4.13",
								},
							},
						},
						Type: "manual",
					},
				},
				Properties: turbulence.Properties{
					TurbulenceAPI: &turbulence.PropertiesTurbulenceAPI{
						Certificate: turbulence.APICertificate,
						CPIJobName:  "warden_cpi",
						Director: turbulence.PropertiesTurbulenceAPIDirector{
							CACert:   turbulence.APIDirectorCACert,
							Host:     "some-bosh-target",
							Password: "some-bosh-password",
							Username: "some-bosh-username",
						},
						Password:   "turbulence-password",
						PrivateKey: turbulence.APIPrivateKey,
					},
					WardenCPI: &iaas.PropertiesWardenCPI{
						Agent: iaas.PropertiesWardenCPIAgent{
							Blobstore: iaas.PropertiesWardenCPIAgentBlobstore{
								Options: iaas.PropertiesWardenCPIAgentBlobstoreOptions{
									Endpoint: "http://10.254.50.4:25251",
									Password: "agent-password",
									User:     "agent",
								},
								Provider: "dav",
							},
							Mbus: "nats://nats:nats-password@10.254.50.4:4222",
						},
						Warden: iaas.PropertiesWardenCPIWarden{
							ConnectAddress: "10.254.50.4:7777",
							ConnectNetwork: "tcp",
						},
					},
				},
			}))
		})
	})

	Describe("FromYAML", func() {
		It("returns a Manifest matching the given YAML", func() {
			turbulenceManifest, err := ioutil.ReadFile("fixtures/turbulence_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest, err := turbulence.FromYAML(turbulenceManifest)
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("turbulence"))
			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "turbulence",
				Version: "latest",
			}))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "bosh-warden-cpi",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "turbulence",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "turbulence",
				Network: "turbulence",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(1))
			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "api",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "turbulence",
					StaticIPs: []string{"10.244.4.12"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "turbulence",
				Templates: []core.JobTemplate{
					{
						Name:    "turbulence_api",
						Release: "turbulence",
					},
					{
						Name:    "warden_cpi",
						Release: "bosh-warden-cpi",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "turbulence",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.11",
							"10.244.4.17-10.244.4.254",
						},
						Static: []string{
							"10.244.4.12",
							"10.244.4.13",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties).To(Equal(turbulence.Properties{
				WardenCPI: &iaas.PropertiesWardenCPI{
					Agent: iaas.PropertiesWardenCPIAgent{
						Blobstore: iaas.PropertiesWardenCPIAgentBlobstore{
							Options: iaas.PropertiesWardenCPIAgentBlobstoreOptions{
								Endpoint: "http://10.254.50.4:25251",
								Password: "agent-password",
								User:     "agent",
							},
							Provider: "dav",
						},
						Mbus: "nats://nats:nats-password@10.254.50.4:4222",
					},
					Warden: iaas.PropertiesWardenCPIWarden{
						ConnectAddress: "10.254.50.4:7777",
						ConnectNetwork: "tcp",
					},
				},
				TurbulenceAPI: &turbulence.PropertiesTurbulenceAPI{
					Certificate: turbulence.APICertificate,
					CPIJobName:  "warden_cpi",
					Director: turbulence.PropertiesTurbulenceAPIDirector{
						CACert:   turbulence.APIDirectorCACert,
						Host:     "some-bosh-target",
						Password: "some-bosh-password",
						Username: "some-bosh-username",
					},
					Password:   "turbulence-password",
					PrivateKey: turbulence.APIPrivateKey,
				},
			}))
		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the turbulence manifest", func() {
			turbulenceManifest, err := ioutil.ReadFile("fixtures/turbulence_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := turbulence.NewManifest(turbulence.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "turbulence",
				IPRange:      "10.244.4.0/24",
				BOSH: turbulence.ConfigBOSH{
					Target:   "some-bosh-target",
					Username: "some-bosh-username",
					Password: "some-bosh-password",
				},
			}, iaas.NewWardenConfig())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(MatchYAML(turbulenceManifest))
		})
	})
})