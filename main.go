package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	EnvPrefix = "UDOCKER_"
)

func exit(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}

func getEnvOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}

func convertStringSetToFlags(set []string) []string {
	flags := []string{}
	for _, env := range set {
		parsed := strings.SplitN(env, "=", 2)
		if len(parsed) != 2 {
			continue
		}

		key := parsed[0]
		value := parsed[1]

		if len(key) == 1 {
			flags = append(flags, "-"+key)
		} else {
			flags = append(flags, "--"+strings.ReplaceAll(key, "_", "-"))
		}

		if len(value) > 0 {
			flags = append(flags, value)
		}
	}

	return flags
}

func transformArgsWithSet(args, set []string, target string) []string {
	newArgs := []string{}
	argName := ""
	argSkips := 0
	targetArg := target

	if targetArg == "" {
		newArgs = append(newArgs, set...)
		newArgs = append(newArgs, args...)
		return newArgs
	}

	parsed := strings.Split(targetArg, ":")
	if len(parsed) == 1 {
		argName = parsed[0]
	}
	if len(parsed) == 2 {
		argName = parsed[0]
		parsedPos, err := strconv.ParseInt(parsed[1], 10, 64)
		if err != nil {
			exit("invalid arg skips: " + parsed[1])
		}

		argSkips = int(parsedPos)
	}
	if len(parsed) > 2 {
		exit("invalid target arg: " + targetArg)
	}

	skips := 0
	for i, arg := range args {
		if arg == argName && skips == argSkips {
			newArgs = append(newArgs, args[:i+1]...)
			newArgs = append(newArgs, set...)
			newArgs = append(newArgs, args[i+1:]...)
			break
		}
	}

	return newArgs
}

func convertEnvToFlags() []string {
	envs := []string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, EnvPrefix) && !strings.HasPrefix(env, EnvPrefix+"_") {
			envs = append(envs, strings.TrimPrefix(env, EnvPrefix))
		}
	}

	return convertStringSetToFlags(envs)
}

func transformArgsWithEnv(args []string) []string {
	return transformArgsWithSet(args, convertEnvToFlags(), os.Getenv("UDOCKER__TARGETARG"))
}

func isDockerPresent() {
	path, err := exec.LookPath(getEnvOrDefault(EnvPrefix+"_"+"DOCKERCLI", "docker"))
	if err != nil {
		exit("docker binary not found")
	}

	exec, err := os.Executable()
	if err != nil {
		exit("failed to evaluate current binary name")
	}

	if exec == path {
		exit("both current binary and \"docker\" point to the same executable")
	}
}

func docker(args []string) {
	newArgs := transformArgsWithEnv(args)
	cmd := exec.Command(getEnvOrDefault(EnvPrefix+"_"+"DOCKERCLI", "docker"), newArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func main() {
	isDockerPresent()
	docker(os.Args[1:])
}
