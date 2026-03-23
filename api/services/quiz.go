package services

import (
	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type QuizService struct {
	quizModel     *models.QuizModel
	questionModel *models.QuestionModel
	db            *goqu.Database
	logger        *zap.Logger
}

func NewQuizService(db *goqu.Database, logger *zap.Logger) *QuizService {
	quizModel := models.InitQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)
	return &QuizService{
		quizModel:     quizModel,
		questionModel: questionModel,
		db:            db,
		logger:        logger,
	}
}

// This function will delete quiz only if no active quiz is present
func (quizSvc *QuizService) DeleteQuizById(quizId string) error {
	isOk := false
	transaction, err := quizSvc.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				quizSvc.logger.Error("error during commit in delete quiz", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				quizSvc.logger.Error("error during rollback in delete quiz", zap.Error(err))
			}
		}
	}()

	err = quizSvc.quizModel.DeleteQuizById(transaction, quizId)
	if err != nil {
		quizSvc.logger.Debug("error in DeleteQuizFromQuizQuestionById", zap.Error(err))
		return err
	}

	isOk = true

	return nil
}

// This function will delete question
func (quizSvc *QuizService) DeleteQuestionById(questionId string) error {
	isOk := false
	transaction, err := quizSvc.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				quizSvc.logger.Error("error during commit in delete quiz", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				quizSvc.logger.Error("error during rollback in delete quiz", zap.Error(err))
			}
		}
	}()

	// Update previous question's next_question pointer (column)
	err = quizSvc.questionModel.UpdatePreviousQuestionById(transaction, questionId)
	if err != nil {
		quizSvc.logger.Debug("error in DeleteQuizFromQuizQuestionById", zap.Error(err))
		return err
	}

	// Delete the question
	err = quizSvc.questionModel.DeleteQuestionById(transaction, questionId)
	if err != nil {
		quizSvc.logger.Debug("error in DeleteQuizFromQuizQuestionById", zap.Error(err))
		return err
	}

	isOk = true

	return nil
}

// Edit question by creating a new question row and rewiring quiz_questions to the new id.
// This preserves historical sessions and reports that still point to the old question id.
func (quizSvc *QuizService) EditQuestionById(quizId, oldQuestionId string, question models.Question) (string, error) {
	isOk := false
	transaction, err := quizSvc.db.Begin()
	if err != nil {
		return "", err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				quizSvc.logger.Error("error during commit in edit question", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				quizSvc.logger.Error("error during rollback in edit question", zap.Error(err))
			}
		}
	}()

	newQuestionId, err := quizSvc.questionModel.CreateQuestion(transaction, question)
	if err != nil {
		return "", err
	}

	err = quizSvc.questionModel.RewireQuizQuestionForEdit(transaction, quizId, oldQuestionId, newQuestionId)
	if err != nil {
		return "", err
	}

	isOk = true
	return newQuestionId.String(), nil
}
