package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/floodedrealms/adventure-archivist/api"
	"github.com/floodedrealms/adventure-archivist/commands"
	"github.com/floodedrealms/adventure-archivist/internal/repository"
	"github.com/floodedrealms/adventure-archivist/internal/services"
	"github.com/floodedrealms/adventure-archivist/internal/util"
	"github.com/floodedrealms/adventure-archivist/webapp"
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

		//Turn on renderer for webpages (will panic if templates are wrong)
		renderer := webapp.NewRenderer()

		sqlRepo, err := repository.NewSqliteRepo("archivist.db", logger)
		util.CheckErr(err)

		//services
		campaignService := services.NewCampaignService(sqlRepo, logger, context.TODO())
		characterService := services.NewCharacterService(sqlRepo, logger, context.TODO())
		adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())
		userService := services.NewUserService(sqlRepo, *logger)
		// campaignActionService := services.NewCampaignActionService(sqlRepo)

		//api exposure
		/*campaignApi := api.NewCampaignApi(campaignService, characterService)
		adventureRecordApi := api.NewAdventureRecordApi(*adventureRecordService, characterService)
		characterApi := api.NewCharacterApi(characterService, *campaignActionService)*/
		userApi := api.NewClientAPI(userService)

		//pages
		homePages := webapp.NewHomePage(*renderer, *campaignService)
		campaignPages := webapp.NewCampaignPage(*campaignService, characterService, *renderer)
		adventurePages := webapp.NewAdventurePage(*adventureRecordService, characterService, *renderer)

		//router
		router := http.NewServeMux()

		// Wrap functions
		/*createCampaign := userApi.RequireValidClient(http.HandlerFunc(campaignApi.CreateCampaign))
		updateCampaign := userApi.RequireValidClient(http.HandlerFunc(campaignApi.UpdateCampaign))
		deleteCampaign := userApi.RequireValidClient(http.HandlerFunc(campaignApi.DeleteCampaign))
		addAdventureToCampaign := userApi.RequireValidClient(http.HandlerFunc(adventureRecordApi.CreateAdventureRecord))
		addCharacterToCampaign := userApi.RequireValidClient(http.HandlerFunc(characterApi.CreateCharacterForCampaign))
		addCampaignActionToCharacter := userApi.RequireValidClient(http.HandlerFunc(characterApi.AddCampaignActivityForCharacter))

		updateAdveture := userApi.RequireValidClient(http.HandlerFunc(adventureRecordApi.UpdateAdventure))

		getAdventure := allowCorsHeaders(http.HandlerFunc(adventureRecordApi.GetAdventure))
		getCharactersForCampaign := allowCorsHeaders(http.HandlerFunc(characterApi.GetCharactersForCampaign))*/

		//Campaign Endpoints
		/*router.Handle("POST /campaigns", createCampaign)
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

		router.HandleFunc(" GET /characters/{characterId}", characterApi.GetCharacterById)*/

		// USER
		router.HandleFunc(" GET /user/validate", userApi.ValidateClient)

		// static
		fs := http.FileServer(http.Dir("./static"))
		router.Handle("/static/", http.StripPrefix("/static/", fs))

		// Webapp Pages
		router.HandleFunc("/", homePages.Index)
		router.HandleFunc("/guild", homePages.GuildLanding)
		router.HandleFunc("/tavern", homePages.TavernLanding)
		router.HandleFunc("/crier", homePages.Campaigns)
		router.HandleFunc("/campaign-list", homePages.LoadNextCampaignSet)
		router.HandleFunc("/about", homePages.About)
		campaignPages.RegisterRoutes(router)
		adventurePages.RegisterRoutes(router)

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
