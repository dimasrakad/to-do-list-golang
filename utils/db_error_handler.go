package utils

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"to-do-list-golang/models/dtos"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func HandleDBError(c *gin.Context, err error) {
	res := dtos.ErrorResponse{}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		res.Error = "Data not found"
		c.JSON(http.StatusNotFound, res)
		return
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062: // Duplicate entry
			re := regexp.MustCompile(`for key '(.+?)'`)
			matches := re.FindStringSubmatch(mysqlErr.Message)

			column := "data"
			if len(matches) > 1 {
				parts := strings.Split(matches[1], ".")
				column = parts[len(parts)-1]
			}

			res.Error = "Column " + column + " already exists, there cannot be duplicates"
			c.JSON(http.StatusBadRequest, res)
			return
		case 1452: // Cannot add or update a child row: foreign key constraint fails
			re := regexp.MustCompile(`FOREIGN KEY \(` + "`(.+?)`" + `\)`)
			matches := re.FindStringSubmatch(mysqlErr.Message)

			column := "foreign key"
			if len(matches) > 1 {
				column = matches[1]
			}

			res.Error = "Related data not found for column " + column + " (foreign key violation)"
			c.JSON(http.StatusBadRequest, res)
			return
		case 1048: // Column cannot be null
			re := regexp.MustCompile(`Column '(.+?)'`)
			matches := re.FindStringSubmatch(mysqlErr.Message)

			column := "column"
			if len(matches) > 1 {
				column = matches[1]
			}

			res.Error = "Required field " + column
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	res.Error = "An internal server error occurred\n" + err.Error()
	c.JSON(http.StatusInternalServerError, res)
}
