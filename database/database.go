package database

const DataSavePath = "database/database.json"

//type Database struct {
//	UserID int
//	Name   string
//}
//
//type Users struct {
//	Users []Database
//}
//
//func DataSave(update *tgbotapi.Update) {
//	dataIn, err := os.ReadFile(assets.GamesSavePath)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	var users Users
//
//	err = json.Unmarshal(dataIn, &users)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	data := Database{
//		UserID: update.Message.From.ID,
//		Name:   update.Message.From.UserName,
//	}
//
//	users.Users = append(users.Users, data)
//
//	dataSave, err := json.MarshalIndent(&users, "", "  ")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	err = os.WriteFile(DataSavePath, dataSave, 0600)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
