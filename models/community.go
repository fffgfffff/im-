package models

import (
	"ginChat/utils"
	"gorm.io/gorm"
	"strconv"
)

type Community struct {
	gorm.Model
	Name    string `json:"name"`
	OwnerId uint   `json:"ownerId"`
	Icon    string `json:"icon"` //群logo
	Desc    string `json:"desc"`
	Type    uint   `json:"type"`
}

func (c *Community) TableName() string {
	return "community"
}
func FindCommunityByStrId(comId string) Community {
	id, _ := strconv.Atoi(comId)
	community := Community{}
	utils.DB.Where("id=? or name=?", id, comId).First(&community)
	return community
}
func CreateCommunity(community Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	c := Community{}
	utils.DB.Where("owner_id=? and name=?", community.OwnerId, community.Name).First(&c)
	if c.ID > 0 {
		return -1, "群已存在"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		return -1, "建群失败"
	}
	if err := CreateContacts(community.OwnerId, community.ID, 2); err != nil {
		return 0, "创建群聊消息失败"
	}
	return 0, "建群成功"
}
func LoadCommunity(ownerId uint) []Community {
	contact := make([]Contact, 0)
	communitys := make([]Community, 0)
	utils.DB.Where("owner_id=? and type=? and target_id>0", ownerId, Contact_group).Find(&contact)
	for i := 0; i < len(contact); i++ {
		community := Community{}
		utils.DB.Where("id=?", contact[i].TargetId).First(&community)
		communitys = append(communitys, community)
	}
	return communitys
}
