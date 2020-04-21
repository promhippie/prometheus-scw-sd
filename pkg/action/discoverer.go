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
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	scwPrefix                     = model.MetaLabelPrefix + "scaleway_"
	allowedActionsLabel           = scwPrefix + "allowed_actions"
	archLabel                     = scwPrefix + "arch"
	bladeLabel                    = scwPrefix + "blade"
	bootTypeLabel                 = scwPrefix + "boot_type"
	bootscriptIdentifierLabel     = scwPrefix + "bootscript_id"
	bootscriptInitrdLabel         = scwPrefix + "bootscript_initrd"
	bootscriptKernelLabel         = scwPrefix + "bootscript_kernel"
	bootscriptTitleLabel          = scwPrefix + "bootscript_title"
	chassisLabel                  = scwPrefix + "chassis"
	clusterLabel                  = scwPrefix + "cluster"
	commercialTypeLabel           = scwPrefix + "commercial_type"
	descriptionLabel              = scwPrefix + "description"
	domainLabel                   = scwPrefix + "domain"
	dynamicIPRequiredLabel        = scwPrefix + "dynamic_ip_required"
	enableIPv6Label               = scwPrefix + "enable_ipv6"
	hostnameLabel                 = scwPrefix + "hostname"
	hypervisorLabel               = scwPrefix + "hypervisor"
	identifierLabel               = scwPrefix + "id"
	imageIdentifierLabel          = scwPrefix + "image_id"
	imageNameLabel                = scwPrefix + "image_name"
	installHostnameLabel          = scwPrefix + "install_hostname"
	installOsLabel                = scwPrefix + "install_os"
	installStatusLabel            = scwPrefix + "install_status"
	ipsLabel                      = scwPrefix + "ips"
	ipv6Label                     = scwPrefix + "ipv6"
	kindLabel                     = scwPrefix + "kind"
	nameLabel                     = scwPrefix + "name"
	nodeLabel                     = scwPrefix + "node"
	offerLabel                    = scwPrefix + "offer"
	orgLabel                      = scwPrefix + "org"
	placementGroupIdentifierLabel = scwPrefix + "placement_group_id"
	placementGroupNameLabel       = scwPrefix + "placement_group_name"
	platformLabel                 = scwPrefix + "platform"
	privateHostLabel              = scwPrefix + "private_host"
	privateIPLabel                = scwPrefix + "private_ipv4"
	projectLabel                  = scwPrefix + "project"
	protectedLabel                = scwPrefix + "protected"
	publicHostLabel               = scwPrefix + "public_host"
	publicIPLabel                 = scwPrefix + "public_ipv4"
	securityGroupIdentifierLabel  = scwPrefix + "security_group_id"
	securityGroupNameLabel        = scwPrefix + "security_group_name"
	stateDetailLabel              = scwPrefix + "state_detail"
	stateLabel                    = scwPrefix + "state"
	statusLabel                   = scwPrefix + "status"
	tagsLabel                     = scwPrefix + "tags"
	zoneLabel                     = scwPrefix + "zone"
)

var (
	// ErrClientFailed defines an error if the client init fails.
	ErrClientFailed = errors.New("failed to initialize client")

	// ErrClientForbidden defines an error if the authentication fails.
	ErrClientForbidden = errors.New("failed to authenticate client")

	// ErrInvalidZone defines an error if an invalid zone have been provided.
	ErrInvalidZone = errors.New("invalid zone provided")
)

