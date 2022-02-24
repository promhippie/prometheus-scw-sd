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
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	// providerPrefix defines the general prefix for all labels.
	providerPrefix = model.MetaLabelPrefix + "scaleway_"

	// Labels defines all available labels for this provider.
	Labels = map[string]string{
		"allowedActions":           providerPrefix + "allowed_actions",
		"arch":                     providerPrefix + "arch",
		"blade":                    providerPrefix + "blade",
		"bootscriptIdentifier":     providerPrefix + "bootscript_id",
		"bootscriptInitrd":         providerPrefix + "bootscript_initrd",
		"bootscriptKernel":         providerPrefix + "bootscript_kernel",
		"bootscriptTitle":          providerPrefix + "bootscript_title",
		"bootType":                 providerPrefix + "boot_type",
		"chassis":                  providerPrefix + "chassis",
		"cluster":                  providerPrefix + "cluster",
		"commercialType":           providerPrefix + "commercial_type",
		"description":              providerPrefix + "description",
		"domain":                   providerPrefix + "domain",
		"dynamicIPRequired":        providerPrefix + "dynamic_ip_required",
		"enableIPv6":               providerPrefix + "enable_ipv6",
		"hostname":                 providerPrefix + "hostname",
		"hypervisor":               providerPrefix + "hypervisor",
		"identifier":               providerPrefix + "id",
		"imageIdentifier":          providerPrefix + "image_id",
		"imageName":                providerPrefix + "image_name",
		"installHostname":          providerPrefix + "install_hostname",
		"installOs":                providerPrefix + "install_os",
		"installStatus":            providerPrefix + "install_status",
		"ips":                      providerPrefix + "ips",
		"ipv6":                     providerPrefix + "ipv6",
		"kind":                     providerPrefix + "kind",
		"name":                     providerPrefix + "name",
		"node":                     providerPrefix + "node",
		"offer":                    providerPrefix + "offer",
		"org":                      providerPrefix + "org",
		"placementGroupIdentifier": providerPrefix + "placement_group_id",
		"placementGroupName":       providerPrefix + "placement_group_name",
		"platform":                 providerPrefix + "platform",
		"privateHost":              providerPrefix + "private_host",
		"privateIP":                providerPrefix + "private_ipv4",
		"project":                  providerPrefix + "project",
		"protected":                providerPrefix + "protected",
		"publicHost":               providerPrefix + "public_host",
		"publicIP":                 providerPrefix + "public_ipv4",
		"securityGroupIdentifier":  providerPrefix + "security_group_id",
		"securityGroupName":        providerPrefix + "security_group_name",
		"state":                    providerPrefix + "state",
		"stateDetail":              providerPrefix + "state_detail",
		"status":                   providerPrefix + "status",
		"tags":                     providerPrefix + "tags",
		"zone":                     providerPrefix + "zone",
	}

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
							model.AddressLabel:                                  model.LabelValue(server.PublicIP.Address.String()),
							model.LabelName(Labels["project"]):                  model.LabelValue(project),
							model.LabelName(Labels["identifier"]):               model.LabelValue(server.ID),
							model.LabelName(Labels["name"]):                     model.LabelValue(server.Name),
							model.LabelName(Labels["org"]):                      model.LabelValue(server.Organization),
							model.LabelName(Labels["allowedActions"]):           model.LabelValue(actionsToString(d.separator, server.AllowedActions)),
							model.LabelName(Labels["tags"]):                     model.LabelValue(tagsToString(d.separator, server.Tags)),
							model.LabelName(Labels["commercialType"]):           model.LabelValue(server.CommercialType),
							model.LabelName(Labels["dynamicIPRequired"]):        model.LabelValue(boolToString(server.DynamicIPRequired)),
							model.LabelName(Labels["enableIPv6"]):               model.LabelValue(boolToString(server.EnableIPv6)),
							model.LabelName(Labels["hostname"]):                 model.LabelValue(server.Hostname),
							model.LabelName(Labels["imageIdentifier"]):          model.LabelValue(imageIdentifier),
							model.LabelName(Labels["imageName"]):                model.LabelValue(imageName),
							model.LabelName(Labels["protected"]):                model.LabelValue(boolToString(server.Protected)),
							model.LabelName(Labels["privateIP"]):                model.LabelValue(*server.PrivateIP),
							model.LabelName(Labels["privateHost"]):              model.LabelValue(fmt.Sprintf("%s.priv.cloud.scaleway.com", server.ID)),
							model.LabelName(Labels["publicIP"]):                 model.LabelValue(server.PublicIP.Address.String()),
							model.LabelName(Labels["publicHost"]):               model.LabelValue(fmt.Sprintf("%s.pub.cloud.scaleway.com", server.ID)),
							model.LabelName(Labels["state"]):                    model.LabelValue(server.State),
							model.LabelName(Labels["cluster"]):                  model.LabelValue(locationCluster),
							model.LabelName(Labels["hypervisor"]):               model.LabelValue(locationHypervisor),
							model.LabelName(Labels["node"]):                     model.LabelValue(locationNode),
							model.LabelName(Labels["platform"]):                 model.LabelValue(locationPlatform),
							model.LabelName(Labels["ipv6"]):                     model.LabelValue(ipv6Address),
							model.LabelName(Labels["bootscriptIdentifier"]):     model.LabelValue(bootscriptIdentifier),
							model.LabelName(Labels["bootscriptTitle"]):          model.LabelValue(bootscriptTitle),
							model.LabelName(Labels["bootscriptKernel"]):         model.LabelValue(bootscriptKernel),
							model.LabelName(Labels["bootscriptInitrd"]):         model.LabelValue(bootscriptInitrd),
							model.LabelName(Labels["bootType"]):                 model.LabelValue(server.BootType),
							model.LabelName(Labels["securityGroupIdentifier"]):  model.LabelValue(securityGroupIdentifier),
							model.LabelName(Labels["securityGroupName"]):        model.LabelValue(securityGroupName),
							model.LabelName(Labels["stateDetail"]):              model.LabelValue(server.StateDetail),
							model.LabelName(Labels["arch"]):                     model.LabelValue(server.Arch),
							model.LabelName(Labels["placementGroupIdentifier"]): model.LabelValue(placementGroupIdentifier),
							model.LabelName(Labels["placementGroupName"]):       model.LabelValue(placementGroupName),
							model.LabelName(Labels["zone"]):                     model.LabelValue(server.Zone),
							model.LabelName(Labels["kind"]):                     model.LabelValue("instance"),
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
							model.AddressLabel:                         model.LabelValue(server.IPs[0].Address.String()),
							model.LabelName(Labels["project"]):         model.LabelValue(project),
							model.LabelName(Labels["identifier"]):      model.LabelValue(server.ID),
							model.LabelName(Labels["org"]):             model.LabelValue(server.OrganizationID),
							model.LabelName(Labels["name"]):            model.LabelValue(server.Name),
							model.LabelName(Labels["description"]):     model.LabelValue(server.Description),
							model.LabelName(Labels["status"]):          model.LabelValue(server.Status),
							model.LabelName(Labels["offer"]):           model.LabelValue(server.OfferID),
							model.LabelName(Labels["installOs"]):       model.LabelValue(installOs),
							model.LabelName(Labels["installHostname"]): model.LabelValue(installHostname),
							model.LabelName(Labels["installStatus"]):   model.LabelValue(installStatus),
							model.LabelName(Labels["tags"]):            model.LabelValue(tagsToString(d.separator, server.Tags)),
							model.LabelName(Labels["ips"]):             model.LabelValue(ipsToString(d.separator, server.IPs)),
							model.LabelName(Labels["domain"]):          model.LabelValue(server.Domain),
							model.LabelName(Labels["bootType"]):        model.LabelValue(server.BootType),
							model.LabelName(Labels["zone"]):            model.LabelValue(server.Zone),
							model.LabelName(Labels["kind"]):            model.LabelValue("baremetal"),
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
