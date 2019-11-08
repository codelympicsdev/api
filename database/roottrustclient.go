package database

// RootTrustClient is a client that can talk to the authentication endpoints
type RootTrustClient struct {
	ID     string `bson:"id"`
	Name   string `bson:"name"`
	Secret string `bson:"secret"`
}

// GetRootTrustClientByID a root trust client from the database by id
func GetRootTrustClientByID(id string) (*RootTrustClient, error) {
	data, err := db.Collection("root_trust_clients").Document(id).Get()
	if err != nil {
		return nil, err
	}

	var client = new(RootTrustClient)
	err = data.DataTo(client)
	if err != nil {
		return nil, err
	}

	client.ID = data.Document.ID

	return client, nil
}

// Save api client data to the database
func (client *RootTrustClient) Save() error {
	if client.ID == "" {
		resp, err := db.Collection("root_trust_clients").Add(client)
		if err == nil {
			client.ID = resp.ID
		}
		return err
	}

	return db.Collection("root_trust_clients").Document(client.ID).Set(client)
}
