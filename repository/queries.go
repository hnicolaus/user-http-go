package repository

var (
	queryInsertUsers         = "INSERT INTO sawitpro_user(full_name, phone_number, password, created_time) VALUES"
	valuesInsertUsersF       = "($%d, $%d, $%d, $%d),"
	returnLastInsertedUserID = "RETURNING id"
)
