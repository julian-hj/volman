package voldiscoverers

import (
	"path/filepath"

	. "code.cloudfoundry.org/csiplugin"
	"code.cloudfoundry.org/csiplugin/oshelper"
	"code.cloudfoundry.org/csishim"
	"code.cloudfoundry.org/goshims/filepathshim"
	"code.cloudfoundry.org/goshims/grpcshim"
	"code.cloudfoundry.org/goshims/osshim"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/volman"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"errors"
)

type csiPluginDiscoverer struct {
	logger          lager.Logger
	pluginRegistry  volman.PluginRegistry
	pluginPaths     []string
	filepathShim    filepathshim.Filepath
	grpcShim        grpcshim.Grpc
	csiShim         csishim.Csi
	osShim          osshim.Os
	csiMountRootDir string
}

func NewCsiPluginDiscoverer(logger lager.Logger, pluginRegistry volman.PluginRegistry, pluginPaths []string, csiMountRootDir string) volman.Discoverer {
	return &csiPluginDiscoverer{
		logger:          logger,
		pluginRegistry:  pluginRegistry,
		pluginPaths:     pluginPaths,
		filepathShim:    &filepathshim.FilepathShim{},
		grpcShim:        &grpcshim.GrpcShim{},
		csiShim:         &csishim.CsiShim{},
		osShim:          &osshim.OsShim{},
		csiMountRootDir: csiMountRootDir,
	}
}

func NewCsiPluginDiscovererWithShims(logger lager.Logger, pluginRegistry volman.PluginRegistry, pluginPaths []string, filepathShim filepathshim.Filepath, grpcShim grpcshim.Grpc, csiShim csishim.Csi, osShim osshim.Os, csiMountRootDir string) volman.Discoverer {
	return &csiPluginDiscoverer{
		logger:          logger,
		pluginRegistry:  pluginRegistry,
		pluginPaths:     pluginPaths,
		filepathShim:    filepathShim,
		grpcShim:        grpcShim,
		csiShim:         csiShim,
		osShim:          osShim,
		csiMountRootDir: csiMountRootDir,
	}
}

func (p *csiPluginDiscoverer) Discover(logger lager.Logger) (map[string]volman.Plugin, error) {
	logger = logger.Session("discover")
	logger.Debug("start")
	defer logger.Debug("end")
	var conns []grpcshim.ClientConn
	defer func() {
		for _, conn := range conns {
			err := conn.Close()
			if err != nil {
				logger.Error("grpc-conn-close", err)
			}
		}
	}()

	plugins := map[string]volman.Plugin{}

	for _, pluginPath := range p.pluginPaths {
		pluginSpecFiles, err := filepath.Glob(pluginPath + "/*.json")
		if err != nil {
			logger.Error("filepath-glob", err, lager.Data{"glob": pluginPath + "/*.json"})
			return plugins, err
		}
		for _, pluginSpecFile := range pluginSpecFiles {
			csiPluginSpec, err := ReadSpec(logger, pluginSpecFile)
			if err != nil {
				logger.Error("read-spec-failed", err, lager.Data{"plugin-path": pluginPath, "plugin-spec-file": pluginSpecFile})
				continue
			}
			// instantiate a volman.Plugin implementation of a csi.NodePlugin
			logger.Debug("rpc-dial", lager.Data{"address": csiPluginSpec.Address})
			conn, err := p.grpcShim.Dial(csiPluginSpec.Address, grpc.WithInsecure())
			conns = append(conns, conn)
			if err != nil {
				logger.Error("grpc-dial", err, lager.Data{"address": csiPluginSpec.Address})
				continue
			}

			identityPlugin := p.csiShim.NewIdentityClient(conn)

			pluginInfo, err := identityPlugin.GetPluginInfo(context.TODO(), &csi.GetPluginInfoRequest{})
			if err != nil {
				logger.Error("plugin-info-error", err)
				continue
			}

			pluginCapabilities, err := identityPlugin.GetPluginCapabilities(context.TODO(), &csi.GetPluginCapabilitiesRequest{})
			if err != nil {
				logger.Error("plugin-capabilities-error", err)
				continue
			}

			pluginHasAccessibilityConstraints := false

			for _, capability := range pluginCapabilities.GetCapabilities() {
				service := capability.GetService()

				if service.GetType() == csi.PluginCapability_Service_ACCESSIBILITY_CONSTRAINTS {
					pluginHasAccessibilityConstraints = true
				}
			}

			if pluginHasAccessibilityConstraints {
				logger.Error("plugin-capability-check", errors.New("accessibility constraints unsupported"))
				continue
			}

			csiPluginName := pluginInfo.Name
			existingPlugin, found := p.pluginRegistry.Plugins()[csiPluginName]
			pluginSpec := volman.PluginSpec{
				Name:    csiPluginName,
				Address: csiPluginSpec.Address,
			}

			if !found || !existingPlugin.Matches(logger, pluginSpec) {
				logger.Info("new-plugin", lager.Data{"address": pluginSpec.Address, "csi-plugin-name": csiPluginName})

				nodePlugin := p.csiShim.NewNodeClient(conn)
				_, err = identityPlugin.Probe(context.TODO(), &csi.ProbeRequest{})
				if err != nil {
					logger.Info("probe-node-unresponsive", lager.Data{"name": csiPluginSpec.Name, "address": csiPluginSpec.Address})
					continue
				}

				plugin := NewCsiPlugin(nodePlugin, pluginSpec, p.grpcShim, p.csiShim, p.osShim, p.csiMountRootDir, oshelper.NewOsHelper())
				plugins[csiPluginName] = plugin
			} else {
				logger.Info("discovered-plugin-ignored", lager.Data{"address": pluginSpec.Address})
				plugins[csiPluginName] = existingPlugin
			}
		}
	}
	return plugins, nil
}
