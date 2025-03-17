package cloudgov

import (
	"fmt"
)

// Stuff we'll need to implement, for ref
//
// mapRoute()
//
// addNetworkPolicy()
// removeNetworkPolicy()
type ClientAPI interface {
	connect(url string, creds *Creds) error

	appGet(id string) (*App, error)
	appPush(m *AppManifest) (*App, error)
	appDelete(id string) error
	appsList() (apps []*App, err error)
}

type CredsGetter interface {
	getCreds() (*Creds, error)
}

type Opts struct {
	CredsGetter
	Creds *Creds

	APIRootURL string
}

type Client struct {
	ClientAPI
	*Opts
}

// TODO: we should pull this out of VCAP_APPLICATION
const apiRootURLDefault = "https://api.fr.cloud.gov"

func New(i ClientAPI, o *Opts) (*Client, error) {
	if o == nil {
		o = &Opts{CredsGetter: EnvCredsGetter{}}
	}
	cg := &Client{ClientAPI: i, Opts: o}
	return cg.Connect()
}

func (c *Client) apiRootURL() string {
	if c.APIRootURL == "" {
		return apiRootURLDefault
	}
	return c.APIRootURL
}

func (c *Client) creds() (*Creds, error) {
	if c.Creds.isEmpty() {
		return c.getCreds()
	}
	return c.Creds, nil
}

func (c *Client) Connect() (*Client, error) {
	creds, err := c.creds()
	if err != nil {
		return nil, err
	}
	if err := c.connect(c.apiRootURL(), creds); err != nil {
		return nil, err
	}
	return c, nil
}

type App struct {
	Name  string // i.e., container ID
	State string
}

func (c *Client) AppGet(id string) (*App, error) {
	return c.appGet(id)
}

func (c *Client) AppDelete(id string) error {
	return c.appDelete(id)
}

func (c *Client) AppsList() ([]*App, error) {
	return c.appsList()
}

// TODO: this abstraction might belong in /cmd,
// unless it can be further generalized to all pushes
func (c *Client) ServicePush(manifest *AppManifest) (*App, error) {
	containerID := manifest.Name

	// check for an old instance of the service, delete if found
	app, err := c.AppGet(containerID)
	if err != nil {
		return nil, fmt.Errorf("error checking for existing service (%v): %w", containerID, err)
	}
	if app != nil {
		err = c.AppDelete(containerID)
	}
	if err != nil {
		return nil, fmt.Errorf("error deleting existing service (%v): %w", containerID, err)
	}

	return c.appPush(manifest)
}

func (c *Client) ServicesPush(manifests []*AppManifest) ([]*App, error) {
	if len(manifests) < 1 {
		return nil, nil
	}

	apps := make([]*App, len(manifests))

	for i, s := range manifests {
		app, err := c.ServicePush(s)
		if err != nil {
			return nil, err
		}
		apps[i] = app
	}

	return apps, nil
}
