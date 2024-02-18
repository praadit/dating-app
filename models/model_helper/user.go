package modelhelper

import (
	"github.com/praadit/dating-apps/models"
)

func CreateUserBens(data map[string]interface{}) models.UserBenefits {
	ben := models.UserBenefits{}
	if val, ok := data["base_swipe"]; ok {
		if temp, tempok := val.(float64); tempok {
			ben.BaseSwipe = int(temp)
		} else {
			ben.BaseSwipe = val.(int)
		}
	}
	if val, ok := data["is_premium"]; ok {
		ben.IsPremium = val.(bool)
	}

	return ben
}
