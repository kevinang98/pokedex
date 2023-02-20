package models

const (
	UserKey    = "pokemon"
	JwtKey     = "pokemon-dex"
	PokemonDB  = "pokemon"
	PokemonCol = "pokedex"

	ErrorNoAccess          = "you have no access"
	ErrorUsernameNotExist  = "username not exist"
	ErrorInvalidPassword   = "invalid password"
	ErrorUserPassEmpty     = "username or password is empty"
	ErrorUserPassTooLong   = "username or password is too long"
	ErrorInvalidRole       = "role must be admin or user"
	ErrorFailGetBody       = "fail to get body"
	ErrorFailMarshalBody   = "fail to marshal body"
	ErrorInvalidBodyReq    = "invalid body request"
	ErrorPIDEmpty          = "pid is empty"
	ErrorPIDZero           = "pid must be greater than 0"
	ErrorInvalidQueryParam = "invalid query param"
	ErrorOffsetLimitParam  = "offset param must be 0 or greater, limit must be greater than 0"
	ErrorSortParam         = "sort param must be name or pid or empty"
	ErrorOrderParam        = "order param must be asc or desc or empty"
)
