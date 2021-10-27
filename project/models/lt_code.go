package models

type LtCode struct {
	Id         string `xorm:"not null pk autoincr UNSIGNED INT(10)"`
	GiftId     string `xorm:"not null default 0 comment('奖品ID，关联lt_gift表') index UNSIGNED INT(10)"`
	Code       string `xorm:"not null default '' comment('虚拟券编码') unique VARCHAR(255)"`
	SysCreated string `xorm:"not null default 0 comment('创建时间') UNSIGNED INT(10)"`
	SysUpdated string `xorm:"not null default 0 comment('更新时间') UNSIGNED INT(10)"`
	SysStatus  string `xorm:"not null default 0 comment('状态，0正常，1作废，2已发放') UNSIGNED SMALLINT(5)"`
}
