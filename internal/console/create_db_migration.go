package console

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var createMigrationCmd = &cobra.Command{
	Use:   "create-migration [filename]",
	Short: "Create new migration file",
	Long:  "Create new migration file with the specified file name",
	Args:  cobra.ExactArgs(1),
	Run:   processCreateMigration,
}

func init() {
	RootCmd.AddCommand(createMigrationCmd)
}

func processCreateMigration(cmd *cobra.Command, args []string) {
	fileName := args[0]
	if err := checkMigrationFolderExists(); err != nil {
		fmt.Println("Error getting folder info:", err)
		return
	}
	migrationFile := []byte(`-- +migrate Up notransaction` + "\n\n" + `-- +migrate Down`)
	migrationFileName := "db/migration/" + createUniqueTime() + "_" + strings.ToLower(fileName) + ".sql"
	if err := ioutil.WriteFile(migrationFileName, migrationFile, 0666); err != nil {
		fmt.Println("Error creating file:", err)
	}
	fmt.Println(migrationFileName + " created")
}

func createUniqueTime() string {
	now := time.Now()
	splitDate := strings.Split(now.Format("01/02/2006"), "/")
	newDate := splitDate[2] + splitDate[0] + splitDate[1]
	hr, min, sc := now.Clock()
	hour := strconv.Itoa(hr)
	minute := strconv.Itoa(min)
	sec := strconv.Itoa(sc)
	if len(hour) == 1 {
		hour = "0" + hour
	}
	if len(minute) == 1 {
		minute = "0" + minute
	}
	if len(sec) == 1 {
		sec = "0" + sec
	}

	return newDate + hour + minute + sec
}

func checkMigrationFolderExists() error {
	_, err := os.Stat("db/migration/")
	if os.IsNotExist(err) { //check if db/migration folder is already exist
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("db/migration/ folder not found, want to create (Y/N)? ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return errors.New("failed when reading user input")
		}
		ans := strings.Contains(strings.ToUpper(input), "Y")
		if !ans {
			return errors.New("cancelled creating migration")
		}
		if err := createMigrationFolder(); err != nil {
			return err
		}
	}
	return nil
}

func createMigrationFolder() error {
	_, err := os.Stat("db/")
	if os.IsNotExist(err) { //check if db folder is already exist
		if err := os.Mkdir("db/", os.ModePerm); err != nil {
			return err
		}
	}
	if err := os.MkdirAll("db/migration/", os.ModePerm); err != nil {
		return err
	}
	return nil
}