// Discoverer implements the Prometheus discoverer interface.
type Discoverer struct {
	clients        map[string]*scw.Client
	logger         log.Logger
	refresh        int
	checkInstance  bool
	instanceZones  []string
	checkBaremetal bool
	baremetalZones []string
	separator      string
	lasts          map[string]struct{}
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
		if d.checkInstance {
			zones := make([]scw.Zone, 0)

			if zone, _ := client.GetDefaultZone(); zone != "" {
				zones = []scw.Zone{
					zone,
				}
			} else {
				for _, raw := range d.instanceZones {
					zones = append(zones, scw.Zone(raw))
				}
			}

			for _, zone := range zones {
				now := time.Now()

				resp, err := instance.NewAPI(client).ListServers(
					&instance.ListServersRequest{
						Zone: zone,
					},
					scw.WithAllPages(),
					scw.WithContext(ctx),
				)

				requestDuration.WithLabelValues(project, "instance", zone.String()).Observe(time.Since(now).Seconds())

				if err != nil {
					level.Warn(d.logger).Log(
						"msg", "Failed to fetch servers",
						"project", project,
						"kind", "instance",
						"zone", zone,
						"err", err,
					)

					requestFailures.WithLabelValues(project, "instance", zone.String()).Inc()
					continue
				}

				level.Debug(d.logger).Log(
					"msg", "Requested servers",
					"project", project,
					"kind", "instance",
					"zone", zone,
					"count", resp.TotalCount,
				)

				for _, server := range resp.Servers {
					var (
						imageIdentifier          string
						imageName                string
						locationCluster          string
						locationHypervisor       string
						locationNode             string
						locationPlatform         string
						bootscriptIdentifier     string
						bootscriptTitle          string
						bootscriptKernel         string
						bootscriptInitrd         string
						securityGroupIdentifier  string
						securityGroupName        string
						placementGroupIdentifier string
						placementGroupName       string
						ipv6Address              string
					)

					if server.Image != nil {
						imageIdentifier = server.Image.ID
						imageName = server.Image.Name
					}

					if server.Location != nil {
						locationCluster = server.Location.ClusterID
						locationHypervisor = server.Location.HypervisorID
						locationNode = server.Location.NodeID
						locationPlatform = server.Location.PlatformID
					}

					if server.Bootscript != nil {
						bootscriptIdentifier = server.Bootscript.ID
						bootscriptTitle = server.Bootscript.Title
						bootscriptKernel = server.Bootscript.Kernel
						bootscriptInitrd = server.Bootscript.Initrd
					}

					if server.SecurityGroup != nil {
						securityGroupIdentifier = server.SecurityGroup.ID
						securityGroupName = server.SecurityGroup.Name
					}

					if server.PlacementGroup != nil {
						placementGroupIdentifier = server.PlacementGroup.ID
						placementGroupName = server.PlacementGroup.Name
					}

					addresses := make([]model.LabelSet, 0)

					addresses = append(addresses, model.LabelSet{
						model.AddressLabel: model.LabelValue(server.PublicIP.Address.String()),
					})

					if server.EnableIPv6 && server.IPv6 != nil {
						ipv6Address = server.IPv6.Address.String()

						addresses = append(addresses, model.LabelSet{
							model.AddressLabel: model.LabelValue(ipv6Address),
						})
					}

					target := &targetgroup.Group{
						Source:  fmt.Sprintf("instance/%s", server.ID),
						Targets: addresses,
						Labels: model.LabelSet{
							model.AddressLabel:                             model.LabelValue(server.PublicIP.Address.String()),
							model.LabelName(projectLabel):                  model.LabelValue(project),
							model.LabelName(identifierLabel):               model.LabelValue(server.ID),
							model.LabelName(nameLabel):                     model.LabelValue(server.Name),
							model.LabelName(orgLabel):                      model.LabelValue(server.Organization),
							model.LabelName(allowedActionsLabel):           model.LabelValue(actionsToString(d.separator, server.AllowedActions)),
							model.LabelName(tagsLabel):                     model.LabelValue(tagsToString(d.separator, server.Tags)),
							model.LabelName(commercialTypeLabel):           model.LabelValue(server.CommercialType),
							model.LabelName(dynamicIPRequiredLabel):        model.LabelValue(boolToString(server.DynamicIPRequired)),
							model.LabelName(enableIPv6Label):               model.LabelValue(boolToString(server.EnableIPv6)),
							model.LabelName(hostnameLabel):                 model.LabelValue(server.Hostname),
							model.LabelName(imageIdentifierLabel):          model.LabelValue(imageIdentifier),
							model.LabelName(imageNameLabel):                model.LabelValue(imageName),
							model.LabelName(protectedLabel):                model.LabelValue(boolToString(server.Protected)),
							model.LabelName(privateIPLabel):                model.LabelValue(*server.PrivateIP),
							model.LabelName(privateHostLabel):              model.LabelValue(fmt.Sprintf("%s.priv.cloud.scaleway.com", server.ID)),
							model.LabelName(publicIPLabel):                 model.LabelValue(server.PublicIP.Address.String()),
							model.LabelName(publicHostLabel):               model.LabelValue(fmt.Sprintf("%s.pub.cloud.scaleway.com", server.ID)),
							model.LabelName(stateLabel):                    model.LabelValue(server.State),
							model.LabelName(clusterLabel):                  model.LabelValue(locationCluster),
							model.LabelName(hypervisorLabel):               model.LabelValue(locationHypervisor),
							model.LabelName(nodeLabel):                     model.LabelValue(locationNode),
							model.LabelName(platformLabel):                 model.LabelValue(locationPlatform),
							model.LabelName(ipv6Label):                     model.LabelValue(ipv6Address),
							model.LabelName(bootscriptIdentifierLabel):     model.LabelValue(bootscriptIdentifier),
							model.LabelName(bootscriptTitleLabel):          model.LabelValue(bootscriptTitle),
							model.LabelName(bootscriptKernelLabel):         model.LabelValue(bootscriptKernel),
							model.LabelName(bootscriptInitrdLabel):         model.LabelValue(bootscriptInitrd),
							model.LabelName(bootTypeLabel):                 model.LabelValue(server.BootType),
							model.LabelName(securityGroupIdentifierLabel):  model.LabelValue(securityGroupIdentifier),
							model.LabelName(securityGroupNameLabel):        model.LabelValue(securityGroupName),
							model.LabelName(stateDetailLabel):              model.LabelValue(server.StateDetail),
							model.LabelName(archLabel):                     model.LabelValue(server.Arch),
							model.LabelName(placementGroupIdentifierLabel): model.LabelValue(placementGroupIdentifier),
							model.LabelName(placementGroupNameLabel):       model.LabelValue(placementGroupName),
							model.LabelName(zoneLabel):                     model.LabelValue(server.Zone),
							model.LabelName(kindLabel):                     model.LabelValue("instance"),
						},
					}

					level.Debug(d.logger).Log(
						"msg", "Server added",
						"project", project,
						"kind", "instance",
						"zone", zone,
						"source", target.Source,
					)

					current[target.Source] = struct{}{}
					targets = append(targets, target)

				}
			}
		}

		if d.checkBaremetal {
			zones := make([]scw.Zone, 0)

			if zone, _ := client.GetDefaultZone(); zone != "" {
				zones = []scw.Zone{
					zone,
				}
			} else {
				for _, raw := range d.baremetalZones {
					zones = append(zones, scw.Zone(raw))
				}
			}

			for _, zone := range zones {
				now := time.Now()

				resp, err := baremetal.NewAPI(client).ListServers(
					&baremetal.ListServersRequest{
						Zone: zone,
					},
					scw.WithAllPages(),
					scw.WithContext(ctx),
				)

				requestDuration.WithLabelValues(project, "baremetal", zone.String()).Observe(time.Since(now).Seconds())

				if err != nil {
					level.Warn(d.logger).Log(
						"msg", "Failed to fetch servers",
						"project", project,
						"kind", "baremetal",
						"zone", zone,
						"err", err,
					)

					requestFailures.WithLabelValues(project, "baremetal", zone.String()).Inc()
					continue
				}

				level.Debug(d.logger).Log(
					"msg", "Requested servers",
					"project", project,
					"kind", "baremetal",
					"zone", zone,
					"count", resp.TotalCount,
				)

				for _, server := range resp.Servers {
					if len(server.IPs) < 1 {
						continue
					}

					var (
						installOs       string
						installHostname string
						installStatus   string
					)

					if server.Install != nil {
						installOs = server.Install.OsID
						installHostname = server.Install.Hostname
						installStatus = server.Install.Status.String()
					}

					addresses := make([]model.LabelSet, len(server.IPs))

					for _, ip := range server.IPs {
						addresses = append(addresses, model.LabelSet{
							model.AddressLabel: model.LabelValue(ip.Address.String()),
						})
					}

					target := &targetgroup.Group{
						Source:  fmt.Sprintf("baremetal/%s", server.ID),
						Targets: addresses,
						Labels: model.LabelSet{
							model.AddressLabel:                    model.LabelValue(server.IPs[0].Address.String()),
							model.LabelName(projectLabel):         model.LabelValue(project),
							model.LabelName(identifierLabel):      model.LabelValue(server.ID),
							model.LabelName(orgLabel):             model.LabelValue(server.OrganizationID),
							model.LabelName(nameLabel):            model.LabelValue(server.Name),
							model.LabelName(descriptionLabel):     model.LabelValue(server.Description),
							model.LabelName(statusLabel):          model.LabelValue(server.Status),
							model.LabelName(offerLabel):           model.LabelValue(server.OfferID),
							model.LabelName(installOsLabel):       model.LabelValue(installOs),
							model.LabelName(installHostnameLabel): model.LabelValue(installHostname),
							model.LabelName(installStatusLabel):   model.LabelValue(installStatus),
							model.LabelName(tagsLabel):            model.LabelValue(tagsToString(d.separator, server.Tags)),
							model.LabelName(ipsLabel):             model.LabelValue(ipsToString(d.separator, server.IPs)),
							model.LabelName(domainLabel):          model.LabelValue(server.Domain),
							model.LabelName(bootTypeLabel):        model.LabelValue(server.BootType),
							model.LabelName(zoneLabel):            model.LabelValue(server.Zone),
							model.LabelName(kindLabel):            model.LabelValue("baremetal"),
						},
					}

					level.Debug(d.logger).Log(
						"msg", "Server added",
						"project", project,
						"kind", "baremetal",
						"zone", zone,
						"source", target.Source,
					)

					current[target.Source] = struct{}{}
					targets = append(targets, target)
				}
			}
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

// actionsToString transforms a slice if actions to a string separated by defined separator.
func actionsToString(separator string, vals []instance.ServerAction) string {
	res := []string{}

	for _, val := range vals {
		res = append(res, val.String())
	}

	return strings.Join(res, separator)
}

// tagsToString transforms a slice if tags to a string separated by defined separator.
func tagsToString(separator string, vals []string) string {
	return strings.Join(vals, separator)
}

// ipsToString transforms a slice if ips to a string separated by defined separator.
func ipsToString(separator string, vals []*baremetal.IP) string {
	res := []string{}

	for _, val := range vals {
		res = append(res, val.Address.String())
	}

	return strings.Join(res, separator)
}

func boolToString(val bool) string {
	if val {
		return "true"
	}

	return "false"
}
