package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Improwised/jovvix/api/config"
	"github.com/spf13/cobra"
)

// Identity represents a Kratos identity.
type Identity struct {
	ID     string `json:"id"`
	Traits struct {
		Email string `json:"email"`
	} `json:"traits"`
}

func GetDeleteOrphanedCommand(cfg config.AppConfig) cobra.Command {
	deleteCommand := cobra.Command{
		Use:   "delete-orphans",
		Short: "To delete the orphan users from kratos.",
		Long:  `To delete the orphan users from kratos, which are already deleted from backend.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Kratos.AdminUrl == "" {
				return fmt.Errorf("KRATOS_ADMIN_URL and ORY_ACCESS_TOKEN must be set as environment variables")
			}

			// Emails to delete
			emailsToDelete := []string{
				"test@Improwised.com",
				"singh@gmail.com",
				"riya123@gmail.com",
				"rishi123@gmail.com",
				"richa123@gmail.com",
				"khushal.mer@improwised.com",
				"husen123@gmail.com",
				"shivani@improwised.com",
				"shivani123456@gmail.com",
				"husen.kureshi@improwised.com",
				"husen.kureshi@improwised1.com",
				"husenkureshi2003@gmail.com",
				"ctkinavar@gmail.com",
				"ankit.jilka@improwised.com",
				"ankitjilka10@gmail.com",
				"angita.shah@improwised.com",
				"19comp.ankit.jilka@gmail.com",
				"19ce152.ankit.jilka@vvpedulink.ac.in",
				"ashvintwst@gmail.com",
				"ashvintest45@gmail.com",
				"ashvin.bambhaniya+test@improwised.com",
				"ashvinbambhaniyatest@gmail.com",
				"ashvin.bambhaniya@improwised.com",
				"shaktirajsinh.zala@improwised.com",
			}

			// Create a map for quick email lookup
			emailsMap := make(map[string]bool)
			for _, email := range emailsToDelete {
				emailsMap[email] = true
			}

			// 1. List all identities
			client := &http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/identities", cfg.Kratos.AdminUrl), nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to list identities: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("failed to read response body: %v", err)
				}
				return fmt.Errorf("failed to list identities: status code %d, body: %s", resp.StatusCode, string(body))
			}

			var identities []Identity
			if err := json.NewDecoder(resp.Body).Decode(&identities); err != nil {
				return fmt.Errorf("failed to decode identities: %v", err)
			}

			fmt.Printf("Found %d total identities. Filtering for users to delete...\n", len(identities))
			isEmailExists := false

			// 2. Filter and delete identities
			for _, identity := range identities {
				if emailsMap[identity.Traits.Email] {
					isEmailExists = true
					fmt.Printf("Attempting to delete user: %s (ID: %s)\n", identity.Traits.Email, identity.ID)

					deleteReq, err := http.NewRequest("DELETE", fmt.Sprintf("%s/identities/%s", cfg.Kratos.AdminUrl, identity.ID), nil)
					if err != nil {
						log.Printf("Failed to create delete request for %s: %v", identity.Traits.Email, err)
						continue
					}

					deleteResp, err := client.Do(deleteReq)
					if err != nil {
						log.Printf("Failed to delete identity for %s: %v", identity.Traits.Email, err)
						continue
					}
					defer deleteResp.Body.Close()

					switch deleteResp.StatusCode {
					case http.StatusNoContent:
						fmt.Printf("✅ Successfully deleted user: %s\n", identity.Traits.Email)
					case http.StatusNotFound:
						fmt.Printf("⚠️ User already deleted (not found): %s\n", identity.Traits.Email)
					default:
						body, _ := io.ReadAll(deleteResp.Body)
						log.Printf("Failed to delete user %s: status code %d, body: %s", identity.Traits.Email, deleteResp.StatusCode, string(body))
					}
				} else {
					fmt.Printf("No action needed for user: %s (ID: %s)\n", identity.Traits.Email, identity.ID)
				}
			}

			if isEmailExists {
				fmt.Println("🎉 Cleanup script finished.")
			} else {
				fmt.Println("⚠️  No specified users found to delete.")
			}

			return nil
		},
	}

	return deleteCommand
}
