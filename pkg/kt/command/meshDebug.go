package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/connect"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	"github.com/alibaba/kt-connect/pkg/kt/command/mesh"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

// NewMeshDebugCommand return new mesh command
func NewMeshDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meshDebug",
		Short: "Redirect marked requests of specified kubernetes service to local",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name of service to mesh is required")
			} else if len(args) > 1 {
				return fmt.Errorf("too many service name are spcified (%s), should be one", strings.Join(args, ","))
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return MeshDebug(args[0])
		},
		Example: "et meshDebug <service-name> [command options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(true))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().MeshDebug, opt.MeshDebugFlags())
	return cmd
}

//Mesh exchange kubernetes workload
func MeshDebug(resourceName string) error {
	*opt.Get().Connect = opt.Get().MeshDebug.ConnectOptions
	*opt.Get().Mesh = opt.Get().MeshDebug.MeshOptions
	ch, err := general.SetupProcess(util.ComponentMeshDebug)
	if err != nil {
		return err
	}
	if !opt.Get().Connect.SkipCleanup {
		go silenceCleanup()
	}

	log.Info().Msgf("Using %s mode", opt.Get().Connect.ConnectMode)
	if opt.Get().Connect.ConnectMode == util.ConnectModeTun2Socks {
		err = connect.ByTun2Socks()
	} else if opt.Get().Connect.ConnectMode == util.ConnectModeShuttle {
		err = connect.BySshuttle()
	} else {
		err = fmt.Errorf("invalid connect mode: '%s', supportted mode are %s, %s", opt.Get().Connect.ConnectMode,
			util.ConnectModeTun2Socks, util.ConnectModeShuttle)
	}
	if err != nil {
		return err
	}
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" All looks good, now you can access to resources in the kubernetes cluster")
	log.Info().Msg("---------------------------------------------------------------")

	if opt.Get().Mesh.SkipPortChecking {
		if port := util.FindBrokenLocalPort(opt.Get().Mesh.Expose); port != "" {
			return fmt.Errorf("no application is running on port %s", port)
		}
	}

	// Get service to mesh
	svc, err := general.GetServiceByResourceName(resourceName, opt.Get().Global.Namespace)
	if err != nil {
		return err
	}

	if port := util.FindInvalidRemotePort(opt.Get().Mesh.Expose, general.GetTargetPorts(svc)); port != "" {
		return fmt.Errorf("target port %s not exists in service %s", port, svc.Name)
	}

	log.Info().Msgf("Using %s mode", opt.Get().Mesh.MeshMode)
	if opt.Get().Mesh.MeshMode == util.MeshModeManual {
		err = mesh.ManualMesh(svc)
	} else if opt.Get().Mesh.MeshMode == util.MeshModeAuto {
		err = mesh.AutoMesh(svc)
	} else {
		err = fmt.Errorf("invalid mesh method '%s', supportted are %s, %s", opt.Get().Mesh.MeshMode,
			util.MeshModeAuto, util.MeshModeManual)
	}
	if err != nil {
		return err
	}

	// watch background process, clean the workspace and exit if background process occur exception
	s := <-ch
	log.Info().Msgf("Terminal Signal is %s", s)
	return nil
}
