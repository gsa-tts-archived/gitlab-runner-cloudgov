package cloudgov

// Stuff we'll need to implement, for ref
//
// mapRoute()
//
// addNetworkPolicy()
// removeNetworkPolicy()
//
// appGet()
// appCmd()
// appPush()
// appDelete()
type ClientAPI interface {
	connect(url string, creds *Creds) error

	appGet(id string) (*App, error)
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
	Id    string
	Name  string
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
