package clerk_utils

import (
	"encoding/json"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

var clerkClient clerk.Client

func ClerkClient() clerk.Client {
	return clerkClient
}

func SetClerk() error {
	_clerkClient, err := clerk.NewClient("sk_test_so7Duxkg3VbkzmrpgtY5AR7Hae7r80LlSEjLtByR2j")
	clerkClient = _clerkClient
	return err
}

func UpdateUserPrivateMeta(uid string, data map[string]interface{}) error {
	u, err := clerkClient.Users().Read(uid)
	if err != nil {
		return err
	}

	meta := u.PrivateMetadata.(map[string]interface{})
	for k, v := range data {
		meta[k] = v
	}

	_, err = clerkClient.Users().Update(uid, &clerk.UpdateUser{})

	if err != nil {
		return err
	}

	return nil
}

func GetClerkUser(uid string) (*clerk.User, error) {
	u, err := clerkClient.Users().Read(uid)
	return u, err
}

func UpdateCalenderConnectionOnClerk(uid string, key string) error {
	u, err := clerkClient.Users().Read(uid)
	if err != nil {
		return err
	}

	meta := u.PublicMetadata.(map[string]interface{})
	if meta["calendarConnections"] == nil {
		meta["calendarConnections"] = []interface{}{}
	}
	meta["calendarConnections"] = append(meta["calendarConnections"].([]interface{}), key)
	b, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	_, err = clerkClient.Users().Update(uid, &clerk.UpdateUser{
		PublicMetadata: string(b),
	})
	if err != nil {
		return err
	}

	return nil
}
