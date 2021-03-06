package prommer

import (
	"errors"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// ServiceProvider provides the services to be monitored
type ServiceProvider struct {
	docker          *client.Client
	monitoringLabel string
}

// CreateServiceProvider creates a new instance of service provider
// Returns error if the client or monitoringLabel is nil
func CreateServiceProvider(cli *client.Client, monitoringLabel string) (*ServiceProvider, error) {
	if cli == nil {
		return nil, errors.New("client can not be nil")
	}
	if &monitoringLabel == nil {
		return nil, errors.New("monitoringLabel can not be nil")
	}

	return &ServiceProvider{
		docker:          cli,
		monitoringLabel: monitoringLabel,
	}, nil
}

// GetServices returns the
func (sp *ServiceProvider) GetServices() ([]*Service, error) {

	filters := filters.NewArgs()
	filters.Add("label", sp.monitoringLabel)

	options := types.ContainerListOptions{
		All:    false,
		Filter: filters,
	}

	list, err := sp.docker.ContainerList(context.Background(), options)

	if err != nil {
		return nil, err
	}

	serviceMap := make(map[string]*Service)

	for _, c := range list {
		var (
			serviceName = c.Labels[sp.monitoringLabel]
			service     *Service
		)

		//fmt.Println("IP: " + port.IP + ", private port: " + strconv.FormatInt(port.PrivatePort, 10) + ", public port: " + strconv.FormatInt(port.PublicPort, 10) + ", type: " + port.Type)

		if service = serviceMap[serviceName]; service == nil {
			service = &Service{
				Name:      serviceName,
				Instances: make([]*Instance, 0),
			}
			serviceMap[serviceName] = service
		}

		s := sp.convert(c, FindPortFromContainer)

		if s != nil {
			service.Instances = append(service.Instances, s)
		}

	}

	services := make([]*Service, 0, len(serviceMap))

	for _, s := range serviceMap {
		services = append(services, s)
	}

	return services, nil
}

func (sp *ServiceProvider) convert(c types.Container, findPort FindPort) *Instance {

	portLabel := sp.monitoringLabel + ".port"

	findPortOptions := &FindPortOptions{
		Label: &portLabel,
	}

	port := findPort(c, findPortOptions)

	if port == 0 {
		return nil
	}

	return &Instance{
		HostIP:   c.NetworkSettings.Networks[c.HostConfig.NetworkMode].IPAddress,
		HostPort: strconv.Itoa(port),
	}
}
