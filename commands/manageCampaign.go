package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type config struct {
	Host string `json:"host"`
	Id   string `json:"id"`
	Key  string `json:"key"`
}

const contentHeader string = "Content-Type"
const jsonContentType string = "application/json"
const apiClientIdHeader string = "X-Archivist-Client-Id"
const apiSecretKeyHeader string = "X-Archivist-API-Key"

var currentConfig config = config{}

var adventureOps = map[string]string{"A": "[A]dd", "E": "[E]dit"}
var campaignOps = map[string]string{"C": "[C]ampiagn Details", "H": "C[h]aracters", "A": "[A]dventures"}
var adventureCharacterOps = map[string]string{"C": "[C]haracter", "H": "[H]enchmen"}

func loadConfig() {
	dirname, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dirname, "archivist", "config.json")
	log.Print(path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatal("Configuration file does not exist")
	}
	err = util.DecodeJSONFileAbsolutePath(path, &currentConfig)
	if err != nil {
		log.Fatal(err)
	}
}

type apiAdventureResponse struct {
	Status string                `json:"status"`
	Data   types.AdventureRecord `json:"data"`
}

type campaignList struct {
	Status    string                 `json:"status"`
	Campaigns []types.CampaignRecord `json:"data"`
}

func getCampaigns() campaignList {
	url := currentConfig.Host + "/campaigns"
	resp, err := http.Get(url)
	util.CheckErr(err)
	defer resp.Body.Close()

	campaigns := campaignList{}
	util.DecodeJSON(resp.Body, &campaigns)
	return campaigns
}

type characterList struct {
	Status     string                  `json:"status"`
	Characters []types.CharacterRecord `json:"data"`
}

func getCharactersForCampaign(id int) characterList {
	url := currentConfig.Host + fmt.Sprintf("/campaigns/%d/characters", id)
	resp, err := http.Get(url)
	util.CheckErr(err)
	defer resp.Body.Close()

	characters := characterList{}
	util.DecodeJSON(resp.Body, &characters)
	return characters

}

func pickCampaign(cl campaignList) int {
	fmt.Println("Please enter the number of the campaign you wish to managed:")
	pickedCampaign := retrieveNumberInRangeFromStdIn(1, len(cl.Campaigns)+1) // rande looks good to users
	pickedCampaign--                                                         // prevent off by 1 error
	fmt.Printf("You picked number %d : %s\n", pickedCampaign+1, cl.Campaigns[pickedCampaign].Name)
	return pickedCampaign

}

func manageAdventures(campaignId int) {
	op := displayAndRetrieveOps(adventureOps)
	switch op {
	case "E":
		fmt.Println("This functionality is not supported")
	case "A":
		addAdventureToCampaign(campaignId)
	}

}

func addAdventureToCampaign(campaignId int) {

	adventureReader := bufio.NewReader(os.Stdin)

	fmt.Println("What is the Adventure Name?")
	name, _ := adventureReader.ReadString('\n')
	name = strings.Trim(name, "\n")
	createdAdventure := types.CreateAdventureRequest{
		Name:          name,
		CampaignID:    campaignId,
		AdventureDate: time.Now(),
	}
	adventureId := makeAdventureCreationRequest(campaignId, createdAdventure)
	updateRequest := &types.UpdateAdventureRequest{
		ID:         adventureId,
		Gems:       make([]types.Loot, 0),
		Jewellery:  make([]types.Loot, 0),
		MagicItems: make([]types.IncomingMagicItem, 0),
		Combat:     make([]types.Loot, 0),
		Characters: make([]types.UpdateAdventureCharacter, 0),
	}

	fmt.Println("How much Copper was recovered? Whole Numbers Only.")
	copper := retreiveNumber()
	fmt.Println("How much Silver was recovered? Whole Numbers Only.")
	silver := retreiveNumber()
	fmt.Println("How much Electrum was recovered? Whole Numbers Only.")
	electrum := retreiveNumber()
	fmt.Println("How much Gold was recovered? Whole Numbers Only.")
	gold := retreiveNumber()
	fmt.Println("How much Platinum was recovered? Whole Numbers Only.")
	platinum := retreiveNumber()
	updateRequest.Copper = copper
	updateRequest.Silver = silver
	updateRequest.Electrum = electrum
	updateRequest.Gold = gold
	updateRequest.Platinum = platinum
	fmt.Println("Would you like to add Gems to the adventure? (Y/N)")
	g := retrieveYesNoFrontStdIn()
	if g == "Y" {
		updateRequest.Gems = createLootArray("Gem")
	}
	fmt.Println("Would you like to add Jewellery to the adventure? (Y/N)")
	g = retrieveYesNoFrontStdIn()
	if g == "Y" {
		updateRequest.Jewellery = createLootArray("Jewellery")
	}
	fmt.Println("Would you like to add Combat to the adventure? (Y/N)")
	g = retrieveYesNoFrontStdIn()
	if g == "Y" {
		updateRequest.Combat = createLootArray("Monster")
	}
	fmt.Println("Would you like to add Magic Items to the adventure? (Y/N)")
	g = retrieveYesNoFrontStdIn()
	if g == "Y" {
		updateRequest.MagicItems = createMagicItemArray()
	}
	newAdventure := makeAdventureUpdateRequest(*updateRequest)
	fullShareXP, halfShare := newAdventure.CalculateXPShares()
	characters := getCharactersForCampaign(campaignId)
	fmt.Println("Which characters would you like to add to the adventure?")
	fmt.Println("1-3 for range, 1,3 for separated list.")
	selectedCharacters := retrieveComboNumbersInRange(1, len(characters.Characters))
	for _, selected := range selectedCharacters {
		op := displayAndRetrieveAdventureCharacterOps(characters.Characters[selected-1].Name)
		c := types.UpdateAdventureCharacter{ID: characters.Characters[selected-1].Id(), Halfshare: false, XpGained: fullShareXP}
		switch op {
		case "H":
			c.Halfshare = true
			c.XpGained = halfShare
			updateRequest.Characters = append(updateRequest.Characters, c)
		case "C":
			c.Halfshare = false
			c.XpGained = fullShareXP
			updateRequest.Characters = append(updateRequest.Characters, c)

		}
	}
	makeAdventureUpdateRequest(*updateRequest)
	fmt.Println("Adventure created")
}

