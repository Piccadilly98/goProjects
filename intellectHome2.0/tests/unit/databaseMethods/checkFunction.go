package data_base_methods

import database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/dataBase"

func CheckExistBoardInfo(db *database.DataBase, boardID string) (bool, error) {
	var exist bool

	err := db.GetPointer().QueryRow(`
		SELECT EXISTS(SELECT 1 FROM boardInfo
		WHERE board_id = $1)
	`, boardID).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}
