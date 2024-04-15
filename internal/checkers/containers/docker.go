package containers

import (
	"context"
	"os/exec"
	"strings"
)

func getRunningContainers(ctx context.Context) (map[string]struct{}, error) {
	cmd := exec.CommandContext(ctx, "docker", "ps", "--format", "'{{.Names}}'")
	output, err := cmd.Output()
	if err != nil {
		return map[string]struct{}{}, err
	}

	result := map[string]struct{}{}
	containers := strings.Split(string(output), "\n")
	for i := range containers {
		result[containers[i]] = struct{}{}
	}
	return result, nil
}
