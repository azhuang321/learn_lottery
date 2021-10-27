package models

type LtBlackip struct {
	Id         string `xorm:"not null pk autoincr UNSIGNED INT(10)"`
	Ip         string `xorm:"not null default '' comment('IP地址') unique VARCHAR(50)"`
	Blacktime  string `xorm:"not null default 0 comment('黑名单限制到期时间') UNSIGNED INT(10)"`
	SysCreated string `xorm:"not null default 0 comment('创建时间') UNSIGNED INT(10)"`
	SysUpdated string `xorm:"not null default 0 comment('修改时间') UNSIGNED INT(10)"`
}
