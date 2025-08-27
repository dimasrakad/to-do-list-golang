package utils

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func HandleDBError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
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

			c.JSON(http.StatusBadRequest, gin.H{"error": "Column " + column + " already exists, there cannot be duplicates"})
			return
		case 1452: // Cannot add or update a child row: foreign key constraint fails
			re := regexp.MustCompile(`FOREIGN KEY \(` + "`(.+?)`" + `\)`)
			matches := re.FindStringSubmatch(mysqlErr.Message)

			column := "foreign key"
			if len(matches) > 1 {
				column = matches[1]
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": "Related data not found for column " + column + " (foreign key violation)"})
			return
		case 1048: // Column cannot be null
			re := regexp.MustCompile(`Column '(.+?)'`)
			matches := re.FindStringSubmatch(mysqlErr.Message)

			column := "column"
			if len(matches) > 1 {
				column = matches[1]
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": "Required field " + column})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred\n" + err.Error()})
}
