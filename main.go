package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/floodedrealms/adventure-archivist/api"
	"github.com/floodedrealms/adventure-archivist/commands"
	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/services"
	"github.com/floodedrealms/adventure-archivist/util"
)

func main() {
	flag.Parse()
	flags := flag.Args()

	if len(flags) == 0 {
		fmt.Print("usage: archivist [operation] <args>")
		return
	}
	operation := flag.Arg(0)
	switch operation {

	case "archive":
		commands.ManageCampaign()
	case "create-user":
		if len(flags) == 1 {
			fmt.Println("usage: archivist create-user [friendly-name]")
			return
		}
		friendlyName := flag.Arg(1)
		if len(flags) == 3 {
			commands.CreateUser("api", friendlyName, false)
		} else {
			commands.CreateUser("api", friendlyName, true)
		}
	case "server":
		//	memRepo, err := repository.NewMemoryRepo()
		//	util.CheckErr(err)
		debug := false
		if len(flags) == 2 {
			if flag.Arg(1) == "true" {
				debug = true
			}
		}

		logger := util.NewLogger(debug)

		sqlRepo, err := repository.NewSqliteRepo("archivist.db", logger)
		util.CheckErr(err)

		campaignService := services.NewCampaignService(sqlRepo, logger, context.TODO())
		characterService := services.NewCharacterService(sqlRepo, logger, context.TODO())
		adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())
		userService := services.NewUserService(sqlRepo, *logger)
		campaignActionService := services.NewCampaignActionService(sqlRepo)

		campaignApi := api.NewCampaignApi(campaignService, characterService)
		adventureRecordApi := api.NewAdventureRecordApi(adventureRecordService, characterService)
		characterApi := api.NewCharacterApi(characterService, *campaignActionService)
		userApi := api.NewClientAPI(userService)

		router := http.NewServeMux()

		// Wrap functions
		createCampaign := userApi.RequireValidClient(http.HandlerFunc(campaignApi.CreateCampaign))
		updateCampaign := userApi.RequireValidClient(http.HandlerFunc(campaignApi.UpdateCampaign))
		deleteCampaign := userApi.RequireValidClient(http.HandlerFunc(campaignApi.DeleteCampaign))
		addAdventureToCampaign := userApi.RequireValidClient(http.HandlerFunc(adventureRecordApi.CreateAdventureRecord))
		addCharacterToCampaign := userApi.RequireValidClient(http.HandlerFunc(characterApi.CreateCharacterForCampaign))
		addCampaignActionToCharacter := userApi.RequireValidClient(http.HandlerFunc(characterApi.AddCampaignActivityForCharacter))

		updateAdveture := userApi.RequireValidClient(http.HandlerFunc(adventureRecordApi.UpdateAdventure))

		getAdventure := allowCorsHeaders(http.HandlerFunc(adventureRecordApi.GetAdventure))
		getCharactersForCampaign := allowCorsHeaders(http.HandlerFunc(characterApi.GetCharactersForCampaign))

		//Campaign Endpoints
		router.Handle("POST /campaigns", createCampaign)
		router.Handle("POST /campaigns/{campaignId}/adventures", addAdventureToCampaign)
		router.Handle("POST /campaigns/{campaignId}/characters", addCharacterToCampaign)

		router.Handle("PATCH /campaigns/{campaignId}", updateCampaign)

		router.HandleFunc("GET /campaigns", campaignApi.ListCampaigns)
		router.HandleFunc("GET /campaigns/{campaignId}", campaignApi.GetCampaign)
		router.HandleFunc("GET /campaigns/{campaignId}/adventures", adventureRecordApi.ListAdventureRecordsForCampaign)
		router.Handle("GET /campaigns/{campaignId}/characters", getCharactersForCampaign)

		router.Handle("DELETE /campaigns/{campaignId}", deleteCampaign)

		//Adventure Endpoints
		//router.HandleFunc("POST/adventures/:adventureId/characters/:characterId", characterApi.ManageCharactersForAdventure)

		router.Handle("PATCH /adventures/{adventureId}", updateAdveture)

		router.Handle("GET /adventures/{adventureId}", getAdventure)

		//Character Endpoints
		router.Handle("POST /characters/{characterId}/campaign-actions", addCampaignActionToCharacter)

		router.HandleFunc(" GET /characters/{characterId}", characterApi.GetCharacterById)

		// USER
		router.HandleFunc(" GET /user/validate", userApi.ValidateClient)

		server := &http.Server{
			Addr:    ":9090",
			Handler: router,
		}
		logger.Print("Listening on 9090")
		for {
			server.ListenAndServe()
			logger.Print("Server crash... attempting restart")
		}

	}

}

/*
func preflight(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader()
	w.Header("Access-Control-Allow-Origin", "*")
	w.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers, content-type")
	w.JSON(http.StatusOK, struct{}{})
}*/

func allowCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}
