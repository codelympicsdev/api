package database

// APIClient is the data stored about an api client in the database
type APIClient struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	RedirectURL string   `json:"redirect_url"`
	Trusted     bool     `json:"trusted"`
	Scopes      []string `json:"scopes"`

	Secret string `json:"secret"`
}

// GetAPIClientByID an api client from the database by id
func GetAPIClientByID(id string) (*APIClient, error) {
	data, err := db.Collection("api_clients").Document(id).Get()
	if err != nil {
		return nil, err
	}

	var client = new(APIClient)
	err = data.DataTo(client)
	if err != nil {
		return nil, err
	}

	client.ID = data.Document.ID

	return client, nil
}

// Save api client data to the database
func (client *APIClient) Save() error {
	if client.ID == "" {
		resp, err := db.Collection("api_clients").Add(client)
		if err == nil {
			client.ID = resp.ID
		}
		return err
	}

	return db.Collection("api_clients").Document(client.ID).Set(client)
}
