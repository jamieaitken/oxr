package oxr

// ClientOption allows a Client to be modified.
type ClientOption func(*Client)

// WithAppID sets the Open Exchange App ID to be used.
func WithAppID(appID string) ClientOption {
	return func(client *Client) {
		client.appID = appID
	}
}

// WithDoer allows clients to specify what http.Client is to be used to perform requests.
func WithDoer(doer Doer) ClientOption {
	return func(client *Client) {
		client.doer = doer
	}
}
