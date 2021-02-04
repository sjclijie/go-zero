package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sjclijie/go-zero/core/stores/cache"
	"github.com/sjclijie/go-zero/core/stores/sqlc"
	"github.com/sjclijie/go-zero/core/stores/sqlx"
	"github.com/sjclijie/go-zero/core/stringx"
	"github.com/sjclijie/go-zero/tools/goctl/model/sql/builderx"
)

var (
	userFieldNames          = builderx.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheUserIdPrefix       = "cache#User#id#"
	cacheUserUserNamePrefix = "cache#User#userName#"
)

type (
	UserModel interface {
		Insert(data User) (sql.Result, error)
		FindOne(id int64) (*User, error)
		FindOneByUserName(userName string) (*User, error)
		Update(data User) error
		Delete(id int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id            int64          `db:"id"`              // 主键
		Uid           int64          `db:"uid"`             // 用户uid
		Email         sql.NullString `db:"email"`           // 邮箱
		DepartId      int64          `db:"depart_id"`       // 所属部门id(最小部门)
		UpdatedAt     time.Time      `db:"updated_at"`      // 更新时间
		CreatedAt     time.Time      `db:"created_at"`      // 申请时间
		UserName      string         `db:"user_name"`       // 真实姓名（中文名）
		DepartName    string         `db:"depart_name"`     // 部门名称
		CnName        string         `db:"cn_name"`         // 中文名
		Status        int64          `db:"status"`          // 用户状态 1-启用，2-禁用
		CityIds       string         `db:"city_ids"`        // 城市id
		LastLoginTime sql.NullTime   `db:"last_login_time"` // 最后一次登录时间
		UidBeisen     string         `db:"uid_beisen"`      // 北森id
		Gender        int64          `db:"gender"`          // 性别：0-男，1-女，2-保密 ----- 从这开始增加字段
		State         int64          `db:"state"`           // ucenter帐号状态:1-已启用(在职)，2-已停用(离职)
		Birthday      sql.NullTime   `db:"birthday"`        // 生日
		HllId         string         `db:"hll_id"`          // 工号
		JobLevel      string         `db:"job_level"`       // 职级
		JobName       string         `db:"job_name"`        // 职务名称
		职务名称

		DepartmentPath string `db:"department_path"` // 部门组织层级

		LevelId   string       `db:"level_id"`   // 点隔开部门id，比如：职能id.技术部id.大数据部id
		CustomId  string       `db:"custom_id"`  // 自定义的部门id，用逗号隔开（需要用英文字母区别lever_id)
		EntryDate sql.NullTime `db:"entry_date"` // 入职时间
		DeletedAt sql.NullTime `db:"deleted_at"` // 逻辑删除/删除时间
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user`",
	}
}

func (m *defaultUserModel) Insert(data User) (sql.Result, error) {
	userUserNameKey := fmt.Sprintf("%s%v", cacheUserUserNamePrefix, data.UserName)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
		return conn.Exec(query, data.Uid, data.Email, data.DepartId, data.UpdatedAt, data.CreatedAt, data.UserName, data.DepartName, data.CnName, data.Status, data.CityIds, data.LastLoginTime, data.UidBeisen, data.Gender, data.State, data.Birthday, data.HllId, data.JobLevel, data.JobName, data.DepartmentPath, data.LevelId, data.CustomId, data.EntryDate, data.DeletedAt)
	}, userUserNameKey)
	return ret, err
}

func (m *defaultUserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOneByUserName(userName string) (*User, error) {
	userUserNameKey := fmt.Sprintf("%s%v", cacheUserUserNamePrefix, userName)
	var resp User
	err := m.QueryRowIndex(&resp, userUserNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_name` = ? limit 1", userRows, m.table)
		if err := conn.QueryRow(&resp, query, userName); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Update(data User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		return conn.Exec(query, data.Uid, data.Email, data.DepartId, data.UpdatedAt, data.CreatedAt, data.UserName, data.DepartName, data.CnName, data.Status, data.CityIds, data.LastLoginTime, data.UidBeisen, data.Gender, data.State, data.Birthday, data.HllId, data.JobLevel, data.JobName, data.DepartmentPath, data.LevelId, data.CustomId, data.EntryDate, data.DeletedAt, data.Id)
	}, userIdKey)
	return err
}

func (m *defaultUserModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	userUserNameKey := fmt.Sprintf("%s%v", cacheUserUserNamePrefix, data.UserName)
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userUserNameKey, userIdKey)
	return err
}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRow(v, query, primary)
}
