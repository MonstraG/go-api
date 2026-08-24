package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/models"
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"net/http"
	"time"

	"gorm.io/gorm"
)

var indexTemplate = pages.ParsePage(
	"nav.gohtml",
	"admin/adminPage.gohtml",
)

type Controller struct {
	db           *gorm.DB
	vpsLoginLink string
}

type PageData struct {
	pages.PageData
	VpsLoginLink string
	SqlResult    MySqlResult
}

type MySqlResult struct {
	Shown bool

	Result       sql.Result
	RowsAffected int64
	Headers      []string
	Rows         [][]string
	Err          error
}

func NewController(config appConfig.AppConfig, db *gorm.DB) *Controller {
	return &Controller{db: db, vpsLoginLink: config.VpsLoginLink}
}

func (controller *Controller) GetAdminPage(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	controller.renderAdminPage(w, r, MySqlResult{})
}

func (controller *Controller) renderAdminPage(w reqRes.MyResponseWriter, r *reqRes.MyRequest, mySqlResult MySqlResult) {
	var pageData = PageData{
		PageData:     pages.NewPageData(r, "Homepage"),
		VpsLoginLink: controller.vpsLoginLink,
		SqlResult:    mySqlResult,
	}

	w.RenderTemplate(indexTemplate, pageData)
}

func (controller *Controller) SetPasswordChangeStatus(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	ok := r.ParseFormRequired(w)
	if !ok {
		return
	}
	username := r.GetFormFieldRequired(w, "username")
	if username == "" {
		return
	}
	canChangePassword := r.GetFormFieldRequired(w, "canChangePassword")
	if canChangePassword == "" {
		return
	}

	user := models.User{Username: username}
	result := controller.db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("Failed to find user %v", result.Error)
		w.Error(message, http.StatusBadRequest)
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Failed to find user %v", result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	user.CanResetPassword = canChangePassword == "true"
	result = controller.db.Save(&user)
	if result.Error != nil {
		message := fmt.Sprintf("Failed to update user %v", result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (controller *Controller) PostSql(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	ok := r.ParseFormRequired(w)
	if !ok {
		w.Error("Empty request", http.StatusBadRequest)
		return
	}

	statement := r.GetFormFieldRequired(w, "sql")
	if statement == "" {
		w.Error("SQL field empty", http.StatusBadRequest)
		return
	}

	result := controller.executeRaw(statement)

	controller.renderAdminPage(w, r, result)
}

func (controller *Controller) executeRaw(query string) MySqlResult {
	result := gorm.WithResult()

	sqlRows, err := gorm.G[any](controller.db, result).Raw(query).Rows(context.Background())
	if err != nil {
		return MySqlResult{
			Shown:        true,
			RowsAffected: 0,
			Headers:      make([]string, 0),
			Rows:         make([][]string, 0),
			Err:          err,
		}
	}
	defer helpers.CloseSafely(sqlRows)

	// column headers
	headers, err := sqlRows.Columns()
	if err != nil {
		return MySqlResult{
			Shown:        true,
			Result:       result.Result,
			RowsAffected: result.RowsAffected,
			Headers:      headers,
			Rows:         make([][]string, 0),
			Err:          err,
		}
	}

	rows := make([][]string, 0)

	scanTargets := make([]interface{}, len(headers))
	raw := make([]interface{}, len(headers))

	for sqlRows.Next() {
		for i := range raw {
			scanTargets[i] = &raw[i]
		}

		err = sqlRows.Scan(scanTargets...)
		if err != nil {
			return MySqlResult{
				Shown:        true,
				Result:       result.Result,
				RowsAffected: result.RowsAffected,
				Headers:      headers,
				Rows:         rows,
				Err:          err,
			}
		}

		rowStrings := make([]string, len(headers))
		for i, value := range raw {
			switch castValue := value.(type) {
			case nil:
				rowStrings[i] = "null"
			case []byte:
				rowStrings[i] = string(castValue)
			case time.Time:
				rowStrings[i] = castValue.Format(time.RFC3339)
			default:
				rowStrings[i] = fmt.Sprintf("%v", castValue)
			}
		}

		rows = append(rows, rowStrings)
	}

	err = sqlRows.Err()
	if err != nil {
		return MySqlResult{
			Shown:        true,
			Result:       result.Result,
			RowsAffected: result.RowsAffected,
			Headers:      headers,
			Rows:         rows,
			Err:          err,
		}
	}

	return MySqlResult{
		Shown:        true,
		Result:       result.Result,
		RowsAffected: result.RowsAffected,
		Headers:      headers,
		Rows:         rows,
		Err:          err,
	}
}
