package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/types"
)

type CampaignApi struct {
	campaignService  services.CampaignService
	characterService services.CharacterService
}

func NewCampaignApi(cs services.CampaignService, chars services.CharacterService) CampaignApi {
	return CampaignApi{
		campaignService:  cs,
		characterService: chars}
}

func (ca *CampaignApi) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var cr types.CreateCampaignRecordRequest
	err := decodeJSONBody(w, r, &cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	clientId := r.Header["X-Archivist-Client-Id"]
	if len(clientId) == 0 {
		http.Error(w, errors.New("no client id supplied").Error(), http.StatusBadRequest)
		return
	}
	//should always be the first Client Id. Need to find a way to expose a possible mismatch from multipe client Id headers
	nc, err := ca.campaignService.CreateCampaign(cr, clientId[0])
	if err != nil {
		if strings.Contains(err.Error(), "Index already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, nc)
}

func (ca *CampaignApi) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	var cr types.CampaignRecord
	err := decodeJSONBody(w, r, &cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nc, err := ca.campaignService.UpdateCampaign(&cr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, nc)

}

func (ca *CampaignApi) ListCampaigns(w http.ResponseWriter, r *http.Request) {
	applyCorsHeaders(w)
	var arr []*types.CampaignRecord
	var err error
	arr, err = ca.campaignService.ListCampaigns()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, arr)

}

func (ca *CampaignApi) GetCampaign(w http.ResponseWriter, r *http.Request) {
	applyCorsHeaders(w)
	id := r.PathValue("campaignId")

	campaign, err := ca.campaignService.GetCampaign(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	campaign.Characters, err = ca.characterService.GetCharactersForCampaign(campaign)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendGoodResponseWithObject(w, campaign)
}

func (ca *CampaignApi) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("campaignId")
	campaign, err := ca.campaignService.DeleteCampaign(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	sendGoodResponseWithObject(w, campaign)
}

func (ca *CampaignApi) EditCampaignPassword(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("campaignId")
	password := r.URL.Query()["password"]
	_, err := ca.campaignService.UpdateCampaignPassword(id, password[1])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendSuccessResponse(w, "password successfuly updated")

}

func (ua CampaignApi) ValidatePassword(w http.ResponseWriter, r *http.Request) {
	var incomingReques *types.APIRequest
	err := decodeJSONBody(w, r, incomingReques)
	if err != nil {
		return
	}

	isValid, err := ua.userService.ValidateApiUser(incomingReques.Auth.ProvidedClientId, incomingReques.Auth.ProvidedAPIKey)
	if !isValid {
		http.Error(w, errors.New("client not valid").Error(), http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (ua CampaignApi) RequirePassword(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("campaignId")
		password := r.URL.Query()["password"]
		givenClientId := r.Header[http.CanonicalHeaderKey("X-Archivist-Client-Id")][0]
		givenAPIKey := r.Header[http.CanonicalHeaderKey("X-Archivist-API-Key")][0]

		isValid, err := ua.userService.ValidateApiUser(givenClientId, givenAPIKey)
		if !isValid {
			http.Error(w, errors.New("user not valid").Error(), http.StatusForbidden)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		next.ServeHTTP(w, r)
	})

}
