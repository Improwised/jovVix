package cli

import (
	"fmt"
	"strings"

	"github.com/Improwised/quizz-app/api/config"
	quizController "github.com/Improwised/quizz-app/api/controllers/api/v1"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetOrCreateDefaults(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	defaultCmd := cobra.Command{
		Use:   "create",
		Short: "To create entity from cmd",
		Long:  `This command created to have short-cut to create entity and reduce UI iteration.`,
		Args:  cobra.MinimumNArgs(1),
	}

	adminCmd := AdminCmd(cfg, logger)

	defaultCmd.AddCommand(adminCmd)

	return defaultCmd
}

func AdminCmd(cfg config.AppConfig, logger *zap.Logger) *cobra.Command {

	var force bool = false

	admin := cobra.Command{
		Use:   "admin <username> <email> <first-name> <last-name>",
		Short: "`create admin <username> <email> <first-name> <last-name>` create a new admin, you need to pass two argument username and email in exact order",
		Long: `This command will create admin with given username, email, first-name and last-name,
	- If username or email exists then it will generate an error
	- But you can override username collision by passing -f flag.
	- Remember, -f will modify username which changes some last characters of the username to make username unique`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			adminObj := models.User{
				Username:  args[0],
				Email:     args[1],
				FirstName: args[2],
				LastName:  args[3],
				Roles:     "admin",
			}

			if len(adminObj.Username) > 12 || len(adminObj.Username) < 6 {
				return fmt.Errorf("the length of the username must in between 6 to 12 chars")
			}

			db, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}

			isValid, err := utils.ValidateGlobalEmail(adminObj.Email)
			if err != nil {
				return err
			}

			if !isValid {
				return fmt.Errorf("%s is not a valid email", adminObj.Email)
			}

			adminObj.Password, err = password.Generate(14, 3, 3, false, false)
			if err != nil {
				return err
			}

			adminObj, err = quizController.CreateQuickUser(db, logger, adminObj, force, true)

			if err != nil {
				return err
			}

			message := fmt.Sprintf(`
				Created new admin with...
				email: %s
				username: %s
				password: %s
				first name: %s
				last name: %s
				- Runtime generated password!!!
				- Please store and change it as early as possible!!!`,
				adminObj.Email,
				adminObj.Username,
				adminObj.Password,
				adminObj.FirstName,
				adminObj.LastName)

			fmt.Println(strings.ReplaceAll(message, "\t", ""))

			return nil
		},
	}

	admin.Flags().BoolVarP(&force, "force", "f", false, "Forcefully create new admin (may have security risk associated!)")

	return &admin
}
