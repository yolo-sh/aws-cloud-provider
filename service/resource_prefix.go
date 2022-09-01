package service

import (
	"github.com/yolo-sh/yolo/entities"
)

func prefixClusterResource(
	clusterNameSlug string,
) func(string) string {

	return func(resourceNameSlug string) string {
		if clusterNameSlug == entities.DefaultClusterName {
			return "yolo-" + resourceNameSlug
		}

		return "yolo-" + clusterNameSlug + "-" + resourceNameSlug
	}
}

func prefixEnvResource(
	clusterNameSlug string,
	envNameSlug string,
) func(string) string {

	return func(resourceNameSlug string) string {
		if clusterNameSlug == entities.DefaultClusterName {
			return "yolo-" + envNameSlug + "-" + resourceNameSlug
		}

		return "yolo-" + clusterNameSlug + "-" + envNameSlug + "-" + resourceNameSlug
	}
}
