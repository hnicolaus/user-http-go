package repository

var (
	queryInsertUsers         = "INSERT INTO user(full_name, phone_number, password, created_time) VALUES"
	valuesInsertUsersF       = "($%d, $%d, $%d, $%d),"
	returnLastInsertedUserID = "RETURNING id"
)

var (
	querySelectUsers     = "SELECT id, full_name, phone_number, password, created_time, updated_time FROM user WHERE true"
	whereUserPhoneNumber = " AND phone_number = $%d"
	whereUserID          = " AND id = $%d"
)

var (
	queryIncrementSuccessfulLoginCount = "UPDATE user SET successful_login_count = successful_login_count + 1 WHERE id = $1"
)

var (
	queryUpdateUserF    = "UPDATE user SET %s WHERE TRUE"
	setUserPhoneNumberF = "phone_number = $%d"
	setUserFullNameF    = "full_name = $%d"
	setUserUpdatedTimeF = "updated_time = $%d"
)
