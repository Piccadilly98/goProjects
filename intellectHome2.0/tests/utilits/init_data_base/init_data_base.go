package init_data_base

import (
	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/storage/dataBase"
)

func InitDataBase(db *database.DataBase) error {
	_, err := db.GetPointer().Exec(`
		INSERT INTO boards(board_id)
		VALUES
		('esp32_1_test'),
		('esp32_2_test'),
		('esp32_3_test');
	`)
	if err != nil {
		return err
	}
	_, err = db.GetPointer().Exec(`
	UPDATE boards
	SET controllers = jsonb_set(
	controllers,
	'{devices,binary}',
 	COALESCE(controllers->'devices'->'binary', '[]') ||
	jsonb_build_object(
	'name', null, 
	'type', 'non-type', 
	'status', false, 
	'pin_number', 1, 
	'created_date', NOW(), 
	'updated_date', null, 
	'controller_id', 'led1'))
	WHERE board_id = 'esp32_2_test';
	`)
	if err != nil {
		return err
	}

	_, err = db.GetPointer().Exec(`
	UPDATE boards
	SET controllers = jsonb_set(
	controllers,
	'{devices,binary}',
 	COALESCE(controllers->'devices'->'binary', '[]') ||
	jsonb_build_object(
	'name', null, 
	'type', 'non-type', 
	'status', false, 
	'pin_number', 1, 
	'created_date', NOW(), 
	'updated_date', null, 
	'controller_id', 'led2'))
	WHERE board_id = 'esp32_1_test';
	`)
	if err != nil {
		return err
	}
	_, err = db.GetPointer().Exec(`
	UPDATE boards
	SET controllers = jsonb_set(
	controllers,
	'{devices,binary}',
 	COALESCE(controllers->'devices'->'binary', '[]') ||
	jsonb_build_object(
	'name', null, 
	'type', 'non-type', 
	'status', false, 
	'pin_number', 1, 
	'created_date', NOW(), 
	'updated_date', null, 
	'controller_id', 'led3'))
	WHERE board_id = 'esp32_3_test';
	`)
	if err != nil {
		return err
	}

	_, err = db.GetPointer().Exec(`
	UPDATE boards
	SET controllers = jsonb_set(
	controllers,
	'{devices,sensor}',
 	COALESCE(controllers->'devices'->'sensor', '[]') ||
	jsonb_build_object(
	'name', null,                                      
	'type', 'non-type',                                
	'unit', '%',                                       
	'value', 0,                                        
	'pin_number', null,                                
	'created_date', NOW(),
	'updated_date', null,                              
	'controller_id', 'led4'))
	WHERE board_id = 'esp32_3_test';
	`)

	if err != nil {
		return err
	}
	return nil
}

func Cleanup(db *database.DataBase) error {
	_, err := db.GetPointer().Exec(`
	DELETE FROM boards;`)
	return err
}
