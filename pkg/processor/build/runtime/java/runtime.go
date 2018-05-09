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

package java

import (
	"bytes"
	"io"
	"os"
	"path"
	"text/template"

	"github.com/nuclio/nuclio/pkg/common"
	"github.com/nuclio/nuclio/pkg/errors"
	"github.com/nuclio/nuclio/pkg/processor/build/runtime"
)


type java struct {
	*runtime.AbstractRuntime
}

// GetName returns the name of the runtime, including version if applicable
func (j *java) GetName() string {
	return "java"
}

// OnAfterStagingDirCreated will build jar if the source is a Java file
// It will set generatedJarPath field
func (j *java) OnAfterStagingDirCreated(stagingDir string) error {

	// create a build script alongside the user's code. if user provided a script, it'll use that
	return j.createGradleBuildScript(stagingDir)
}

func (j *java) createGradleBuildScript(stagingBuildDir string) error {
	gradleBuildScriptPath := path.Join(stagingBuildDir, "handler", "build.gradle")

	// if user supplied gradle build script - use it
	if common.IsFile(gradleBuildScriptPath) {
		j.Logger.DebugWith("Found user gradle build script, using it", "path", gradleBuildScriptPath)
		return nil
	}

	gradleBuildScriptTemplate, err := template.New("gradleBuildScript").Parse(j.getGradleBuildScriptTemplateContents())
	if err != nil {
		return errors.Wrap(err, "Failed to create gradle build script template")
	}

	buildFile, err := os.Create(gradleBuildScriptPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to create gradle build script file @ %s", gradleBuildScriptPath)
	}

	defer buildFile.Close() // nolint: errcheck

	data := map[string]interface{}{
		"Dependencies": j.FunctionConfig.Spec.Build.Dependencies,
	}

	var gradleBuildScriptTemplateBuffer bytes.Buffer
	err = gradleBuildScriptTemplate.Execute(io.MultiWriter(&gradleBuildScriptTemplateBuffer, buildFile), data)

	j.Logger.DebugWith("Created gradle build script",
		"path", gradleBuildScriptPath,
		"content", gradleBuildScriptTemplateBuffer.String())

	return err
}

func (j *java) getGradleBuildScriptTemplateContents() string {
	return `plugins {
  id 'com.github.johnrengelman.shadow' version '2.0.2'
  id 'java'
}

repositories {
    mavenCentral()
}

shadowJar {
   baseName = 'user-handler'
   classifier = null  // Don't append "all" to jar name
}

task userHandler(dependsOn: shadowJar)
`
}

// GetProcessorDockerfilePath returns the contents of the appropriate Dockerfile, with which we'll build
// the processor image
func (j *java) GetProcessorDockerfileContents() string {
	return `ARG NUCLIO_TAG=latest
ARG NUCLIO_ARCH=amd64
ARG NUCLIO_BASE_IMAGE=openjdk:9-slim

# Supplies processor, handler.jar
FROM nuclio/handler-builder-java-onbuild:${NUCLIO_TAG}-${NUCLIO_ARCH} as builder

# Supplies uhttpc, used for healthcheck
FROM nuclio/uhttpc:latest-amd64 as uhttpc

# From the base image
FROM ${NUCLIO_BASE_IMAGE}

# Copy required objects from the suppliers
COPY --from=builder /home/nuclio/bin/processor /usr/local/bin/processor
COPY --from=builder /home/nuclio/src/wrapper/build/libs/nuclio-java-wrapper.jar /opt/nuclio/nuclio-java-wrapper.jar
COPY --from=uhttpc /home/nuclio/bin/uhttpc /usr/local/bin/uhttpc

# Readiness probe
HEALTHCHECK --interval=1s --timeout=3s CMD /usr/local/bin/uhttpc --url http://localhost:8082/ready || exit 1

# Run processor with configuration and platform configuration
CMD [ "processor", "--config", "/etc/nuclio/config/processor/processor.yaml", "--platform-config", "/etc/nuclio/config/platform/platform.yaml" ]
`
}
