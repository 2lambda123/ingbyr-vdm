/*
 @Author: ingbyr
*/

package task

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/ingbyr/vdm/app/media"
	"github.com/ingbyr/vdm/pkg/store"
)

// DTask is a media downloading task
type DTask struct {
	*store.Model

	// Media is selected media format info
	*media.Info `json:"media" gorm:"embedded"`

	// FormatId is media format id from media.Format
	FormatId string `json:"formatId" form:"formatId"`

	// Engine is one of download engines
	Engine string `json:"engine" gorm:"engine" form:"engine"`

	// ExtArgs store the download engine command args from user input
	ExtArgs string `json:"extArgs" gorm:"ext_args"`

	// StoragePath is a media storage path
	StoragePath string `json:"storagePath" gorm:"storage_path" form:"storagePath"`

	// Progress will be updated in downloading operation and be sent to websocket client
	Progress Progress `json:"progress" gorm:"embedded"`

	Ctx    context.Context    `json:"-" gorm:"-"`
	Cancel context.CancelFunc `json:"-" gorm:"-"`
}

type Progress struct {
	// ID is from DTask id field
	ID        snowflake.ID `json:"id" gorm:"-"`
	Status    status       `json:"status" gorm:"status"`
	CmdOutput string       `json:"cmdOutput" gorm:"cmd_output"`
	Percent   string       `json:"progress" gorm:"percent"`
	Speed     string       `json:"speed" gorm:"speed"`
}

func NewDTask() *DTask {
	model := store.NewModel()
	return &DTask{
		Model: model,
		Progress: Progress{
			ID: model.ID,
		},
	}
}

func (dtask *DTask) Save() {
	store.DB.Save(dtask)
}

func (dtask *DTask) Find(page *store.Page) *store.Page {
	page.Data = &[]DTask{}
	tx := store.DB.Model(dtask)
	if dtask.Title != "" {
		tx.Where("title LIKE ?", "%"+dtask.Title+"%")
		dtask.Title = ""
	}
	if dtask.Desc != "" {
		tx.Where("desc LIKE ?", "%"+dtask.Desc+"%")
		dtask.Desc = ""
	}

	tx.Where(dtask).Order("status DESC")
	return store.PagingQuery(tx, page)
}

func (dtask *DTask) SaveProgress() {
	dtaskUpdater := DTask{Progress: dtask.Progress}
	store.DB.Model(dtask).Updates(dtaskUpdater)
}

func (dtask *DTask) SameTasks(page *store.Page) *store.Page {
	query := DTask{
		Info:        dtask.Info,
		FormatId:    dtask.FormatId,
		Engine:      dtask.Engine,
		ExtArgs:     dtask.ExtArgs,
		StoragePath: dtask.StoragePath,
	}
	return query.Find(page)
}
