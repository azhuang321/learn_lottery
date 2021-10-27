package models

type LtUserday struct {
	Id         string `xorm:"not null pk autoincr UNSIGNED INT(10)"`
	Uid        string `xorm:"not null default 0 comment('用户ID') unique(uid_day) UNSIGNED INT(10)"`
	Day        string `xorm:"not null default 0 comment('日期，如：20180725') unique(uid_day) UNSIGNED INT(10)"`
	Num        string `xorm:"not null default 0 comment('次数') UNSIGNED INT(10)"`
	SysCreated string `xorm:"not null default 0 comment('创建时间') UNSIGNED INT(10)"`
	SysUpdated string `xorm:"not null default 0 comment('修改时间') UNSIGNED INT(10)"`
}
