package utils

import (
	"omega_launcher/plantform"
	"path/filepath"
)

func IsDocker() bool {
	return plantform.GetPlantform() == plantform.Linux_x86_64 && IsFile(filepath.Join("ome", "launcher_liliya"))
}
