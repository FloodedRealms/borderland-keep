package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/floodedrealms/borderland-keep/archivist"
	"github.com/floodedrealms/borderland-keep/commands"
	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
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
		debug := false
		if len(flags) == 2 {
			if flag.Arg(1) == "true" {
				debug = true
			}
		}

		logger := util.NewLogger(debug)

		//Turn on renderer for webpages (will panic if templates are wrong)
		renderer := archivist.NewRenderer()

		sqlRepo, err := repository.NewSqliteRepo("archivist.db", logger)
		util.CheckErr(err)

		//services
		campaignService := services.NewCampaignService(sqlRepo, logger, context.TODO())
		characterService := services.NewCharacterService(sqlRepo, logger, context.TODO())
		adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())
		// userService := services.NewUserService(sqlRepo, *logger)
		// campaignActionService := services.NewCampaignActionService(sqlRepo)

		//pages
		homePages := archivist.NewHomePage(*renderer, *campaignService)
		campaignPages := archivist.NewCampaignPage(*campaignService, *characterService, *adventureRecordService, *renderer)
		adventurePages := archivist.NewAdventurePage(*adventureRecordService, *characterService, *renderer)

		//router
		router := http.NewServeMux()

		// USER
		// router.HandleFunc(" GET /user/validate", userApi.ValidateClient)

		// static
		fs := http.FileServer(http.Dir("./static"))
		router.Handle("/static/", http.StripPrefix("/static/", fs))

		// archivist Pages
		router.HandleFunc("/", homePages.Index)
		router.HandleFunc("/guild", homePages.GuildLanding)
		router.HandleFunc("/tavern", homePages.TavernLanding)
		router.HandleFunc("/archivist", homePages.Campaigns)
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
