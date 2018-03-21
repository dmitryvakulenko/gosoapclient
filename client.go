package gosoapclient

type Client struct {
	url             string
	typesNamespaces map[string]string
	namespacesAlias map[string]string
}

func NewClient(url string, typesNs, nsAlias map[string]string) *Client {
	return &Client{
		url:             url,
		typesNamespaces: typesNs,
		namespacesAlias: nsAlias}
}


