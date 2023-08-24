// Copyright (c) 2023, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.
package singularity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	scslibclient "github.com/sylabs/scs-library-client/client"
	"github.com/sylabs/singularity/pkg/sylog"
)

// RemoteGetLoginPassword retrieves cli token from oci library shim
func RemoteGetLoginPassword(config *scslibclient.Config) (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	path := "/v1/rbac/users/current"
	endPoint := config.BaseURL + path

	req, err := http.NewRequest(http.MethodGet, endPoint, nil)
	if err != nil {
		return "", fmt.Errorf("error creating new request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", config.AuthToken))
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		sylog.Debugf("Status Code: %v", res.StatusCode)
		if res.StatusCode == http.StatusUnauthorized {
			return "", fmt.Errorf("must be logged in to retrieve token")
		}
		return "", fmt.Errorf("status is not ok: %v", res.StatusCode)
	}

	var u oUser
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return "", fmt.Errorf("error decoding json response: %v", err)
	}

	if u.OidcUserMeta.Secret == "" {
		return "", fmt.Errorf("user does not have cli token set")
	}

	return u.OidcUserMeta.Secret, nil
}

type oidcUserMeta struct {
	Secret string `json:"secret"`
}

type oUser struct {
	OidcUserMeta oidcUserMeta `json:"oidc_user_meta"`
}
