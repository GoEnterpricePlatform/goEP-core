package model

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

type OptionNoSqlModel struct {
	Name        string `bson:"name"`
	VarOptName  string `bson:"var_opt_name"`
	VarOptValue string `bson:"var_opt_value"`
}

func (m *OptionNoSqlModel) ToDomain(o *domain.Option) {
	if m == nil {
		return
	}

	o.Name = m.Name
	o.VarOptName = m.VarOptName
	o.VarOptValue = m.VarOptValue
}

func FromDomainOption(o *domain.Option) *OptionNoSqlModel {
	if o == nil {
		return nil
	}

	return &OptionNoSqlModel{
		Name:        o.Name,
		VarOptName:  o.VarOptName,
		VarOptValue: o.VarOptValue,
	}
}
