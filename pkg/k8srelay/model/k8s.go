/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package model

import "gopkg.in/yaml.v3"

type Cluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}
type Clusters struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}
type Context struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}
type Contexts struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}
type User struct {
	Token                 string `yaml:"token"`
	ClientCertificateData string `yaml:"client-certificate-data,omitempty"`
	ClientKeyData         string `yaml:"client-key-data,omitempty"`
}
type Users struct {
	User User   `json:"user"`
	Name string `json:"name"`
}
type KubeConfigCluster struct {
	ApiVersion     string     `yaml:"apiVersion"`
	Clusters       []Clusters `yaml:"clusters"`
	Contexts       []Contexts `yaml:"contexts"`
	CurrentContext string     `yaml:"current-context"`
	Kind           string     `yaml:"kind"`
	Users          []Users    `yaml:"users"`
}

func NewRelayKubeConfig(routeID string, relayURL string) ([]byte, error) {
	x := new(KubeConfigCluster)
	x.ApiVersion = "v1"
	x.CurrentContext = "context0"
	x.Kind = "Config"

	x.Clusters = make([]Clusters, 1)
	x.Clusters[0].Cluster.Server = relayURL
	x.Clusters[0].Cluster.CertificateAuthorityData = loadCA()
	x.Clusters[0].Name = "cluster0"
	x.Users = make([]Users, 1)
	x.Users[0].Name = "user0"
	x.Users[0].User.Token = routeID
	x.Contexts = make([]Contexts, 1)
	x.Contexts[0].Name = "context0"
	x.Contexts[0].Context.User = "user0"
	x.Contexts[0].Context.Cluster = "cluster0"
	ret, err := yaml.Marshal(x)
	return ret, err
}

// loadCA gets a base 64 of the cert authority certificate
func loadCA() string {
	return "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURoRENDQW13Q0NRQ01pd1VDK0hkTTB6QU5CZ2txaGtpRzl3MEJBUXNGQURDQmd6RUxNQWtHQTFVRUJoTUMKUlVNeER6QU5CZ05WQkFnTUJrMWhibUZpYVRFT01Bd0dBMVVFQnd3RlRXRnVkR0V4RkRBU0JnTlZCQW9NQzBWdQpaMmx1WldWeWFXNW5NUXd3Q2dZRFZRUUxEQU5FWlhZeERUQUxCZ05WQkFNTUJISnZiM1F4SURBZUJna3Foa2lHCjl3MEJDUUVXRVcxaGMyOXVZa0J1WlhSaGNIQXVZMjl0TUI0WERUSXpNRE15TWpFMk5UWXhNMW9YRFRJNE1ETXkKTURFMk5UWXhNMW93Z1lNeEN6QUpCZ05WQkFZVEFrVkRNUTh3RFFZRFZRUUlEQVpOWVc1aFlta3hEakFNQmdOVgpCQWNNQlUxaGJuUmhNUlF3RWdZRFZRUUtEQXRGYm1kcGJtVmxjbWx1WnpFTU1Bb0dBMVVFQ3d3RFJHVjJNUTB3CkN3WURWUVFEREFSeWIyOTBNU0F3SGdZSktvWklodmNOQVFrQkZoRnRZWE52Ym1KQWJtVjBZWEJ3TG1OdmJUQ0MKQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFMUkh4LzMybnZORHF6RWJTeDZXWVZjdgptT3pEcWhDVUpmd0dPbHVhU1UzTytlbkMrNFVoWkdndVR3NXl3MEQ2YXltZzJNUWdHTlRjT254V1NPeE5mS0V3ClhGa0t5UnYrVXZPOE9MeFFXYlFOUzBGNW1mZzJDNnBPc21uMG8vVmprTHloWCtYbkM0TzVjd3hkd0c0SjVvYlUKTGJNSldZb0FJVWxuK0ZkL21VSGU3bWx0RnNsRVNmNTRoS2VTYkdRc212NzBBalBoV0h6Ulh1dG9JVGR6a1NwUQp3WXB4Vm8zeWY4aUc5cDE4WHc5N1k1T1pNYzZQQTRIekFaUUduVHpmZDhRUUlVTWxzMVhwNlJJMGQ1ekpnRzNDCnpqK3BUMGdtbFlCSE00cVlsQTJjR2QyYWQ0UjlRRzZZdjE1ZXBlYU1FZkVmL2REUkFpdG1LUDFzdWtKMXhLRUMKQXdFQUFUQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFYd2g1cVhLU2dqOXlNeU9NUTZwREdSSlpJb2VhNDVEdgpwdFJFR0RiOWEwS042WHhYL1BscXpPQ25lNExkem4vMzRYWHJMV2hBZTEvM3FpaGlpUER3MmM1Vjl5SWIvMW9aClBGSlVobU1tNm5lWEQ1aWRSVDdXQUxlZzljU0h1U0IyeHZxTi9hYkd4RGROSTRhZXJxNGNzc1FvaTQ2WjNVZVcKMU82YWpoSUUyYUs3eFNPT2xWdGhDV2hTa1drb20yUHIrNE80RzFkaWdyV2Z2V1hqQ2k1SnFZcUVleWExSTVMaAptU3lCeW5UZy8vVy9OQ01LTEErOUk4eitmQis0Q3ByeUJzYkVWMngyM2RWemVseHYxZkI5eCtTT2ZQOFN4MkI0CmEzZjdyaVQwUldjMFgvY293UnVpYlJSbzdhRDRMYzFqaERJaTBUSm9zZVBWY3AzSlo5NUVWUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
}