func makeAdventureCreationRequest(c int, a types.CreateAdventureRequest) int {
	url := currentConfig.Host + fmt.Sprintf("/campaigns/%d/adventures", c)
	data, err := json.Marshal(a)
	if err != nil {
		log.Fatalf("Failed to marshal Data.\nError: %e", err)
	}
	payload := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	req.Close = true
	if err != nil {
		log.Fatalf("Failed to create request.\nError: %e", err)
	}
	req.Header.Set(apiClientIdHeader, currentConfig.Id)
	req.Header.Set(apiSecretKeyHeader, currentConfig.Key)
	req.Header.Set(contentHeader, jsonContentType)
	log.Print(req.Header.Values(currentConfig.Key))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request.\nError: %s", err.Error())
	}
	defer resp.Body.Close()
	log.Print(resp)
	newA := apiAdventureResponse{}
	err = util.DecodeJSON(resp.Body, &newA)
	if err != nil {
		log.Fatalf("Failed to decode response.\nError: %s", err.Error())
	}
	return newA.Data.ID
}
func makeAdventureUpdateRequest(a types.UpdateAdventureRequest) types.AdventureRecord {
	url := currentConfig.Host + fmt.Sprintf("/adventures/%d", a.ID)
	data, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err.Error())
	}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Failed to maked Adventure Update Request\nError: %s", err.Error())
	}
	addHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to complete Adventure Update Request\nError: %s", err.Error())
	}
	defer resp.Body.Close()
	updateRequest := apiAdventureResponse{}
	err = util.DecodeJSON(resp.Body, &updateRequest)
	if err != nil {
		log.Fatalf("Failed to decode Adventure after Update Request\nError: %s\nBody: %s", err.Error())
	}
	return updateRequest.Data
}

func createLootArray(lType string) []types.Loot {
	loot := make([]types.Loot, 0)
	finished := false
	for !finished {
		fmt.Printf("Name of %s\n", lType)
		name := retreiveString()
		fmt.Println("Xp Value")
		value := retreiveNumber()
		fmt.Println("Number Recovered")
		amount := retreiveNumber()
		loot = append(loot, types.Loot{Name: name, XPValueOfOne: float64(value), NumberOfItem: amount})
		fmt.Printf("Would you like to create another %s?\n", lType)
		a := retrieveYesNoFrontStdIn()
		if a == "N" {
			finished = true
		}
	}
	return loot
}

func createMagicItemArray() []types.IncomingMagicItem {
	loot := make([]types.IncomingMagicItem, 0)
	finished := false
	for !finished {
		fmt.Println("Name of Magic Item")
		name := retreiveString()
		fmt.Println("Xp Value")
		value := retreiveNumber()
		fmt.Println("Actual Value")
		amount := retreiveNumber()
		loot = append(loot, types.IncomingMagicItem{Name: name, ApparentValue: value, ActualValue: amount})
		fmt.Println("Would you like to create another Magic Item (Y/N)")
		a := retrieveYesNoFrontStdIn()
		if a == "N" {
			finished = true
		}
	}
	return loot
}

func manageCharacters(campaignId int) {
	fmt.Println("This Functionality is not yet supported")
}
func editCampaignDetails(campaignId int) {
	fmt.Println("This Functioanality is not yet supported")
}

