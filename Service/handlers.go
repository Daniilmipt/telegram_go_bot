package Service

import (
	"bot/Domain"
	"bot/Repository/DataBase"
	"bot/Repository/Interface"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const (
	helpCmd   = "help"
	listCmd   = "list"
	addCmd    = "add"
	changeCmd = "change"
	deleteCmd = "delete"
)

var Service *UserService

func New(repository Interface.Repository) *UserService {
	Service = &UserService{
		repository: repository,
	}
	return Service
}

type UserService struct {
	repository Interface.Repository
}

//---------

func listFunc(string) string {
	data := DataBase.Map.List()
	res := make([]string, 0, len(data))

	for _, value := range data {
		res = append(res, value.GetData())
	}

	return strings.Join(res, "\n")
}

func helpFunc(string) string {
	return "help - list commands\n" +
		"list - list data\n" +
		"add <name> <age> - add new person \n" +
		"change <id> <name> <age> - change info about person\n" +
		"delete <id> - delete person"
}

// addFunc handles the "add" command
// It takes a single string of space-separated values
// The first value is the name of the new user
// The second value is their age
// It returns a string with the result of the operation
// If there is an error, it will be in the form of an error message
func addFunc(data string) string {
	// Split the data into separate values
	params := strings.Split(data, " ")
	// Check that we have the expected number of parameters
	if len(params) != 2 {
		// If not, return an error with the number of params and the values we got
		return errors.Wrapf(errors.New("bad argument"), "%d items: <%v>", len(params), params).Error()
	}

	// Try to parse the second parameter (age) as an integer
	age, err := strconv.Atoi(params[1])
	if err != nil {
		return err.Error()
	}

	// Try to create a new user entity with the given name and age
	user, err := Domain.NewEntity(params[0], uint(age))
	if err != nil {
		return err.Error()
	}

	// Try to add the new user to the database
	err = DataBase.Map.Add(user)
	if err != nil {
		return err.Error()
	}

	// If everything went well, return a success message with some info about the new user
	return fmt.Sprintf("user (%d) %s:%d added", user.Id, user.Name, user.Age)
}


// changeFunc handles the "change" command
// It takes a single string of space-separated values
// The first value is the ID of the user to be updated
// The second value is the new name of the user
// The third value is the new age of the user
// It returns a string with the result of the operation
// If there is an error, it will be in the form of an error message
func changeFunc(data string) string {
	// Split the data into separate values
	params := strings.Split(data, " ")
	// Check that we have the expected number of parameters
	if len(params) != 3 {
		// If not, return an error with the number of params and the values we got
		return errors.Wrapf(errors.New("bad argument"), "%d items: <%v>", len(params), params).Error()
	}

	// Try to parse the third parameter (age) as an integer
	age, _ := strconv.Atoi(params[2])

	// Try to create a new user entity with the given name and age
	user := Domain.Entity{}

	// Try to set the name of the user
	if err := user.SetName(params[1]); err != nil {
		fmt.Println(err)
		return ""
	}

	// Try to set the age of the user
	if err := user.SetAge(uint(age)); err != nil {
		fmt.Println(err)
		return ""
	}

	// Try to parse the first parameter (ID) as an integer
	id, _ := strconv.Atoi(params[0])

	// Try to update the user with the given ID with the new info
	err := DataBase.Map.Update(&user, uint(id))
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("user info was changed")
}

// deleteFunc is the function that handles the "delete" command
// It takes a single string argument which is the ID of the user to be deleted
// It returns a string with the result of the operation
// If there is an error, it will be in the form of an error message
func deleteFunc(data string) string {
	// Split the input into separate values
	params := strings.Split(data, "\n")
	// Check that we have exactly one value (the user ID)
	if len(params) != 1 {
		return errors.Wrapf(errors.New("bad argument"), "%d items: <%v>", len(params), params).Error()
	}

	// Try to parse the input (the user ID) as an integer
	id, _ := strconv.Atoi(params[0])

	// Try to delete the user with the given ID from the database
	err := DataBase.Map.Delete(uint(id))
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("user deleted")
}


func AddHandlers(bot *Bot) {
	bot.RegisterHandler(helpCmd, helpFunc)
	bot.RegisterHandler(listCmd, listFunc)
	bot.RegisterHandler(addCmd, addFunc)
	bot.RegisterHandler(changeCmd, changeFunc)
	bot.RegisterHandler(deleteCmd, deleteFunc)
}
