# The Adventure Archivist API
This is the API project for the Adventure Archivist. The API itself will have a front end hosted on borderlandkeep.com, but is available for self hosting should people desire. 

# Local Development
If you want to develop the API locally, you will need to first create a user in the database. Simply run ```go run main.go create-user <name>``` and then copy the displayed client_id and api_key down to whatever you are using to make test requests.

# API Structure
The API is primarily concerned with simple data management. It will include functionality for calculating experience on characters, based on the assumptions made in the Adventure Conqueror King System core rules. Other systems will be up to front ends to implement themselves, if desired. The primary functions of the API, as well as their endpoints will be listed here. 

## Campaign Crud functionality
### POST /campaigns
Creates a new campaign
### POST /campaings/{campaignId}/adventures
Creates a new adventure associated with the given camapaignId
### GET /campaigns
Lists all campaigns in the database 
### GET /campaigns/{camapaignId}
Retrieves a CAMPAIGN object
### DELETE /campaigns/campaignId
Deletes a campaign from the database, including all adventures and characters

## Adventure management functionality

### GET /adventures/{campaignId}
Lists all adventures in the database associated with the given campaignId 
### GET /adventures/{adventureId}
Retrieves an ADVENTURE RECORD object
### POST /adventures/{adventureId}/loot/{type}
Adds a new loot record to the adventure. The possible loot record types are:
- Coins
- Gems
- Jewellery
- Magic Items
### PATCH /adventures/{adventureId}/{lootId}/[{attribute,value}]
Updates a loot record for the adventure.
### POST /adventures/{adventureId}/combat
Adds a new combat record to the adventure
### POST /adventures/{adventureId}/{combatId}
Updates a new combat record for the adventure
### GET /adventures/{adventureId}/exeperience
Returns the total amount of experience earned for a given adventure
### POST /adventures/{adventureId}/{characterId}/{operation}
Adds or removes a character from an adventure
### DELETE /adventures/{adventureId}
Deletes an adventure

## Character management functionality
### POST /characters/{campaignId}
Adds a new character associated with the camapaignId
### GET /characters/{campaignId}
Adds a new character associated with the camapaignId
### GET /characters/{characterId}
Retrieves a CHARACTER RECORD object
### PATCH /characters/{characterId}/[{attribute,value}]
Updates the given list of attributes to the new values
### DELETE /characters/{characterId}
Deletes a character

