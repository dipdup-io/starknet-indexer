package storage

import (
	"context"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
	"github.com/uptrace/bun"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock -typed
type IClassReplace interface {
	storage.Table[*ClassReplace]

	ByHeight(ctx context.Context, height uint64) ([]ClassReplace, error)
}

// ClassReplace -
type ClassReplace struct {
	bun.BaseModel `bun:"class_replace" comment:"Table with class replace history"`

	ID          uint64 `bun:"id,pk,notnull,autoincrement" comment:"Unique internal identity"`
	ContractId  uint64 `bun:"contract_id"                 comment:"Contract id"`
	PrevClassId uint64 `bun:"prev_class_id"               comment:"Previous class id"`
	NextClassId uint64 `bun:"next_class_id"               comment:"Next class id"`
	Height      uint64 `bun:"height"                      comment:"Block height"`

	Contract  Address `bun:"rel:belongs-to" hasura:"table:address,field:contract_id,remote_field:id,type:oto,name:contract"`
	PrevClass Class   `bun:"rel:belongs-to,join:prev_class_id=id" hasura:"table:class,field:prev_class_id,remote_field:id,type:oto,name:prev_class"`
	NextClass Class   `bun:"rel:belongs-to,join:next_class_id=id" hasura:"table:class,field:next_class_id,remote_field:id,type:oto,name:next_class"`
}

// TableName -
func (ClassReplace) TableName() string {
	return "class_replace"
}
