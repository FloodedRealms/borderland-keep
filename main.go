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
	}
	operation := flag.Arg(0)
	switch operation {

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
	default:
		//	memRepo, err := repository.NewMemoryRepo()
		//	util.CheckErr(err)

		logger := util.NewLogger(true)

		sqlRepo, err := repository.NewSqliteRepo("archivist.db", logger)
		util.CheckErr(err)

		campaignService := services.NewCampaignService(sqlRepo, logger, context.TODO())
		characterService := services.NewCharacterService(sqlRepo, logger, context.TODO())
		adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())
		//userService := services.NewUserService(sqlRepo, *logger)

		campaignApi := api.NewCampaignApi(campaignService, characterService)
		adventureRecordApi := api.NewAdventureRecordApi(adventureRecordService)
		//characterApi := api.NewCharacterApi(characterService)
		//userApi := api.NewUserApi(userService)

		//characterService := services.NewCharacterService(memRepo, context.TODO())
		//characterApi := api.NewCharacterApi(characterService)

		router := http.NewServeMux()

		//Campaign Endpoints
		router.HandleFunc("POST /api/campaigns", campaignApi.CreateCampaign)
		router.HandleFunc("POST /api/campaigns/{campaignId}/adventures", adventureRecordApi.CreateAdventureRecord)
		//	router.HandleFunc("POST /api/campaigns/:campaignId/characters", characterApi.CreateCharacterForCampaign)

		router.HandleFunc("PATCH /api/campaigns/{campaignId}", campaignApi.UpdateCampaign)

		router.HandleFunc("GET /api/campaigns", campaignApi.ListCampaigns)
		router.HandleFunc("GET /api/campaigns/{campaignId}", campaignApi.GetCampaign)
		router.HandleFunc("GET /api/campaigns/{campaignId}/adventures", adventureRecordApi.ListAdventureRecordsForCampaign)

		router.HandleFunc("DELETE /api/campaigns/{campaignId}", campaignApi.DeleteCampaign)

		//Adventure Endpoints
		//		router.HandleFunc(" POST/api/adventures/:adventureId/characters/:characterId", characterApi.ManageCharactersForAdventure)

		router.HandleFunc("PATCH /api/adventures/{adventureId}", adventureRecordApi.UpdateAdventure)

		router.HandleFunc("GET /api/adventures/{adventureId}", adventureRecordApi.GetAdventure)

		//Character Endpoints
		//		router.HandleFunc(" GET /api/characters/:characterId", characterApi.GetCharacterById)

		// USER
		//		router.HandleFunc(" GET /api/user/validate", userApi.ValidateApiUser)

		server := &http.Server{
			Addr:    ":9090",
			Handler: router,
		}
		logger.Print("Listening on 9090")
		for true {
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
