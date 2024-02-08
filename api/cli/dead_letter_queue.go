package cli

import (
	"go.uber.org/zap"

	"github.com/Improwised/quizz-app/api/cli/workers"
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/spf13/cobra"
)

type DeadLetterQ struct {
	Handler    string `json:"handler_poisoned"`
	Reason     string `json:"reason_poisoned"`
	Subscriber string `json:"subscriber_poisoned"`
	Topic      string `json:"topic_poisoned"`
}

// GetAPICommandDef runs app
func GetDeadQueueCommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {

	workerCommand := cobra.Command{
		Use:   "dead-letter-queue",
		Short: "To start dead-letter queue",
		Long:  `This queue is used to store failed job from all worker`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Init worker
			subscriber, err := watermill.InitSubscriber(cfg, true)
			if err != nil {
				return err
			}

			// run worker with topic(queue name) and process function
			// it will run failed job until it get success
			err = subscriber.Run(cfg.MQ.DeadLetterQ, "dead_letter_queue", workers.Process)
			return err
		},
	}
	return workerCommand
}
