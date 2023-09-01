package kubectl

import (
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/manifests"
	"github.com/alibaba/kt-connect/pkg/kt/service/cluster"
	"github.com/rs/zerolog/log"
	"strings"
)

func InstallVenv() error {
	err := InstallVenvController()
	if err != nil {
		return err
	}

	err = InstallVenvOperator()
	if err != nil {
		return err
	}
	return nil
}

func InstallVenvController() error {
	certDeploy, _ := cluster.Ins().GetAllDeploymentInNamespace("et-virtual-environment")
	if len(certDeploy.Items) == 0 {
		log.Info().Msgf("install venv controller start")
		fileRepo := manifests.NewRepo("venv")
		data, err := fileRepo.ReadFile("venv/templates/venv_crd.yaml")
		err = cluster.Ins().CreateByFile(data, "")
		if err != nil {
			return err
		}
		data, err = fileRepo.ReadFile("venv/templates/venv_webhook.yaml")
		if opt.Get().VenvInstall.VenvVersion == "v0.6.0" {
			data = []byte(strings.Replace(string(data), "v0.6.1", "v0.6.0", -1))
		}
		err = cluster.Ins().CreateByFile(data, "")
		if err != nil {
			return err
		}
	}
	log.Info().Msgf("venv controller is ready")
	return nil
}

func InstallVenvOperator() error {
	namespace := opt.Get().Global.Namespace
	if "" == namespace {
		namespace = "default"
	}

	fileRepo := manifests.NewRepo("venv")

	log.Info().Msgf("label namespace environment-tag-injection be enabled")
	_, err := cluster.Ins().PatchNamespaceLabel(namespace, "environment-tag-injection", "enabled")
	if err != nil {
		return err
	}

	log.Info().Msgf("label namespace istio-injection be enabled")
	_, err = cluster.Ins().PatchNamespaceLabel(namespace, "istio-injection", "enabled")
	if err != nil {
		return err
	}

	log.Info().Msgf("install venv service_account start")
	data, err := fileRepo.ReadFile("venv/templates/venv_service_account.yaml")
	if err != nil {
		return err
	}
	err = cluster.Ins().CreateByFile(data, namespace)
	if err != nil {
		return err
	}

	log.Info().Msgf("install venv operator start")
	data, err = fileRepo.ReadFile("venv/templates/venv_operator.yaml")
	if err != nil {
		return err
	}
	if opt.Get().VenvInstall.VenvVersion == "v0.6.0" {
		data = []byte(strings.Replace(string(data), "v0.6.1", "v0.6.0", -1))
	}
	err = cluster.Ins().CreateByFile(data, namespace)
	if err != nil {
		return err
	}

	log.Info().Msgf("install venv resource start")
	data, err = fileRepo.ReadFile("venv/templates/venv_virtualenv.yaml")
	data = []byte(strings.Replace(string(data), "{{headerKey}}", opt.Get().VenvInstall.EnvHeader, -1))
	if err != nil {
		return err
	}
	err = cluster.Ins().CreateByFile(data, namespace)
	if err != nil {
		return err
	}

	return nil
}

func UninstallVenvOperator() error {
	//etckDeploy, _ := cluster.Ins().GetAllDeploymentInNamespace("etck-system")
	//if len(etckDeploy.Items) == 0 {
	//	return nil
	//} else {
	namespace := opt.Get().Global.Namespace
	fileRepo := manifests.NewRepo("venv")

	log.Info().Msgf("uninstall venv resource start")
	data, err := fileRepo.ReadFile("venv/templates/venv_virtualenv.yaml")
	if err != nil {
		return err
	}
	_ = cluster.Ins().DeleteByFile(data, namespace)
	//if err != nil {
	//	return err
	//}

	log.Info().Msgf("uninstall venv service_account start")
	data, err = fileRepo.ReadFile("venv/templates/venv_service_account.yaml")
	if err != nil {
		return err
	}
	err = cluster.Ins().DeleteByFile(data, namespace)
	if err != nil {
		return err
	}

	log.Info().Msgf("uninstall venv operator start")
	data, err = fileRepo.ReadFile("venv/templates/venv_operator.yaml")
	if err != nil {
		return err
	}
	err = cluster.Ins().DeleteByFile(data, namespace)
	if err != nil {
		return err
	}

	log.Info().Msgf("label namespace environment-tag-injection be disabled")
	_, err = cluster.Ins().PatchNamespaceLabel(namespace, "environment-tag-injection", "disabled")
	if err != nil {
		return err
	}
	return nil
}

func UninstallVenvController() error {
	//etckDeploy, _ := cluster.Ins().GetAllDeploymentInNamespace("etck-system")
	//if len(etckDeploy.Items) == 0 {
	//	return nil
	//} else {
	fileRepo := manifests.NewRepo("venv")

	log.Info().Msgf("uninstall venv resource start")
	data, err := fileRepo.ReadFile("venv/templates/venv_webhook.yaml")
	if err != nil {
		return err
	}
	err = cluster.Ins().DeleteByFile(data, "")
	if err != nil {
		return err
	}

	log.Info().Msgf("uninstall venv service_account start")
	data, err = fileRepo.ReadFile("venv/templates/venv_crd.yaml")
	if err != nil {
		return err
	}
	err = cluster.Ins().DeleteByFile(data, "")
	if err != nil {
		return err
	}

	return nil
}
