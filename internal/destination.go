package internal

import "strings"

type Destination struct {
	Path       string
	RemoteHost string
	RemotePath string
}

func ParseDestination(path string) *Destination {
	remoteHost := ""
	remotePath := ""
	components := strings.SplitN(path, ":", 2)
	if len(components) == 2 {
		remoteHost = components[0]
		remotePath = components[1]
	}
	return &Destination{
		Path:       path,
		RemotePath: remotePath,
		RemoteHost: remoteHost,
	}
}
