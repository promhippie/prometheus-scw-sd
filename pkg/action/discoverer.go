package action

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/scaleway/go-scaleway"
)

const (
	scwPrefix            = model.MetaLabelPrefix + "scaleway_"
	projectLabel         = scwPrefix + "project"
	nameLabel            = scwPrefix + "name"
	identifierLabel      = scwPrefix + "id"
	archLabel            = scwPrefix + "arch"
	imageIdentifierLabel = scwPrefix + "image_id"
	imageNameLabel       = scwPrefix + "image_name"
	publicIPLabel        = scwPrefix + "public_ipv4"
	publicHostLabel      = scwPrefix + "public_host"
	stateLabel           = scwPrefix + "state"
	privateIPLabel       = scwPrefix + "private_ipv4"
	privateHostLabel     = scwPrefix + "private_host"
	hostnameLabel        = scwPrefix + "hostname"
	orgLabel             = scwPrefix + "org"
	commercialTypeLabel  = scwPrefix + "commercial_type"
	platformLabel        = scwPrefix + "platform"
	hypervisorLabel      = scwPrefix + "hypervisor"
	nodeLabel            = scwPrefix + "node"
	bladeLabel           = scwPrefix + "blade"
	chassisLabel         = scwPrefix + "chassis"
	clusterLabel         = scwPrefix + "cluster"
	zoneLabel            = scwPrefix + "zone"
	tagsLabel            = scwPrefix + "tags"
)

var (
	// ErrClientFailed defines an error if the client init fails.
	ErrClientFailed = errors.New("failed to initialize client")

	// ErrClientForbidden defines an error if the authentication fails.
	ErrClientForbidden = errors.New("failed to authenticate client")
)

// Discoverer implements the Prometheus discoverer interface.
type Discoverer struct {
	clients   map[string]*api.ScalewayAPI
	logger    log.Logger
	refresh   int
	separator string
	lasts     map[string]struct{}
}

// Run initializes fetching the targets for service discovery.
func (d Discoverer) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	ticker := time.NewTicker(time.Duration(d.refresh) * time.Second)

	for {
		targets, err := d.getTargets(ctx)

		if err == nil {
			ch <- targets
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (d *Discoverer) getTargets(ctx context.Context) ([]*targetgroup.Group, error) {
	current := make(map[string]struct{})
	targets := make([]*targetgroup.Group, 0)

	for project, client := range d.clients {

		now := time.Now()
		servers, err := client.GetServers(false, 0)
		requestDuration.WithLabelValues(project).Observe(time.Since(now).Seconds())

		if err != nil {
			level.Warn(d.logger).Log(
				"msg", "Failed to fetch servers",
				"project", project,
				"err", err,
			)

			requestFailures.WithLabelValues(project).Inc()
			return nil, err
		}

		level.Debug(d.logger).Log(
			"msg", "Requested servers",
			"project", project,
			"count", len(*servers),
		)

		for _, server := range *servers {
			target := &targetgroup.Group{
				Source: fmt.Sprintf("scaleway/%s", server.Identifier),
				Targets: []model.LabelSet{
					{
						model.AddressLabel: model.LabelValue(server.PublicAddress.IP),
					},
				},
				Labels: model.LabelSet{
					model.AddressLabel:                    model.LabelValue(server.PublicAddress.IP),
					model.LabelName(projectLabel):         model.LabelValue(project),
					model.LabelName(nameLabel):            model.LabelValue(server.Name),
					model.LabelName(identifierLabel):      model.LabelValue(server.Identifier),
					model.LabelName(archLabel):            model.LabelValue(server.Arch),
					model.LabelName(imageIdentifierLabel): model.LabelValue(server.Image.Identifier),
					model.LabelName(imageNameLabel):       model.LabelValue(server.Image.Name),
					model.LabelName(publicIPLabel):        model.LabelValue(server.PublicAddress.IP),
					model.LabelName(publicHostLabel):      model.LabelValue(fmt.Sprintf("%s.pub.cloud.scaleway.com", server.Identifier)),
					model.LabelName(stateLabel):           model.LabelValue(server.State),
					model.LabelName(privateIPLabel):       model.LabelValue(server.PrivateIP),
					model.LabelName(privateHostLabel):     model.LabelValue(fmt.Sprintf("%s.priv.cloud.scaleway.com", server.Identifier)),
					model.LabelName(hostnameLabel):        model.LabelValue(server.Hostname),
					model.LabelName(orgLabel):             model.LabelValue(server.Organization),
					model.LabelName(commercialTypeLabel):  model.LabelValue(server.CommercialType),
					model.LabelName(platformLabel):        model.LabelValue(server.Location.Platform),
					model.LabelName(hypervisorLabel):      model.LabelValue(server.Location.Hypervisor),
					model.LabelName(nodeLabel):            model.LabelValue(server.Location.Node),
					model.LabelName(bladeLabel):           model.LabelValue(server.Location.Blade),
					model.LabelName(chassisLabel):         model.LabelValue(server.Location.Chassis),
					model.LabelName(clusterLabel):         model.LabelValue(server.Location.Cluster),
					model.LabelName(zoneLabel):            model.LabelValue(server.Location.ZoneID),
					model.LabelName(tagsLabel):            model.LabelValue(strings.Join(server.Tags, d.separator)),
				},
			}

			level.Debug(d.logger).Log(
				"msg", "Server added",
				"project", project,
				"source", target.Source,
			)

			current[target.Source] = struct{}{}
			targets = append(targets, target)
		}

	}

	for k := range d.lasts {
		if _, ok := current[k]; !ok {
			level.Debug(d.logger).Log(
				"msg", "Server deleted",
				"source", k,
			)

			targets = append(
				targets,
				&targetgroup.Group{
					Source: k,
				},
			)
		}
	}

	d.lasts = current
	return targets, nil
}
