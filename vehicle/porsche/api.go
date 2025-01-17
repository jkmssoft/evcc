package porsche

import (
	"fmt"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
	"github.com/evcc-io/evcc/util/transport"
	"golang.org/x/oauth2"
)

const (
	ApiURI = "https://api.porsche.com"

	PairingComplete = "PAIRINGCOMPLETE"
)

// API is an api.Vehicle implementation for Porsche PHEV cars
type API struct {
	*request.Helper
}

// NewAPI creates a new vehicle
func NewAPI(log *util.Logger, identity oauth2.TokenSource) *API {
	v := &API{
		Helper: request.NewHelper(log),
	}

	v.Client.Transport = &transport.Decorator{
		Base: &oauth2.Transport{
			Source: identity,
			Base:   v.Client.Transport,
		},
		Decorator: transport.DecorateHeaders(map[string]string{
			"apikey": OAuth2Config.ClientID,
		}),
	}

	return v
}

// Vehicles implements the vehicle list response
func (v *API) Vehicles() ([]Vehicle, error) {
	var res []Vehicle
	uri := fmt.Sprintf("%s/core/api/v3/de/de_DE/vehicles", ApiURI)
	err := v.GetJSON(uri, &res)
	return res, err
}

// PairingStatus implements the vehicle pairing status response
func (v *API) PairingStatus(vin string) (VehiclePairingResponse, error) {
	var res VehiclePairingResponse
	uri := fmt.Sprintf("%s/core/api/v3/de/de_DE/vehicles/%s/pairing", ApiURI, vin)
	err := v.GetJSON(uri, &res)
	return res, err
}

// Status implements the vehicle status response
func (v *API) Status(vin string) (StatusResponse, error) {
	var res StatusResponse
	uri := fmt.Sprintf("%s/vehicle-data/de/de_DE/status/%s", ApiURI, vin)
	err := v.GetJSON(uri, &res)
	return res, err
}