func ManageCampaign() {
	loadConfig()
	campaignList := getCampaigns()
	campaignString := fmt.Sprintf("Which Campaign do you wish to update [1-%d]:", len(campaignList.Campaigns))
	fmt.Println(campaignString)
	fmt.Println("----|campaign name| database id|---")
	for i, campaign := range campaignList.Campaigns {
		outString := fmt.Sprintf("%d----| %s | %d", i+1, campaign.Name, campaign.ID)
		fmt.Println(outString)
	}
	pickedCampaign := pickCampaign(campaignList)
	cId := campaignList.Campaigns[pickedCampaign].ID
	operation := displayAndRetrieveOps(campaignOps)
	switch operation {
	case "A":
		manageAdventures(cId)
	case "C":
		editCampaignDetails(cId)
	case "H":
		manageCharacters(cId)
	}
}

func retreiveString() string {
	var s string
	strReader := bufio.NewReader(os.Stdin)
	input, err := strReader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	s = strings.Trim(input, "\n")
	return s
}

func retreiveNumber() int {
	numberPicked := false
	var num int
	for !numberPicked {
		numReader := bufio.NewReader(os.Stdin)
		picked, _ := numReader.ReadString('\n')
		n, err := strconv.Atoi(strings.Trim(picked, "\n"))
		if err != nil {
			log.Print(err.Error())
			fmt.Println("Please enter a number")
			continue
		}
		num = n
		numberPicked = true
	}
	return num
}

func retrieveNumberInRangeFromStdIn(s, e int) int {
	num := -1
	for num < s || num > e {
		numReader := bufio.NewReader(os.Stdin)
		picked, _ := numReader.ReadString('\n')
		id, err := strconv.Atoi(strings.Trim(picked, "\n"))
		if err != nil {
			log.Print(err.Error())
			fmt.Println("Please enter a number")
			continue
		}
		num = id
	}
	return num
}

func retrieveComboNumbersInRange(s, e int) []int {
	num := make([]int, 0)
	for len(num) < 1 {
		numReader := bufio.NewReader(os.Stdin)
		input, _ := numReader.ReadString('\n')
		input = strings.Trim(input, "\n")
		blocks := strings.Split(input, ",")
		for _, block := range blocks {
			numbers := strings.Split(block, "-")
			if len(numbers) == 1 {
				number, err := inputIsNumericAndInRange(numbers[0], s, e)
				if err != nil {
					break
				}
				num = append(num, number)
			} else {
				lowNumber, lErr := inputIsNumericAndInRange(numbers[0], s, e)
				if lErr != nil {
					break
				}
				highNumber, hErr := inputIsNumericAndInRange(numbers[1], s, e)
				if hErr != nil {
					break
				}
				num = append(num, lowNumber)
				num = append(num, highNumber)
			}
		}
	}
	return num
}

func retrieveCharacterInSetFromStdIn(validChar map[string]string) string {
	char := ""
	for char == "" {
		charReader := bufio.NewReader(os.Stdin)
		picked, _ := charReader.ReadString('\n')
		picked = strings.Trim(picked, "\n")
		picked = strings.ToUpper(picked) // cover cases
		_, ok := validChar[picked]
		if ok {
			char = picked
			return char
		}
		fmt.Println("Please enter a valid choice")

	}
	return char
}

func retrieveYesNoFrontStdIn() string {
	set := make(map[string]string)
	set["Y"] = "Yes"
	set["N"] = "No"
	return retrieveCharacterInSetFromStdIn(set)

}

func displayAndRetrieveOps(ops map[string]string) string {
	opstring := "Would you like to: "
	for _, op := range ops {
		opstring += fmt.Sprintf("%s, ", op)
	}
	fmt.Println(opstring[0:strings.LastIndex(opstring, ",")])
	return retrieveCharacterInSetFromStdIn(ops)
}

func displayAndRetrieveAdventureCharacterOps(name string) string {
	opstring := fmt.Sprintf("Is % a %s or a %s: ", name, adventureCharacterOps["C"], adventureCharacterOps["H"])
	fmt.Println(opstring[0:strings.LastIndex(opstring, ",")])
	return retrieveCharacterInSetFromStdIn(adventureCharacterOps)
}

func inputIsNumericAndInRange(input string, s, e int) (int, error) {
	number, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Character %s not valid \n", input)
		return 0, errors.New("Invalid Character")
	}
	return numberIsInRange(number, s, e)
}

func numberIsInRange(number, s, e int) (int, error) {
	if number < s || number > e {
		fmt.Printf("Number %d not in range\n", number)
		return 0, errors.New("Number not in range")
	}
	return number, nil
}

func addHeaders(req *http.Request) {
	req.Header.Set(apiClientIdHeader, currentConfig.Id)
	req.Header.Set(apiSecretKeyHeader, currentConfig.Key)
	req.Header.Set(contentHeader, jsonContentType)
}
