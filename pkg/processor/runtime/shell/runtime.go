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

package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/nuclio/nuclio/pkg/errors"
	"github.com/nuclio/nuclio/pkg/processor/runtime"

	"github.com/nuclio/nuclio-sdk"
)

type shell struct {
	*runtime.AbstractRuntime
	configuration *runtime.Configuration
	handlerPath   string
	env           []string
	ctx           context.Context
}

func NewRuntime(parentLogger nuclio.Logger, configuration *runtime.Configuration) (runtime.Runtime, error) {

	runtimeLogger := parentLogger.GetChild("shell")

	// create base
	abstractRuntime, err := runtime.NewAbstractRuntime(runtimeLogger, configuration)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create abstract runtime")
	}

	// create the command string
	newShellRuntime := &shell{
		AbstractRuntime: abstractRuntime,
		ctx:             context.Background(),
		configuration:   configuration,
	}

	// update it with some stuff so that we don't have to do this each invocation
	newShellRuntime.handlerPath = newShellRuntime.getHandlerPath()
	newShellRuntime.env = newShellRuntime.getEnvFromConfiguration()

	// set permissions of handler such that if it wasn't executable before, it's executable now
	os.Chmod(newShellRuntime.handlerPath, 0755)

	return newShellRuntime, nil
}

func (s *shell) ProcessEvent(event nuclio.Event, functionLogger nuclio.Logger) (interface{}, error) {
	s.Logger.DebugWith("Executing shell",
		"name", s.configuration.Meta.Name,
		"version", s.configuration.Spec.Version,
		"eventID", event.GetID(),
		"handlerPath", s.handlerPath)

	// create a timeout context
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	// create a command
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", s.handlerPath)

	cmd.Stdin = strings.NewReader(string(event.GetBody()))

	// set the command env
	cmd.Env = s.env

	// add event stuff to env
	cmd.Env = append(cmd.Env, s.getEnvFromEvent(event)...)

	// run the command
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to run shell command")
	}

	s.Logger.DebugWith("Shell executed",
		"out", string(out),
		"eventID", event.GetID())

	return out, nil
}

func (s *shell) getHandlerPath() string {
	handler := s.configuration.Spec.Handler
	targetName := strings.Split(handler, ":")[0]

	shellHandlerPath := os.Getenv("NUCLIO_SHELL_HANDLER_DIR")
	if shellHandlerPath == "" {
		shellHandlerPath = "/opt/nuclio/"
	}

	command := path.Join(shellHandlerPath, targetName)

	return command
}

func (s *shell) getEnvFromConfiguration() []string {
	return []string{
		fmt.Sprintf("NUCLIO_FUNCTION_NAME=%s", s.configuration.Meta.Name),
		fmt.Sprintf("NUCLIO_FUNCTION_DESCRIPTION=%s", s.configuration.Spec.Description),
		fmt.Sprintf("NUCLIO_FUNCTION_VERSION=%d", s.configuration.Spec.Version),
	}
}

func (s *shell) getEnvFromEvent(event nuclio.Event) []string {
	return []string{
		fmt.Sprintf("NUCLIO_EVENT_ID=%s", event.GetID()),
		fmt.Sprintf("NUCLIO_EVENT_SOURCE_CLASS=%s", event.GetSource().GetClass()),
		fmt.Sprintf("NUCLIO_EVENT_SOURCE_KIND=%s", event.GetSource().GetKind()),
		fmt.Sprintf("NUCLIO_EVENT_BODY=%s", event.GetBody()),
	}
}
