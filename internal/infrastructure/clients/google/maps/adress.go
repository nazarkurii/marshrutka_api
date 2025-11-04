package googleMaps

import (
	"context"
	"fmt"
	"maryan_api/config"
	rfc7807 "maryan_api/pkg/problem"
	"net/http"
	"strings"
)

func VerifyAdressID(ctx context.Context, client *http.Client, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&key=%s", id, config.GooglePlacesApiKey()),
		strings.NewReader(""),
	)

	if err != nil {
		return rfc7807.Internal("Google Request Composing Error", err.Error())
	}

	_, err = client.Do(req)

	if err != nil {
		return rfc7807.BadRequest("non-existing-google-places-id", "Non-existing Google Places ID", err.Error())
	}

	return nil

}
