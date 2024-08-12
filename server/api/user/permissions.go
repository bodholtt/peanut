package user

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"peanutserver/pcfg"
	"peanutserver/types"
	"reflect"
	"strconv"
)

func CheckPermissions(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	var rank int

	user, err := database.GetUser(id)
	if err != nil {
		rank = 0
	} else {
		rank = user.Rank
	}

	// compare user rank to all permissions and respond with a { permission: bool } map where 0 = denied 1 = allowed

	userPerms := &pcfg.Permissions{}
	value := reflect.ValueOf(*pcfg.Perms)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if int64(rank) >= field.Int() {
			newField := reflect.New(value.Field(i).Type())
			newField.Elem().SetInt(1)
			filteredPermissionsValue := reflect.ValueOf(userPerms).Elem()
			filteredPermissionsValue.FieldByName(value.Type().Field(i).Name).Set(newField.Elem())
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  userPerms,
		Error: "",
	})
}
