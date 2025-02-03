package service

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/hcpss-banderson/orikal/model"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"regexp"
	"strings"
)

type KamalService struct {
	projectDir string
	configFile string
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func NewKamalService(projectDir, configFile string) *KamalService {
	return &KamalService{projectDir, configFile}
}

func (k *KamalService) GetRoles() []string {
	logs := k.RunContainer([]string{"config", "--config-file", k.configFile})
	var out []byte
	buffer := make([]byte, 256)
	for {
		n, readerr := logs.Read(buffer)
		if readerr == nil || readerr == io.EOF {
			out = append(out, buffer[:n]...)
		} else {
			panic(readerr)
		}

		if readerr == io.EOF {
			break
		}
	}

	src := strings.NewReader(string(out))
	stdout := &bytes.Buffer{}
	errout := &bytes.Buffer{}
	stdcopy.StdCopy(stdout, errout, src)
	content := stdout.Bytes()

	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(content, m)
	if err != nil {
		panic(err)
	}

	var roles []string
	for _, role := range m[":roles"].([]interface{}) {
		roles = append(roles, role.(string))
	}

	return roles
}

func (k *KamalService) RunContainer(command []string) io.ReadCloser {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	reader, err := apiClient.ImagePull(context.Background(), "banderson/kamal:latest", image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	io.Copy(io.Discard, reader)

	cont, err := apiClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        "banderson/kamal:latest",
			Env:          []string{"SSH_AUTH_SOCK=/run/host-services/ssh-auth.sock"},
			Cmd:          command,
			AttachStdout: true,
		},
		&container.HostConfig{
			Binds: []string{
				os.Getenv("HOME") + "/.ssh/id_rsa:/root/.ssh/id_rsa",
				"/run/host-services/ssh-auth.sock:/run/host-services/ssh-auth.sock",
				"/var/run/docker.sock:/var/run/docker.sock",
				k.projectDir + ":/workdir",
			},
			AutoRemove: true,
		},
		nil,
		&ocispec.Platform{OS: "linux", Architecture: "arm64"},
		"orikal",
	)
	if err != nil {
		panic(err)
	}

	err = apiClient.ContainerStart(context.Background(), cont.ID, container.StartOptions{})
	if err != nil {
		panic(err)
	}

	logs, err := apiClient.ContainerLogs(context.Background(), cont.ID, container.LogsOptions{
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		panic(err)
	}

	return logs
}

func (k *KamalService) AppExec(command string, acronymChan chan string, dataChan chan model.Payload) {
	logs := k.RunContainer([]string{"app", "exec", command, "--config-file", k.configFile, "--reuse"})

	var myout []byte
	acronyms := []string{}
	buffer := make([]byte, 256)
	for {
		n, readerr := logs.Read(buffer)
		if readerr == nil || readerr == io.EOF {
			myout = append(myout, buffer[:n]...)
			reAcronym := regexp.MustCompile(`[a-z]+-schools-([a-z]{2,5})-[a-z0-9]{40} drush ms`)
			matches := reAcronym.FindAllStringSubmatch(string(myout), -1)
			for _, match := range matches {
				acronym := match[1]
				if !stringInSlice(acronym, acronyms) {
					acronyms = append(acronyms, acronym)
					acronymChan <- acronym
				}
			}
		} else {
			panic(readerr)
		}

		if readerr == io.EOF {
			close(acronymChan)
			break
		}
	}

	src := strings.NewReader(string(myout))
	stddest := &bytes.Buffer{}
	errdest := &bytes.Buffer{}
	stdcopy.StdCopy(stddest, errdest, src)

	content := stddest.String()
	stringSlice := strings.Split(content, "Running docker exec")
	stringSlice = stringSlice[1:]
	for _, v := range stringSlice {
		reAcronym := regexp.MustCompile(`[a-z]+-schools-([a-z]{2,5})-[a-z0-9]{40} drush ms`)
		matches := reAcronym.FindStringSubmatch(v)

		acronym := matches[1]
		reJson := regexp.MustCompile(`(?ms)^\[(.*?)\]`)
		matchesJson := reJson.FindStringSubmatch(v)
		jsonString := "[" + matchesJson[1] + "]"

		dataChan <- model.Payload{acronym, jsonString}
	}

	close(dataChan)
}
