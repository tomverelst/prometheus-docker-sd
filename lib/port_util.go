package prommer

import (
	"strconv"

	dockertypes "github.com/docker/engine-api/types"
)

// FindPortOptions holds the optional parameters for FindPort
type FindPortOptions struct {
	Label *string
}

// FindPort attempts to find the correct metrics port
func FindPort(c dockertypes.Container, options *FindPortOptions) int {

	portFromLabels := findPortFromLabels(c, options)

	if portFromLabels != 0 {
		return portFromLabels
	}

	publicPorts := onlyPublicPorts(c.Ports)
	amountOfPorts := len(publicPorts)

	var port *dockertypes.Port

	if amountOfPorts == 0 {
		return 0
	}

	if amountOfPorts == 1 {
		return publicPorts[0].PublicPort
	}

	// If there are multiple ports
	// prefer the one that forwards to port 80
	port = findPortEighty(publicPorts)

	if port == nil {
		port = &publicPorts[0]
	}

	if port == nil {
		return 0
	}

	return port.PublicPort
}

func findPortFromLabels(c dockertypes.Container, options *FindPortOptions) int {
	if options != nil && options.Label != nil {
		portString := c.Labels[*options.Label]
		if &portString != nil {
			port, err := strconv.Atoi(portString)
			if err != nil {
				return 0
			}
			return port
		}
	}
	return 0
}

func onlyPublicPorts(ports []dockertypes.Port) []dockertypes.Port {
	var publicPorts []dockertypes.Port
	for _, port := range ports {
		if port.PublicPort != 0 {
			publicPorts = append(publicPorts, port)
		}
	}
	return publicPorts
}

func findPortEighty(ports []dockertypes.Port) *dockertypes.Port {
	for _, port := range ports {
		if port.PrivatePort == 80 {
			return &port
		}
	}
	return nil
}