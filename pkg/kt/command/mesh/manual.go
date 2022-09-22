package mesh

import (
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/service/cluster"
	"github.com/alibaba/kt-connect/pkg/kt/util"
	"github.com/rs/zerolog/log"
	coreV1 "k8s.io/api/core/v1"
	"strings"
)

func ManualMesh(svc *coreV1.Service) error {
	service := svc.Name
	namespace := opt.Get().Global.Namespace
	drName := opt.Get().Mesh.DrName
	if drName == "" {
		drName = service
		opt.Get().Mesh.DrName = service
	}
	vsName := opt.Get().Mesh.VsName
	if vsName == "" {
		vsName = service
		opt.Get().Mesh.VsName = service
	}
	if _, err := cluster.Ins().GetVirtualService(vsName, namespace); err != nil {
		return err
	}
	if _, err := cluster.Ins().GetDestinationRule(drName, namespace); err != nil {
		return err
	}

	meshKey, meshVersion := getVersion(opt.Get().Mesh.VersionMark)
	shadowPodName := svc.Name + util.MeshPodInfix + meshVersion
	labels := getMeshLabels(meshKey, meshVersion, svc)
	annotations := make(map[string]string)
	if err := general.CreateShadowAndInbound(shadowPodName, opt.Get().Mesh.Expose, labels,
		annotations, general.GetTargetPorts(svc)); err != nil {
		return err
	}
	log.Info().Msgf("update Istio rule start...")

	cluster.Ins().PatchDestinationRule(drName, namespace, "add", meshKey, meshVersion)
	opt.Store.DestinationRulePatch = true
	cluster.Ins().PatchVirtualService(vsName, svc.Name, namespace, "add", meshKey, meshVersion)
	opt.Store.VirtualServicePatch = true
	log.Info().Msg("update Istio rule finish...")
	cluster.Ins().PatchDestinationRule(drName, namespace, "add", meshKey, meshVersion)
	opt.Store.DestinationRulePatch = true
	cluster.Ins().PatchVirtualService(vsName, svc.Name, namespace, "add", meshKey, meshVersion)
	opt.Store.VirtualServicePatch = true
	log.Info().Msg("update Istio rule finish...")
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" Now you can access your service by header '%s: %s' ", strings.ToUpper(meshKey), meshVersion)
	log.Info().Msg("---------------------------------------------------------------")

	return nil
}

func getMeshLabels(meshKey, meshVersion string, svc *coreV1.Service) map[string]string {
	labels := map[string]string{}
	if svc != nil {
		for k, v := range svc.Spec.Selector {
			labels[k] = v
		}
	}
	labels[util.KtRole] = util.RoleMeshShadow
	labels[meshKey] = meshVersion
	return labels
}
