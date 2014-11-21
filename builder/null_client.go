package builder

import "github.com/fsouza/go-dockerclient"

type nullClient struct{}

func (client *nullClient) Client() *docker.Client {
	return nil
}

func (client *nullClient) LatestImageIDByName(name string) (string, error) {
	return name, nil
}

func (client *nullClient) LatestImageIDByTag(tag string) (string, error) {
	return tag, nil
}

func (client *nullClient) LatestImageByRegex(regex string) (*docker.APIImages, error) {
	return nil, nil
}
