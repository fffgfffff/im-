package models

import (
	"fmt"
	"ginChat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name       string `json:"name"`
	Password   string `json:"password"`
	Phone      string `valid:"matches(^1[3-9]{1}\\d{9}$)" json:"phone"`
	Email      string `valid:"email" json:"email"`
	Identity   string `json:"identity"`
	ClientIp   string `json:"clientIp"`
	ClientPort string `json:"clientPort"`
	//随机数
	Salt          string    `json:"salt"`
	LoginTime     time.Time `json:"loginTime"`
	HeartbeatTime time.Time `json:"heartbeatTime"`
	LoginOutTime  time.Time `json:"loginOutTime"`
	IsLogout      bool      `json:"isLogout"`
	DeviceInfo    string    `json:"deviceInfo"`
	Memo          string    `json:"memo"` //头像
	Avatar        string    `json:"avatar"`
}

func (t *UserBasic) TableName() string {
	return "user_basic"
}
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)

	return data
}
func CreateUser(u UserBasic) {
	utils.DB.Create(&u)
}
func UpdateUser(u UserBasic) {
	utils.DB.Model(&u).Updates(&UserBasic{Name: u.Name, Password: u.Password})
}
func DelUser(u UserBasic) {
	utils.DB.Delete(&u)
}
func FindUserById(id int) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("id=?", id).First(user)
	return user
}
func FindUserByName(name string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name=?", name).First(user)
	return user
}
func FindUserByPhone(phone string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("phone=?", phone).First(user)
	return user
}
func FindUserByEmail(email string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("email=?", email).First(user)
	return user
}
func FindUserByPhoneAndPwd(phone, password string) *UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone=? and password=?", phone, password).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	//安全起见，实时刷新token
	utils.DB.Model(&user).Where("id=?", user.ID).Update("identity", temp)
	return &user
}
func CheckToken(userId int, token string) bool {
	ub := FindUserById(userId)
	return ub.Identity == token
}
