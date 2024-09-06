// A generated module for HelloDagger functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/hello-dagger/internal/dagger"
	"fmt"
	"math"
	"math/rand/v2"
)

type HelloDagger struct{}

// Returns a container that echoes whatever string argument is provided
func (m *HelloDagger) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", "hello", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *HelloDagger) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

// Publish the application container after building and testing it on-the-fly
func (m *HelloDagger) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}
	return m.Build(source).
		Publish(ctx, fmt.Sprintf("ttl.sh/hello-dagger-%.0f", math.Floor(rand.Float64()*10000000))) //#nosec
}

// Build the application container
func (m *HelloDagger) Build(source *dagger.Directory) *dagger.Container {
	// get the build environment container
	// by calling another Dagger Function
	build := m.BuildEnv(source).
		// build the application
		WithExec([]string{"npm", "run", "build"}).
		// get the build output directory
		Directory("./dist")
	// start from a slim NGINX container
	return dag.Container().From("nginx:1.25-alpine").
		// copy the build output directory to the container
		WithDirectory("/usr/share/nginx/html", build).
		// expose the container port
		WithExposedPort(80)
}

// Return the result of running unit tests
func (m *HelloDagger) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	// get the build environment container
	// by calling another Dagger Function
	return m.BuildEnv(source).
		// call the test runner
		WithExec([]string{"npm", "run", "test:unit", "run"}).
		// capture and return the command output
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *HelloDagger) BuildEnv(source *dagger.Directory) *dagger.Container {

	// create a Dagger cache volume for dependencies
	nodeCache := dag.CacheVolume("node")
	return dag.Container().

		// start from a base Node.js container
		From("node:21-slim").

		// add the source code at /src
		WithDirectory("/src", source).

		// mount the cache volume at /src/node_modules
		WithMountedCache("/src/node_modules", nodeCache).

		// change the working directory to /src
		WithWorkdir("/src").

		// run npm install to install dependencies
		WithExec([]string{"npm", "install"})
}
