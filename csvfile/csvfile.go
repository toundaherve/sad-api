package csvfile

import (
	"encoding/csv"
	"os"

	"github.com/toundaherve/sad-api/user"
)

type CSVFile struct {
	Filename string
}

func New(filename string) *CSVFile {
	return &CSVFile{
		Filename: filename,
	}
}

func (c *CSVFile) InitFile() error {
	f, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{"Name", "Email", "Country", "City", "Password"}); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (c *CSVFile) CreateUser(newUser *user.User) error {
	f, err := os.OpenFile(c.Filename, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.Write([]string{newUser.Name, newUser.Email, newUser.Country, newUser.City, newUser.Password})
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (c *CSVFile) GetUserByEmail(email string) (*user.User, error) {
	return nil, nil
}
