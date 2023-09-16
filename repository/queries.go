package repository

var (
	queryInsertUsers         = "INSERT INTO sawitpro_user(full_name, phone_number, password, created_time) VALUES"
	valuesInsertUsersF       = "($%d, $%d, $%d, $%d),"
	returnLastInsertedUserID = "RETURNING id"
)

var (
	querySelectUsers     = "SELECT id, full_name, phone_number, password, created_time, updated_time FROM sawitpro_user WHERE true"
	whereUserPhoneNumber = " AND phone_number = $%d"
	whereUserID          = " AND id = $%d"
)

var (
	queryIncrementSuccessfulLoginCount = "UPDATE sawitpro_user SET successful_login_count = successful_login_count + 1 WHERE id = $1"
)
