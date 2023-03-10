package kubectl

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/manifests"
	"github.com/alibaba/kt-connect/pkg/kt/service/cluster"
	"github.com/rs/zerolog/log"
	"time"
)

func InstallEtck() error {
	err := waitCertManagerReady(120, 0)
	if err != nil {
		return err
	}

	err = installEtckController()
	if err != nil {
		return err
	}
	return nil
}

func waitCertManagerReady(timeoutSec int, times int) error {
	const interval = 6
	if times > timeoutSec/interval {
		return fmt.Errorf("cert-manager failed to start")
	}

	certDeploy, _ := cluster.Ins().GetAllDeploymentInNamespace("cert-manager")
	if len(certDeploy.Items) == 0 {
		log.Info().Msgf("install cert-manager start")
		fileRepo := manifests.NewRepo("etck")
		data, err := fileRepo.ReadFile("etck/templates/cert-manager.yaml")
		err = cluster.Ins().CreateByFile(data)
		if err != nil {
			return err
		}
		time.Sleep(90 * time.Second)
		waitCertManagerReady(timeoutSec, times+1)
	} else if len(certDeploy.Items) == 3 {
		if certDeploy.Items[0].Status.ReadyReplicas == 1 &&
			certDeploy.Items[1].Status.ReadyReplicas == 1 &&
			certDeploy.Items[2].Status.ReadyReplicas == 1 {
			log.Info().Msgf("cert-manager is ready")
			return nil
		}
	} else {
		time.Sleep(interval * time.Second)
		waitCertManagerReady(timeoutSec, times+1)
	}
	return nil
}

func installEtckController() error {
	etckDeploy, _ := cluster.Ins().GetAllDeploymentInNamespace("etck-system")
	if len(etckDeploy.Items) == 1 && etckDeploy.Items[0].Status.ReadyReplicas == 1 {
		return nil
	} else {
		fileRepo := manifests.NewRepo("etck")
		data, err := fileRepo.ReadFile("etck/templates/operator-bundle.yaml")
		if err != nil {
			return err
		}
		log.Info().Msgf("install etck-controller start")
		err = cluster.Ins().CreateByFile(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func UninstallEtck() error {
	//etckDeploy, _ := cluster.Ins().GetAllDeploymentInNamespace("etck-system")
	//if len(etckDeploy.Items) == 0 {
	//	return nil
	//} else {

	fileRepo := manifests.NewRepo("etck")
	data, err := fileRepo.ReadFile("etck/templates/operator-bundle.yaml")
	if err != nil {
		return err
	}
	log.Info().Msgf("uninstall etck-controller start")
	err = cluster.Ins().DeleteByFile(data)
	if err != nil {
		return err
	}

	data, err = fileRepo.ReadFile("etck/templates/cert-manager.yaml")
	if err != nil {
		return err
	}
	log.Info().Msgf("uninstall cert-manager start")
	err = cluster.Ins().DeleteByFile(data)
	if err != nil {
		return err
	}
	return nil
}
