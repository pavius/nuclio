/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"time"

	nuclioioclient "github.com/nuclio/nuclio/pkg/platform/kube/client/clientset/versioned"
	"github.com/nuclio/nuclio/pkg/platform/kube/functionres"
	"github.com/nuclio/nuclio/pkg/platformconfig"

	"github.com/nuclio/errors"
	"github.com/nuclio/logger"
	"github.com/v3io/version-go"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	logger                logger.Logger
	namespace             string
	kubeClientSet         kubernetes.Interface
	nuclioClientSet       nuclioioclient.Interface
	functionresClient     functionres.Client
	imagePullSecrets      string
	functionOperator      *functionOperator
	projectOperator       *projectOperator
	functionEventOperator *functionEventOperator
	cronJobMonitoring     *CronJobMonitoring
	platformConfiguration *platformconfig.Config
}

func NewController(parentLogger logger.Logger,
	namespace string,
	imagePullSecrets string,
	kubeClientSet kubernetes.Interface,
	nuclioClientSet nuclioioclient.Interface,
	functionresClient functionres.Client,
	resyncInterval time.Duration,
	cronJobStaleResourcesCleanupInterval time.Duration,
	platformConfiguration *platformconfig.Config,
	functionOperatorNumWorkers int,
	functionEventOperatorNumWorkers int,
	projectOperatorNumWorkers int) (*Controller, error) {
	var err error

	// replace "*" with "", which is actually "all" in kube-speak
	if namespace == "*" {
		namespace = ""
	}

	newController := &Controller{
		logger:                parentLogger,
		namespace:             namespace,
		imagePullSecrets:      imagePullSecrets,
		kubeClientSet:         kubeClientSet,
		nuclioClientSet:       nuclioClientSet,
		functionresClient:     functionresClient,
		platformConfiguration: platformConfiguration,
	}

	newController.logger.DebugWith("Read configuration",
		"platformConfig", newController.platformConfiguration,
		"version", version.Get())

	// set ourselves as the platform configuration provider of the function resource client (it needs it to do
	// stuff when creating stuff)
	functionresClient.SetPlatformConfigurationProvider(newController)

	// create a function operator
	newController.functionOperator, err = newFunctionOperator(parentLogger,
		newController,
		&resyncInterval,
		imagePullSecrets,
		functionresClient,
		functionOperatorNumWorkers)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create functions operator")
	}

	// create a function event operator
	newController.functionEventOperator, err = newFunctionEventOperator(parentLogger,
		newController,
		&resyncInterval,
		functionEventOperatorNumWorkers)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create function event operator")
	}

	// create a project operator
	newController.projectOperator, err = newProjectOperator(parentLogger,
		newController,
		&resyncInterval,
		projectOperatorNumWorkers)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create project operator")
	}

	// create cron job monitoring
	if platformConfiguration.CronTriggerCreationMode == platformconfig.KubeCronTriggerCreationMode {
		newController.cronJobMonitoring = NewCronJobMonitoring(parentLogger,
			newController,
			&cronJobStaleResourcesCleanupInterval)
	}

	return newController, nil
}

func (c *Controller) Start() error {
	c.logger.InfoWith("Starting", "namespace", c.namespace)

	// start the function operator
	if err := c.functionOperator.start(); err != nil {
		return errors.Wrap(err, "Failed to start function operator")
	}

	// start the project operator
	if err := c.projectOperator.start(); err != nil {
		return errors.Wrap(err, "Failed to start project operator")
	}

	// start the function event operator
	if err := c.functionEventOperator.start(); err != nil {
		return errors.Wrap(err, "Failed to start function event operator")
	}

	if c.cronJobMonitoring != nil {

		// start cron job monitoring
		c.cronJobMonitoring.start()
	}

	return nil
}

func (c *Controller) GetPlatformConfiguration() *platformconfig.Config {
	return c.platformConfiguration
}
