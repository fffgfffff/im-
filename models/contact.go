package models

import (
	"ginChat/utils"
	"gorm.io/gorm"
)

// Contact 好友和群都存在这个表里面
// 可根据具体业务做拆分
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对应的谁
	Type     int  //1私聊 2群聊
	Desc     string
}

const (
	Contact_user  = 0x01
	Contact_group = 0x02
)

func (c *Contact) TableName() string {
	return "contact"
}
func CreateContacts(userId, targetId uint, typ int) error {
	return utils.DB.Create(&Contact{OwnerId: userId, TargetId: targetId, Type: typ}).Error
}
func FindContactsByOwnerIdAndTargetIdAndType(userId, targetId, typ int) Contact {
	contact := Contact{}
	utils.DB.Where("owner_id=? and target_id=? and type=?",
		userId, targetId, typ).First(&contact)
	return contact
}

// SearchFriend 找朋友
func SearchFriend(userId int) []UserBasic {
	contacts := make([]Contact, 0)
	ub := make([]UserBasic, 0)
	utils.DB.Where("owner_id=? and type=?", userId, Contact_user).Find(&contacts)
	for _, contact := range contacts {
		ub1 := UserBasic{}
		utils.DB.Where("id=?", contact.TargetId).First(&ub1)
		ub = append(ub, ub1)
	}
	return ub
}

// SearchCommunity 查找type=2群聊
func SearchCommunity(userId int64) []Contact {
	conconts := make([]Contact, 0)
	utils.DB.Where("owner_id=? and type=?", userId, Contact_group).Find(&conconts)
	return conconts
}
func SearchUserByGroupId(targetId int64) []Contact {
	conconts := make([]Contact, 0)
	utils.DB.Where("target_id=? and type=?", targetId, Contact_group).Find(&conconts)
	return conconts
}
