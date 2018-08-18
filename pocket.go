package pocket

type PocketClient struct {
	AccessToken string
	ConsumerKey string
	URL         string
	Username    string
}

type PocketClientOptions struct {
	AccessToken string
	ConsumerKey string
	URL         string
}

func NewPocketClient(o *PocketClientOptions) *PocketClient {
	// TODO check for ConsumerKey
	if o.URL == "" {
		o.URL = "https://getpocket.com"
	}
	return &PocketClient{
		URL:         o.URL,
		ConsumerKey: o.ConsumerKey,
		AccessToken: o.AccessToken,
	}
}
