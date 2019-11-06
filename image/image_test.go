package image

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
)

func generateContainerStatus(name, image, imageID string) v1.ContainerStatus {
	cs := v1.ContainerStatus{
		Name:    name,
		Image:   image,
		ImageID: imageID,
	}
	return cs
}

var imageTable = []struct {
	imageID     string
	throwsError bool

	expectedHost       string
	expectedNamespace  string
	expectedRepository string
	expectedDigest     string
}{
	{
		"docker-pullable://quay.io/quay/redis@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		false,
		"quay.io",
		"quay",
		"redis",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
	},
	{
		"docker-pullable://nginx@sha256:0d71ff22db29a08ac7399d1b35b0311c5b0cbe68d878993718275758811f652a",
		false,
		"docker.io",
		"library",
		"nginx",
		"sha256:0d71ff22db29a08ac7399d1b35b0311c5b0cbe68d878993718275758811f652a",
	},
	{
		"docker-pullable://quay/redis@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		false,
		"docker.io",
		"quay",
		"redis",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
	},
	{
		"docker-pullable://quay/redis--test@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		false,
		"docker.io",
		"quay",
		"redis--test",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
	},
	{
		"quay/redis--test@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		true,
		"docker.io",
		"quay",
		"redis--test",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
	},
}

func TestParseImageID(t *testing.T) {
	for _, tt := range imageTable {
		var image = &Image{
			Host:       tt.expectedHost,
			Namespace:  tt.expectedNamespace,
			Repository: tt.expectedRepository,
			Digest:     tt.expectedDigest,
		}
		parsedImageID, err := ParseImageID(tt.imageID)
		if tt.throwsError {
			assert.Error(t, err)
		} else if !reflect.DeepEqual(image, parsedImageID) {
			t.Errorf("Incorrectly parsed %s as %+v: %s", tt.imageID, parsedImageID, err)
		}
	}
}

var containerStatusTable = []struct {
	// Container status
	name    string
	image   string
	imageID string

	// Expected values
	containername string
	host          string
	namespace     string
	repository    string
	digest        string
	tag           string
}{
	{
		"redis",
		"quay.io/quay/redis:latest",
		"docker-pullable://quay.io/quay/redis@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",

		"redis",
		"quay.io",
		"quay",
		"redis",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		"latest",
	},
	{
		"nginx",
		"nginx:latest",
		"docker-pullable://nginx@sha256:0d71ff22db29a08ac7399d1b35b0311c5b0cbe68d878993718275758811f652a",

		"nginx",
		"docker.io",
		"library",
		"nginx",
		"sha256:0d71ff22db29a08ac7399d1b35b0311c5b0cbe68d878993718275758811f652a",
		"latest",
	},
	{
		"redis",
		"quay/redis:latest",
		"docker-pullable://quay/redis@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",

		"redis",
		"docker.io",
		"quay",
		"redis",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		"latest",
	},
	{
		"redis",
		"quay/redis--test:latest",
		"docker-pullable://quay/redis--test@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",

		"redis",
		"docker.io",
		"quay",
		"redis--test",
		"sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		"latest",
	},
}

func TestParseContainerStatus(t *testing.T) {
	for _, tt := range containerStatusTable {
		containerStatus := generateContainerStatus(tt.name, tt.image, tt.imageID)
		var image = &Image{
			ContainerName: tt.containername,
			Host:          tt.host,
			Namespace:     tt.namespace,
			Repository:    tt.repository,
			Digest:        tt.digest,
			Tag:           tt.tag,
		}

		parsedContainerStatus, err := ParseContainerStatus(containerStatus)
		if err != nil {
			t.Errorf("%s", err)
		}
		if !reflect.DeepEqual(image, parsedContainerStatus) {
			t.Errorf("Incorrectly parsed %+v as %+v", containerStatus, parsedContainerStatus)
		}
	}
}

var imageIDTable = []struct {
	imageID        string
	tag            string
	expectedString string
}{
	{
		imageID:        "docker-pullable://quay.io/quay/redis@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		tag:            "testTag",
		expectedString: "quay.io/quay/redis:testTag",
	},
	{
		imageID:        "docker-pullable://nginx@sha256:0d71ff22db29a08ac7399d1b35b0311c5b0cbe68d878993718275758811f652a",
		tag:            "testTag",
		expectedString: "nginx:testTag",
	},
	{
		imageID:        "docker-pullable://quay/redis@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		tag:            "testTag",
		expectedString: "quay/redis:testTag",
	},
	{
		imageID:        "docker-pullable://quay/redis--test@sha256:94033a42da840b970fd9d2b04dae5fec56add2714ca674a758d030ce5acba27e",
		tag:            "testTag",
		expectedString: "quay/redis--test:testTag",
	},
}

func TestStringIDRepresentations(t *testing.T) {
	for _, tt := range imageIDTable {
		image, _ := ParseImageID(tt.imageID)
		image.Tag = tt.tag
		expectedImageID := strings.TrimPrefix(tt.imageID, "docker-pullable://")
		assert.Equal(t, expectedImageID, image.IDString())
		assert.Equal(t, tt.expectedString, image.String())
	}
}
