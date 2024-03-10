package service

import (
	"errors"
	"ginChat/models"
	"ginChat/utils"
)

type ContactService struct {
}

func (service *ContactService) SearchFriends(userId int) []models.UserBasic {
	return models.SearchFriend(userId)
}
func (service *ContactService) AddFriend(userId, targetId int) error {
	if userId == targetId {
		return errors.New("不能添加自己为好友")
	}
	if ub := models.FindUserById(targetId); ub.ID < 1 || ub.Identity == "" {
		return errors.New("该用户不存在")
	}
	//查询是否已经是好友
	contact := models.FindContactsByOwnerIdAndTargetIdAndType(userId, targetId, models.Contact_user)
	//如果存在记录说明已经是好友了不加
	if contact.ID > 0 {
		return errors.New("该用户已经被添加过")
	}
	//启动事务，插入两条数据
	db := utils.DB.Begin()
	err1 := models.CreateContacts(uint(userId), uint(targetId), models.Contact_user)
	err2 := models.CreateContacts(uint(targetId), uint(userId), models.Contact_user)
	if err1 == nil && err2 == nil {
		//提交
		db.Commit()
		return nil
	} else {
		//回滚
		db.Rollback()
		if err1 != nil {
			return err1
		} else {
			return err2
		}
	}
}
func (service *ContactService) SearchCommunityIds(userId int64) []int64 {
	comIds := make([]int64, 0)
	//todo 获取用户全部群ID
	conconts := models.SearchCommunity(userId)
	for _, v := range conconts {
		comIds = append(comIds, int64(v.TargetId))
	}
	return comIds
}
func (service *ContactService) JoinCommunity(userId int, comId string) (int, string) {
	if comId == "" {
		return -1, "未输入"
	}
	community := models.FindCommunityByStrId(comId)
	if community.ID < 0 {
		return -1, "没有找到群"
	}
	contact := models.FindContactsByOwnerIdAndTargetIdAndType(
		userId, int(community.ID), models.Contact_group)
	if contact.ID > 0 {
		return -1, "群已加过"
	}

	if err := models.CreateContacts(uint(userId), community.ID,
		models.Contact_group); err != nil {
		return -1, "加群失败"
	}
	return int(community.ID), "加群成功"
}
