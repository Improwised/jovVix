package cli

import (
	"fmt"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetOrCreateDefaults(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	var force bool = false

	defaultCmd := cobra.Command{
		Use:   "create",
		Short: "To create default",
		Long:  `This command create default`,
		Args:  cobra.MinimumNArgs(1),
	}

	admin := cobra.Command{
		Use:   "admin <username> <email>",
		Short: "`create admin <username> <email>` create a new admin, you need to pass two argument username and email in excet order",
		Long: `This command will create admin with given username and email, 
	- If username or email exists then it will generate an error
	- But you can override username coolition by passing -f flag.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			adminInfo := struct {
				userName string
				email    string
			}{args[0], args[1]}

			if len(adminInfo.userName) > 12 || len(adminInfo.userName) < 6 {
				return fmt.Errorf("The length of the username must in betweeen 6 to 12 chars.")
			}

			db, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}

			isValid, err := utils.ValidateGlobalEmail(adminInfo.email)
			if err != nil {
				return err
			}

			if !isValid {
				return fmt.Errorf("%s is not a valid email", adminInfo.email)
			}

			userModel, err := models.InitUserModel(db)
			if err != nil {
				return err
			}

			userSvc := services.NewUserService(&userModel)

			passwordStr, err := password.Generate(14, 3, 3, false, false)
			if err != nil {
				return err
			}

			user := models.User{
				FirstName: "default",
				LastName:  "admin",
				Email:     adminInfo.email,
				Username:  adminInfo.userName,
				Password:  passwordStr,
				Roles:     "admin",
			}

			_, err = userSvc.RegisterUser(user, events.NewEventBus(logger))

			if err != nil {

				// Check is there unique-key error
				if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {

					if !(force && pqErr.Constraint == constants.UserUkey) {
						return fmt.Errorf("Unique key violation on %v", pqErr.Constraint)
					}

					// force flag is set then it tries agin with manipulated username tring to make new admin-username

					user.Username = GenerateNewUserName(user.Username, 5, 12)

					_, err = userSvc.RegisterUser(user, events.NewEventBus(logger))

					if err != nil {
						return fmt.Errorf("SomeError during register admin with new username %s", user.Username)
					}
				}

			}

			fmt.Printf("\nCreated new admin with... \nemail: %s username %s\npassword: %s \n- Runtime generated password!!! \n- Please store and change it as early as possible\n", adminInfo.email, adminInfo.userName, passwordStr)

			return nil
		},
	}

	admin.Flags().BoolVarP(&force, "force", "f", false, "Forcefully create new admin (may have security risk associated!)")

	defaultCmd.AddCommand(&admin)
	// Migration commands up, down

	return defaultCmd
}

func GenerateNewUserName(currentUserName string, randomStringLen int, maxLength int) string {
	random_str := "_" + utils.GenerateRandomString(randomStringLen-1)

	truncate_at := len(currentUserName)
	if truncate_at+len(random_str) > maxLength {
		truncate_at = maxLength - len(random_str)
	}
	return currentUserName[:truncate_at] + random_str
}
