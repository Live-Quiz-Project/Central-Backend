package v1

import (
	"context"
	"sort"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repo Repository) Service {
	return &service{
		Repository: repo,
		timeout:    time.Duration(3) * time.Second,
	}
}

func (s *service) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx, err := s.Repository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *service) CommitTransaction(ctx context.Context, tx *gorm.DB) error {
	err := s.Repository.CommitTransaction(tx)
	if err != nil {
		return err
	}
	return nil
}

// ---------- Quiz related service methods ---------- //
func (s *service) CreateQuiz(ctx context.Context, tx *gorm.DB, req *CreateQuizRequest, uid uuid.UUID) (*CreateQuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	q := &Quiz{
		ID:             req.ID,
		CreatorID:      uid,
		Title:          req.Title,
		Description:    req.Description,
		CoverImage:     req.CoverImage,
		Visibility:     req.Visibility,
		TimeLimit:      req.TimeLimit,
		HaveTimeFactor: req.HaveTimeFactor,
		TimeFactor:     req.TimeFactor,
		FontSize:       req.FontSize,
		Mark:           req.Mark,
		// SelectUpTo:     req.SelectUpTo,
		SelectMin:     req.SelectMin,
		SelectMax:     req.SelectMax,
		CaseSensitive: req.CaseSensitive,
	}

	qh := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         q.ID,
		CreatorID:      q.CreatorID,
		Title:          q.Title,
		Description:    q.Description,
		CoverImage:     q.CoverImage,
		Visibility:     q.Visibility,
		TimeLimit:      q.TimeLimit,
		HaveTimeFactor: q.HaveTimeFactor,
		TimeFactor:     q.TimeFactor,
		FontSize:       q.FontSize,
		Mark:           q.Mark,
		SelectMin:      q.SelectMin,
		SelectMax:      q.SelectMax,
		CaseSensitive:  q.CaseSensitive,
	}

	quiz, err := s.Repository.CreateQuiz(c, tx, q)
	if err != nil {
		return &CreateQuizResponse{}, err
	}
	quizH, er := s.Repository.CreateQuizHistory(c, tx, qh)
	if er != nil {
		return &CreateQuizResponse{}, er
	}

	return &CreateQuizResponse{
		QuizResponse: QuizResponse{
			Quiz: Quiz{
				ID:             quiz.ID,
				CreatorID:      quiz.CreatorID,
				Title:          quiz.Title,
				Description:    quiz.Description,
				CoverImage:     quiz.CoverImage,
				Visibility:     quiz.Visibility,
				TimeLimit:      quiz.TimeLimit,
				HaveTimeFactor: quiz.HaveTimeFactor,
				TimeFactor:     quiz.TimeFactor,
				FontSize:       quiz.FontSize,
				Mark:           quiz.Mark,
				SelectMin:      quiz.SelectMin,
				SelectMax:      quiz.SelectMax,
				CaseSensitive:  quiz.CaseSensitive,
				CreatedAt:      quiz.CreatedAt,
				UpdatedAt:      quiz.UpdatedAt,
				DeletedAt:      quiz.DeletedAt,
			},
		},
		QuizHistoryID: quizH.ID,
	}, nil
}

func (s *service) GetQuizzes(ctx context.Context, uid uuid.UUID) ([]QuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quizzes, err := s.Repository.GetQuizzesByUserID(c, uid)
	if err != nil {
		return nil, err
	}

	var res []QuizResponse
	for _, q := range quizzes {
		res = append(res, QuizResponse{
			Quiz: Quiz{
				ID:             q.ID,
				CreatorID:      q.CreatorID,
				Title:          q.Title,
				Description:    q.Description,
				CoverImage:     q.CoverImage,
				Visibility:     q.Visibility,
				TimeLimit:      q.TimeLimit,
				HaveTimeFactor: q.HaveTimeFactor,
				TimeFactor:     q.TimeFactor,
				FontSize:       q.FontSize,
				Mark:           q.Mark,
				SelectMin:      q.SelectMin,
				SelectMax:      q.SelectMax,
				CaseSensitive:  q.CaseSensitive,
				CreatedAt:      q.CreatedAt,
				UpdatedAt:      q.UpdatedAt,
				DeletedAt:      q.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetQuizByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetQuizByID(c, id)
	if err != nil {
		return nil, err
	}

	return &QuizResponse{
		Quiz: Quiz{
			ID:             quiz.ID,
			CreatorID:      quiz.CreatorID,
			Title:          quiz.Title,
			Description:    quiz.Description,
			CoverImage:     quiz.CoverImage,
			Visibility:     quiz.Visibility,
			TimeLimit:      quiz.TimeLimit,
			HaveTimeFactor: quiz.HaveTimeFactor,
			TimeFactor:     quiz.TimeFactor,
			FontSize:       quiz.FontSize,
			Mark:           quiz.Mark,
			SelectMin:      quiz.SelectMin,
			SelectMax:      quiz.SelectMax,
			CaseSensitive:  quiz.CaseSensitive,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			DeletedAt:      quiz.DeletedAt,
		},
	}, nil
}

func (s *service) GetDeleteQuizByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetDeleteQuizByID(c, id)
	if err != nil {
		return nil, err
	}

	return &QuizResponse{
		Quiz: Quiz{
			ID:             quiz.ID,
			CreatorID:      quiz.CreatorID,
			Title:          quiz.Title,
			Description:    quiz.Description,
			CoverImage:     quiz.CoverImage,
			Visibility:     quiz.Visibility,
			TimeLimit:      quiz.TimeLimit,
			HaveTimeFactor: quiz.HaveTimeFactor,
			TimeFactor:     quiz.TimeFactor,
			FontSize:       quiz.FontSize,
			Mark:           quiz.Mark,
			SelectMin:      quiz.SelectMin,
			SelectMax:      quiz.SelectMax,
			CaseSensitive:  quiz.CaseSensitive,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			DeletedAt:      quiz.DeletedAt,
		},
	}, nil
}

func (s *service) GetQuizHistories(ctx context.Context, uid uuid.UUID) ([]QuizHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quizHistories, err := s.Repository.GetQuizHistoriesByUserID(c, uid)
	if err != nil {
		return nil, err
	}

	var res []QuizHistoryResponse
	for _, q := range quizHistories {
		res = append(res, QuizHistoryResponse{
			QuizHistory: QuizHistory{
				ID:             q.ID,
				QuizID:         q.QuizID,
				CreatorID:      q.CreatorID,
				Title:          q.Title,
				Description:    q.Description,
				CoverImage:     q.CoverImage,
				Visibility:     q.Visibility,
				TimeLimit:      q.TimeLimit,
				HaveTimeFactor: q.HaveTimeFactor,
				TimeFactor:     q.TimeFactor,
				FontSize:       q.FontSize,
				Mark:           q.Mark,
				SelectMin:      q.SelectMin,
				SelectMax:      q.SelectMax,
				CaseSensitive:  q.CaseSensitive,
				CreatedAt:      q.CreatedAt,
				UpdatedAt:      q.UpdatedAt,
				DeletedAt:      q.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetQuizHistoryByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetQuizHistoryByID(c, id)
	if err != nil {
		return nil, err
	}

	return &QuizHistoryResponse{
		QuizHistory: QuizHistory{
			ID:             quiz.ID,
			QuizID:         quiz.QuizID,
			CreatorID:      quiz.CreatorID,
			Title:          quiz.Title,
			Description:    quiz.Description,
			CoverImage:     quiz.CoverImage,
			Visibility:     quiz.Visibility,
			TimeLimit:      quiz.TimeLimit,
			HaveTimeFactor: quiz.HaveTimeFactor,
			TimeFactor:     quiz.TimeFactor,
			FontSize:       quiz.FontSize,
			Mark:           quiz.Mark,
			SelectMin:      quiz.SelectMin,
			SelectMax:      quiz.SelectMax,
			CaseSensitive:  quiz.CaseSensitive,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			DeletedAt:      quiz.DeletedAt,
		},
	}, nil
}

func (s *service) UpdateQuiz(ctx context.Context, tx *gorm.DB, req *UpdateQuizRequest, uid uuid.UUID, id uuid.UUID) (*UpdateQuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetQuizByID(c, id)
	if err != nil {
		return &UpdateQuizResponse{}, err
	}
	if req.Title != "" {
		quiz.Title = req.Title
	}
	if req.Description != "" {
		quiz.Description = req.Description
	}
	if req.CoverImage != "" {
		quiz.CoverImage = req.CoverImage
	}
	if req.Title != "" {
		quiz.Title = req.Title
	}
	if req.Description != "" {
		quiz.Description = req.Description
	}
	if req.CoverImage != "" {
		quiz.CoverImage = req.CoverImage
	}
	if req.Visibility != "" {
		quiz.Visibility = req.Visibility
	}
	if req.TimeLimit != 0 {
		quiz.TimeLimit = req.TimeLimit
	}
	if req.HaveTimeFactor {
		quiz.HaveTimeFactor = req.HaveTimeFactor
	}
	if req.TimeFactor != 0 {
		quiz.TimeFactor = req.TimeFactor
	}
	if req.FontSize != 0 {
		quiz.FontSize = req.FontSize
	}
	if req.Mark != 0 {
		quiz.Mark = req.Mark
	}
	// if req.SelectMin != 0 {
	// 	quiz.SelectMin = req.SelectMin
	// }
	// if req.SelectMax != 0 {
	// 	quiz.SelectMax = req.SelectMax
	// }
	if !req.CaseSensitive {
		quiz.CaseSensitive = req.CaseSensitive
	}

	qh := &QuizHistory{
		ID:             uuid.New(),
		QuizID:         quiz.ID,
		CreatorID:      quiz.CreatorID,
		Title:          quiz.Title,
		Description:    quiz.Description,
		CoverImage:     quiz.CoverImage,
		Visibility:     quiz.Visibility,
		TimeLimit:      quiz.TimeLimit,
		HaveTimeFactor: quiz.HaveTimeFactor,
		TimeFactor:     quiz.TimeFactor,
		FontSize:       quiz.FontSize,
		Mark:           quiz.Mark,
		SelectMin:      quiz.SelectMin,
		SelectMax:      quiz.SelectMax,
		CaseSensitive:  quiz.CaseSensitive,
	}

	quiz, er := s.Repository.UpdateQuiz(c, tx, quiz)
	if er != nil {
		return &UpdateQuizResponse{}, er
	}

	_, e := s.Repository.CreateQuizHistory(c, tx, qh)
	if e != nil {
		return &UpdateQuizResponse{}, e
	}

	return &UpdateQuizResponse{
		QuizResponse: QuizResponse{
			Quiz: Quiz{
				ID:             quiz.ID,
				CreatorID:      quiz.CreatorID,
				Title:          quiz.Title,
				Description:    quiz.Description,
				CoverImage:     quiz.CoverImage,
				Visibility:     quiz.Visibility,
				TimeLimit:      quiz.TimeLimit,
				HaveTimeFactor: quiz.HaveTimeFactor,
				TimeFactor:     quiz.TimeFactor,
				FontSize:       quiz.FontSize,
				Mark:           quiz.Mark,
				SelectMin:      quiz.SelectMin,
				SelectMax:      quiz.SelectMax,
				CaseSensitive:  quiz.CaseSensitive,
				CreatedAt:      quiz.CreatedAt,
				UpdatedAt:      quiz.UpdatedAt,
				DeletedAt:      quiz.DeletedAt,
			},
		},
		QuizHistoryID: qh.ID,
	}, nil
}

func (s *service) DeleteQuiz(ctx context.Context, tx *gorm.DB, quizID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteQuiz(c, tx, quizID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreQuiz(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreQuiz(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

// ---------- Question Pool related service methods ---------- //
func (s *service) CreateQuestionPool(ctx context.Context, tx *gorm.DB, req *QuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID) (*CreateQuestionPoolResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	qp := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         quizID,
		Order:          req.Order,
		PoolOrder:      req.PoolOrder,
		Content:        req.Content,
		Note:           req.Note,
		Media:          req.Media,
		MediaType:      req.MediaType,
		TimeLimit:      req.TimeLimit,
		HaveTimeFactor: req.HaveTimeFactor,
		TimeFactor:     req.TimeFactor,
		FontSize:       req.FontSize,
	}

	qph := &QuestionPoolHistory{
		ID:             uuid.New(),
		QuestionPoolID: qp.ID,
		QuizID:         quizHistoryID,
		Order:          qp.Order,
		PoolOrder:      req.PoolOrder,
		Content:        qp.Content,
		Note:           qp.Note,
		Media:          qp.Media,
		MediaType:      qp.MediaType,
		TimeLimit:      qp.TimeLimit,
		HaveTimeFactor: qp.HaveTimeFactor,
		TimeFactor:     qp.TimeFactor,
		FontSize:       qp.FontSize,
	}

	questionPool, err := s.Repository.CreateQuestionPool(c, tx, qp)
	if err != nil {
		return &CreateQuestionPoolResponse{}, err
	}
	questionPoolH, er := s.Repository.CreateQuestionPoolHistory(c, tx, qph)
	if er != nil {
		return &CreateQuestionPoolResponse{}, er
	}

	return &CreateQuestionPoolResponse{
		QuestionPoolResponse: QuestionPoolResponse{
			QuestionPool: QuestionPool{
				ID:             questionPool.ID,
				QuizID:         questionPool.QuizID,
				Order:          questionPool.Order,
				PoolOrder:      questionPool.PoolOrder,
				Content:        questionPool.Content,
				Note:           questionPool.Note,
				Media:          questionPool.Media,
				MediaType:      questionPool.MediaType,
				TimeLimit:      questionPool.TimeLimit,
				HaveTimeFactor: questionPool.HaveTimeFactor,
				TimeFactor:     questionPool.TimeFactor,
				FontSize:       questionPool.FontSize,
				CreatedAt:      questionPool.CreatedAt,
				UpdatedAt:      questionPool.UpdatedAt,
				DeletedAt:      questionPool.DeletedAt,
			},
		},
		QuestionPoolHistoryID: questionPoolH.ID,
	}, nil
}

func (s *service) GetQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questionPools, err := s.Repository.GetQuestionPoolsByQuizID(c, quizID)
	if err != nil {
		return nil, err
	}

	var res []QuestionPoolResponse
	for _, qp := range questionPools {
		res = append(res, QuestionPoolResponse{
			QuestionPool: QuestionPool{
				ID:             qp.ID,
				QuizID:         qp.QuizID,
				Order:          qp.Order,
				PoolOrder:      qp.PoolOrder,
				Content:        qp.Content,
				Note:           qp.Note,
				Media:          qp.Media,
				MediaType:      qp.MediaType,
				TimeLimit:      qp.TimeLimit,
				HaveTimeFactor: qp.HaveTimeFactor,
				TimeFactor:     qp.TimeFactor,
				FontSize:       qp.FontSize,
				CreatedAt:      qp.CreatedAt,
				UpdatedAt:      qp.UpdatedAt,
				DeletedAt:      qp.DeletedAt,
			},
		})
	}
	return res, nil
}

func (s *service) GetDeleteQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questionPools, err := s.Repository.GetDeleteQuestionPoolsByQuizID(c, quizID)
	if err != nil {
		return nil, err
	}

	var res []QuestionPoolResponse
	for _, qp := range questionPools {
		res = append(res, QuestionPoolResponse{
			QuestionPool: QuestionPool{
				ID:             qp.ID,
				QuizID:         qp.QuizID,
				Order:          qp.Order,
				PoolOrder:      qp.PoolOrder,
				Content:        qp.Content,
				Note:           qp.Note,
				Media:          qp.Media,
				MediaType:      qp.MediaType,
				TimeLimit:      qp.TimeLimit,
				HaveTimeFactor: qp.HaveTimeFactor,
				TimeFactor:     qp.TimeFactor,
				FontSize:       qp.FontSize,
				CreatedAt:      qp.CreatedAt,
				UpdatedAt:      qp.UpdatedAt,
				DeletedAt:      qp.DeletedAt,
			},
		})
	}
	return res, nil
}

func (s *service) UpdateQuestionPool(ctx context.Context, tx *gorm.DB, req *QuestionRequest, user_id uuid.UUID, questionPoolID uuid.UUID, quizHistoryID uuid.UUID) (*UpdateQuestionPoolResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questionPool, err := s.Repository.GetQuestionPoolByID(c, questionPoolID)
	if err != nil {
		return &UpdateQuestionPoolResponse{}, err
	}

	if req.Order != 0 {
		questionPool.Order = req.Order
	}
	if req.Content != "" {
		questionPool.Content = req.Content
	}
	if req.Note != "" {
		questionPool.Note = req.Note
	}
	if req.Media != "" {
		questionPool.Media = req.Media
	}
	if req.TimeLimit != 0 {
		questionPool.TimeLimit = req.TimeLimit
	}
	if !req.HaveTimeFactor {
		questionPool.HaveTimeFactor = req.HaveTimeFactor
	}
	if req.TimeFactor != 0 {
		questionPool.TimeFactor = req.TimeFactor
	}
	if req.FontSize != 0 {
		questionPool.FontSize = req.FontSize
	}

	qph := &QuestionPoolHistory{
		ID:             uuid.New(),
		QuestionPoolID: questionPool.ID,
		QuizID:         quizHistoryID,
		Order:          questionPool.Order,
		PoolOrder:      questionPool.PoolOrder,
		Content:        questionPool.Content,
		Note:           questionPool.Note,
		Media:          questionPool.Media,
		MediaType:      questionPool.MediaType,
		TimeLimit:      questionPool.TimeLimit,
		HaveTimeFactor: questionPool.HaveTimeFactor,
		TimeFactor:     questionPool.TimeFactor,
		FontSize:       questionPool.FontSize,
	}

	questionPool, er := s.Repository.UpdateQuestionPool(c, tx, questionPool)
	if er != nil {
		return &UpdateQuestionPoolResponse{}, er
	}

	_, e := s.Repository.CreateQuestionPoolHistory(c, tx, qph)
	if e != nil {
		return &UpdateQuestionPoolResponse{}, e
	}

	return &UpdateQuestionPoolResponse{
		QuestionPoolResponse: QuestionPoolResponse{
			QuestionPool: QuestionPool{
				ID:             questionPool.ID,
				QuizID:         questionPool.QuizID,
				Order:          questionPool.Order,
				PoolOrder:      questionPool.PoolOrder,
				Content:        questionPool.Content,
				Note:           questionPool.Note,
				Media:          questionPool.Media,
				MediaType:      questionPool.MediaType,
				TimeLimit:      questionPool.TimeLimit,
				HaveTimeFactor: questionPool.HaveTimeFactor,
				TimeFactor:     questionPool.TimeFactor,
				FontSize:       questionPool.FontSize,
				CreatedAt:      questionPool.CreatedAt,
				UpdatedAt:      questionPool.UpdatedAt,
				DeletedAt:      questionPool.DeletedAt,
			},
		},
		QuestionPoolHistoryID: qph.ID,
	}, nil
}

func (s *service) DeleteQuestionPool(ctx context.Context, tx *gorm.DB, questionPoolID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteQuestionPool(c, tx, questionPoolID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreQuestionPool(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreQuestionPool(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetQuestionPoolHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questionPools, err := s.Repository.GetQuestionPoolHistoriesByQuizID(c, quizID)
	if err != nil {
		return nil, err
	}

	var res []QuestionPoolHistoryResponse
	for _, qp := range questionPools {
		res = append(res, QuestionPoolHistoryResponse{
			QuestionPoolHistory: QuestionPoolHistory{
				ID:             qp.ID,
				QuestionPoolID: qp.QuestionPoolID,
				QuizID:         qp.QuizID,
				Order:          qp.Order,
				PoolOrder:      qp.PoolOrder,
				Content:        qp.Content,
				Note:           qp.Note,
				Media:          qp.Media,
				MediaType:      qp.MediaType,
				TimeLimit:      qp.TimeLimit,
				HaveTimeFactor: qp.HaveTimeFactor,
				TimeFactor:     qp.TimeFactor,
				FontSize:       qp.FontSize,
				CreatedAt:      qp.CreatedAt,
				UpdatedAt:      qp.UpdatedAt,
				DeletedAt:      qp.DeletedAt,
			},
		})
	}
	return res, nil
}

// ---------- Question related service methods ---------- //
func (s *service) CreateQuestion(ctx context.Context, tx *gorm.DB, req *QuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID, questionPoolID *uuid.UUID, questionPoolHistoryID *uuid.UUID, uid uuid.UUID) (*CreateQuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	q := &Question{
		ID:             uuid.New(),
		QuizID:         quizID,
		QuestionPoolID: questionPoolID,
		Type:           req.Type,
		Order:          req.Order,
		PoolOrder:      req.PoolOrder,
		PoolRequired:   req.PoolRequired,
		Content:        req.Content,
		Note:           req.Note,
		Media:          req.Media,
		MediaType:      req.MediaType,
		UseTemplate:    req.UseTemplate,
		TimeLimit:      req.TimeLimit,
		HaveTimeFactor: req.HaveTimeFactor,
		TimeFactor:     req.TimeFactor,
		FontSize:       req.FontSize,
		LayoutIdx:      req.LayoutIdx,
		SelectMin:      req.SelectMin,
		SelectMax:      req.SelectMax,
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     q.ID,
		QuizID:         quizHistoryID,
		QuestionPoolID: questionPoolHistoryID,
		Type:           q.Type,
		Order:          q.Order,
		PoolOrder:      q.PoolOrder,
		PoolRequired:   q.PoolRequired,
		Content:        q.Content,
		Note:           q.Note,
		Media:          q.Media,
		MediaType:      q.MediaType,
		UseTemplate:    q.UseTemplate,
		TimeLimit:      q.TimeLimit,
		HaveTimeFactor: q.HaveTimeFactor,
		TimeFactor:     q.TimeFactor,
		FontSize:       q.FontSize,
		LayoutIdx:      q.LayoutIdx,
		SelectMin:      q.SelectMin,
		SelectMax:      q.SelectMax,
	}

	question, err := s.Repository.CreateQuestion(c, tx, q)
	if err != nil {
		return &CreateQuestionResponse{}, err
	}

	_, er := s.Repository.CreateQuestionHistory(c, tx, qh)
	if er != nil {
		return &CreateQuestionResponse{}, er
	}

	options := make([]any, 0)
	options = append(options, req.Options...)

	return &CreateQuestionResponse{
		QuestionResponse: QuestionResponse{
			Question: Question{
				ID:             question.ID,
				QuizID:         question.QuizID,
				QuestionPoolID: question.QuestionPoolID,
				Type:           question.Type,
				Order:          question.Order,
				PoolOrder:      question.PoolOrder,
				PoolRequired:   question.PoolRequired,
				Content:        question.Content,
				Note:           question.Note,
				Media:          question.Media,
				MediaType:      question.MediaType,
				UseTemplate:    question.UseTemplate,
				TimeLimit:      question.TimeLimit,
				HaveTimeFactor: question.HaveTimeFactor,
				TimeFactor:     question.TimeFactor,
				FontSize:       question.FontSize,
				LayoutIdx:      question.LayoutIdx,
				SelectMin:      question.SelectMin,
				SelectMax:      question.SelectMax,
				CreatedAt:      question.CreatedAt,
				UpdatedAt:      question.UpdatedAt,
				DeletedAt:      question.DeletedAt,
			},
			Options: options,
		},
		QuestionHistoryID: qh.ID,
	}, nil
}

func (s *service) GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questions, err := s.Repository.GetQuestionsByQuizID(c, quizID)
	if err != nil {
		return nil, err
	}

	var res []QuestionResponse
	for _, q := range questions {
		res = append(res, QuestionResponse{
			Question: Question{
				ID:             q.ID,
				QuizID:         q.QuizID,
				QuestionPoolID: q.QuestionPoolID,
				Type:           q.Type,
				Order:          q.Order,
				PoolOrder:      q.PoolOrder,
				PoolRequired:   q.PoolRequired,
				Content:        q.Content,
				Note:           q.Note,
				Media:          q.Media,
				MediaType:      q.MediaType,
				UseTemplate:    q.UseTemplate,
				TimeLimit:      q.TimeLimit,
				HaveTimeFactor: q.HaveTimeFactor,
				TimeFactor:     q.TimeFactor,
				FontSize:       q.FontSize,
				LayoutIdx:      q.LayoutIdx,
				SelectMin:      q.SelectMin,
				SelectMax:      q.SelectMax,
				CreatedAt:      q.CreatedAt,
				UpdatedAt:      q.UpdatedAt,
				DeletedAt:      q.DeletedAt,
			},
		})
	}
	return res, nil
}

func (s *service) GetDeleteQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questions, err := s.Repository.GetDeleteQuestionsByQuizID(c, quizID)
	if err != nil {
		return nil, err
	}

	var res []QuestionResponse
	for _, q := range questions {
		res = append(res, QuestionResponse{
			Question: Question{
				ID:             q.ID,
				QuizID:         q.QuizID,
				QuestionPoolID: q.QuestionPoolID,
				Type:           q.Type,
				Order:          q.Order,
				PoolOrder:      q.PoolOrder,
				PoolRequired:   q.PoolRequired,
				Content:        q.Content,
				Note:           q.Note,
				Media:          q.Media,
				MediaType:      q.MediaType,
				UseTemplate:    q.UseTemplate,
				TimeLimit:      q.TimeLimit,
				HaveTimeFactor: q.HaveTimeFactor,
				TimeFactor:     q.TimeFactor,
				FontSize:       q.FontSize,
				LayoutIdx:      q.LayoutIdx,
				SelectMin:      q.SelectMin,
				SelectMax:      q.SelectMax,
				CreatedAt:      q.CreatedAt,
				UpdatedAt:      q.UpdatedAt,
				DeletedAt:      q.DeletedAt,
			},
		})
	}
	return res, nil
}

func (s *service) GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*Question, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	ques, err := s.Repository.GetQuestionByQuizIDAndOrder(c, quizID, order)
	if err != nil {
		return &Question{}, err
	}

	return ques, nil
}

func (s *service) GetQuestionCountByQuizID(ctx context.Context, quizID uuid.UUID) (int, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	count, err := s.Repository.GetQuestionCountByQuizID(c, quizID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *service) UpdateQuestion(ctx context.Context, tx *gorm.DB, req *QuestionRequest, user_id uuid.UUID, questionID uuid.UUID, quizHistoryID uuid.UUID, questionPoolHistoryID *uuid.UUID) (*UpdateQuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	question, err := s.Repository.GetQuestionByID(c, questionID)
	if err != nil {
		return &UpdateQuestionResponse{}, err
	}

	if req.Type != "" {
		question.Type = req.Type
	}
	if req.Order != 0 {
		question.Order = req.Order
	}

	question.QuestionPoolID = req.QuestionPoolID

	if req.Content != "" {
		question.Content = req.Content
	}
	if req.Note != "" {
		question.Note = req.Note
	}
	if req.Media != "" {
		question.Media = req.Media
	}
	if req.TimeLimit != 0 {
		question.TimeLimit = req.TimeLimit
	}
	if !req.HaveTimeFactor {
		question.HaveTimeFactor = req.HaveTimeFactor
	}
	if req.TimeFactor != 0 {
		question.TimeFactor = req.TimeFactor
	}
	if req.FontSize != 0 {
		question.FontSize = req.FontSize
	}
	if req.LayoutIdx != 0 {
		question.LayoutIdx = req.LayoutIdx
	}
	// if req.SelectUpTo != 0 {
	// 	question.SelectUpTo = req.SelectUpTo
	// }

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     question.ID,
		QuizID:         quizHistoryID,
		QuestionPoolID: questionPoolHistoryID,
		Type:           question.Type,
		Order:          question.Order,
		PoolOrder:      question.PoolOrder,
		PoolRequired:   question.PoolRequired,
		Content:        question.Content,
		Note:           question.Note,
		Media:          question.Media,
		MediaType:      question.MediaType,
		UseTemplate:    question.UseTemplate,
		TimeLimit:      question.TimeLimit,
		HaveTimeFactor: question.HaveTimeFactor,
		TimeFactor:     question.TimeFactor,
		FontSize:       question.FontSize,
		LayoutIdx:      question.LayoutIdx,
		SelectMin:      question.SelectMin,
		SelectMax:      question.SelectMax,
	}

	question, er := s.Repository.UpdateQuestion(c, tx, question)
	if er != nil {
		return &UpdateQuestionResponse{}, er
	}

	_, e := s.Repository.CreateQuestionHistory(c, tx, qh)
	if e != nil {
		return &UpdateQuestionResponse{}, e
	}

	options := make([]any, 0)
	options = append(options, req.Options...)

	return &UpdateQuestionResponse{
		QuestionResponse: QuestionResponse{
			Question: Question{
				ID:             question.ID,
				QuizID:         question.QuizID,
				QuestionPoolID: question.QuestionPoolID,
				Type:           question.Type,
				Order:          question.Order,
				PoolOrder:      question.PoolOrder,
				PoolRequired:   question.PoolRequired,
				Content:        question.Content,
				Note:           question.Note,
				Media:          question.Media,
				MediaType:      question.MediaType,
				TimeLimit:      question.TimeLimit,
				HaveTimeFactor: question.HaveTimeFactor,
				TimeFactor:     question.TimeFactor,
				FontSize:       question.FontSize,
				LayoutIdx:      question.LayoutIdx,
				SelectMin:      question.SelectMin,
				SelectMax:      question.SelectMax,
				CreatedAt:      question.CreatedAt,
				UpdatedAt:      question.UpdatedAt,
				DeletedAt:      question.DeletedAt,
			},
			Options: options,
		},
		QuestionHistoryID: qh.ID,
	}, nil
}

func (s *service) DeleteQuestion(ctx context.Context, tx *gorm.DB, questionID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteQuestion(c, tx, questionID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreQuestion(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreQuestion(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetQuestionHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	questions, err := s.Repository.GetQuestionHistoriesByQuizID(c, quizID)
	if err != nil {
		return nil, err
	}

	var res []QuestionHistoryResponse
	for _, q := range questions {
		res = append(res, QuestionHistoryResponse{
			QuestionHistory: QuestionHistory{
				ID:             q.ID,
				QuestionID:     q.QuestionID,
				QuizID:         q.QuizID,
				QuestionPoolID: q.QuestionPoolID,
				Type:           q.Type,
				Order:          q.Order,
				PoolOrder:      q.PoolOrder,
				PoolRequired:   q.PoolRequired,
				Content:        q.Content,
				Note:           q.Note,
				Media:          q.Media,
				MediaType:      q.MediaType,
				UseTemplate:    q.UseTemplate,
				TimeLimit:      q.TimeLimit,
				HaveTimeFactor: q.HaveTimeFactor,
				TimeFactor:     q.TimeFactor,
				FontSize:       q.FontSize,
				LayoutIdx:      q.LayoutIdx,
				SelectMin:      q.SelectMin,
				SelectMax:      q.SelectMax,
				CreatedAt:      q.CreatedAt,
				UpdatedAt:      q.UpdatedAt,
				DeletedAt:      q.DeletedAt,
			},
		})
	}
	return res, nil
}

func (s *service) GetQuestionHistoryByID(ctx context.Context, id uuid.UUID) (*QuestionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	q, err := s.Repository.GetQuestionHistoryByID(c, id)
	if err != nil {
		return nil, err
	}

	return &QuestionHistoryResponse{
		QuestionHistory: QuestionHistory{
			ID:             q.ID,
			QuestionID:     q.QuestionID,
			QuizID:         q.QuizID,
			QuestionPoolID: q.QuestionPoolID,
			Type:           q.Type,
			Order:          q.Order,
			PoolOrder:      q.PoolOrder,
			PoolRequired:   q.PoolRequired,
			Content:        q.Content,
			Note:           q.Note,
			Media:          q.Media,
			MediaType:      q.MediaType,
			UseTemplate:    q.UseTemplate,
			TimeLimit:      q.TimeLimit,
			HaveTimeFactor: q.HaveTimeFactor,
			TimeFactor:     q.TimeFactor,
			FontSize:       q.FontSize,
			LayoutIdx:      q.LayoutIdx,
			SelectMin:      q.SelectMin,
			SelectMax:      q.SelectMax,
			CreatedAt:      q.CreatedAt,
			UpdatedAt:      q.UpdatedAt,
			DeletedAt:      q.DeletedAt,
		},
	}, nil
}

// ---------- Options related service methods ---------- //
// Choice related service methods
func (s *service) CreateChoiceOption(ctx context.Context, tx *gorm.DB, req *ChoiceOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateChoiceOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	oc := &ChoiceOption{
		ID:         uuid.New(),
		QuestionID: questionID,
		Order:      req.Order,
		Content:    req.Content,
		Mark:       req.Mark,
		Color:      req.Color,
		Correct:    req.Correct,
	}

	och := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: oc.ID,
		QuestionID:     questionHistoryID,
		Order:          oc.Order,
		Content:        oc.Content,
		Mark:           oc.Mark,
		Color:          oc.Color,
		Correct:        oc.Correct,
	}

	optionChoice, err := s.Repository.CreateChoiceOption(c, tx, oc)
	if err != nil {
		return &CreateChoiceOptionResponse{}, err
	}

	_, er := s.Repository.CreateChoiceOptionHistory(c, tx, och)
	if er != nil {
		return &CreateChoiceOptionResponse{}, er
	}

	return &CreateChoiceOptionResponse{
		ChoiceOption: ChoiceOption{
			ID:         optionChoice.ID,
			QuestionID: optionChoice.QuestionID,
			Order:      optionChoice.Order,
			Content:    optionChoice.Content,
			Mark:       optionChoice.Mark,
			Color:      optionChoice.Color,
			Correct:    optionChoice.Correct,
			CreatedAt:  optionChoice.CreatedAt,
			UpdatedAt:  optionChoice.UpdatedAt,
			DeletedAt:  optionChoice.DeletedAt,
		},
	}, nil
}

func (s *service) GetChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoices, err := s.Repository.GetChoiceOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []ChoiceOptionResponse
	for _, oc := range optionChoices {
		res = append(res, ChoiceOptionResponse{
			ChoiceOption: ChoiceOption{
				ID:         oc.ID,
				QuestionID: oc.QuestionID,
				Order:      oc.Order,
				Content:    oc.Content,
				Mark:       oc.Mark,
				Color:      oc.Color,
				Correct:    oc.Correct,
				CreatedAt:  oc.CreatedAt,
				UpdatedAt:  oc.UpdatedAt,
				DeletedAt:  oc.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetDeleteChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoices, err := s.Repository.GetDeleteChoiceOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []ChoiceOptionResponse
	for _, oc := range optionChoices {
		res = append(res, ChoiceOptionResponse{
			ChoiceOption: ChoiceOption{
				ID:         oc.ID,
				QuestionID: oc.QuestionID,
				Order:      oc.Order,
				Content:    oc.Content,
				Mark:       oc.Mark,
				Color:      oc.Color,
				Correct:    oc.Correct,
				CreatedAt:  oc.CreatedAt,
				UpdatedAt:  oc.UpdatedAt,
				DeletedAt:  oc.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetChoiceAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoices, err := s.Repository.GetChoiceAnswersByQuestionID(c, id)
	if err != nil {
		return nil, err
	}

	var res []ChoiceOptionResponse
	for _, oc := range optionChoices {
		res = append(res, ChoiceOptionResponse{
			ChoiceOption: ChoiceOption{
				ID:         oc.ID,
				QuestionID: oc.QuestionID,
				Order:      oc.Order,
				Content:    oc.Content,
				Mark:       oc.Mark,
				Color:      oc.Color,
				Correct:    oc.Correct,
				CreatedAt:  oc.CreatedAt,
				UpdatedAt:  oc.UpdatedAt,
				DeletedAt:  oc.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) UpdateChoiceOption(ctx context.Context, tx *gorm.DB, req *ChoiceOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateChoiceOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoice, err := s.Repository.GetChoiceOptionByID(c, optionID)
	if err != nil {
		return &UpdateChoiceOptionResponse{}, err
	}

	if req.Order != 0 {
		optionChoice.Order = req.Order
	}
	if req.Content != "" {
		optionChoice.Content = req.Content
	}
	if req.Mark != 0 {
		optionChoice.Mark = req.Mark
	}
	if req.Color != "" {
		optionChoice.Color = req.Color
	}
	if !req.Correct {
		optionChoice.Correct = req.Correct
	}

	och := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: optionChoice.ID,
		QuestionID:     questionHistoryID,
		Order:          optionChoice.Order,
		Content:        optionChoice.Content,
		Mark:           optionChoice.Mark,
		Color:          optionChoice.Color,
		Correct:        optionChoice.Correct,
	}

	optionChoice, er := s.Repository.UpdateChoiceOption(c, tx, optionChoice)
	if er != nil {
		return &UpdateChoiceOptionResponse{}, er
	}

	_, e := s.Repository.CreateChoiceOptionHistory(c, tx, och)
	if e != nil {
		return &UpdateChoiceOptionResponse{}, e
	}

	return &UpdateChoiceOptionResponse{
		ChoiceOption: ChoiceOption{
			ID:         optionChoice.ID,
			QuestionID: optionChoice.QuestionID,
			Order:      optionChoice.Order,
			Content:    optionChoice.Content,
			Mark:       optionChoice.Mark,
			Color:      optionChoice.Color,
			Correct:    optionChoice.Correct,
			CreatedAt:  optionChoice.CreatedAt,
			UpdatedAt:  optionChoice.UpdatedAt,
			DeletedAt:  optionChoice.DeletedAt,
		},
	}, nil
}

func (s *service) DeleteChoiceOption(ctx context.Context, tx *gorm.DB, choiceOptionID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteChoiceOption(c, tx, choiceOptionID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreChoiceOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreChoiceOption(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetChoiceOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoices, err := s.Repository.GetChoiceOptionHistoriesByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []ChoiceOptionHistoryResponse
	for _, oc := range optionChoices {
		res = append(res, ChoiceOptionHistoryResponse{
			ChoiceOptionHistory: ChoiceOptionHistory{
				ID:             oc.ID,
				ChoiceOptionID: oc.ChoiceOptionID,
				QuestionID:     oc.QuestionID,
				Order:          oc.Order,
				Content:        oc.Content,
				Mark:           oc.Mark,
				Color:          oc.Color,
				Correct:        oc.Correct,
				CreatedAt:      oc.CreatedAt,
				UpdatedAt:      oc.UpdatedAt,
				DeletedAt:      oc.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetChoiceOptionHistoryByQuestionIDAndContent(ctx context.Context, questionID uuid.UUID, content string) (*ChoiceOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	oc, err := s.Repository.GetChoiceOptionHistoryByQuestionIDAndContent(c, questionID, content)
	if err != nil {
		return nil, err
	}

	return &ChoiceOptionHistoryResponse{
		ChoiceOptionHistory: ChoiceOptionHistory{
			ID:             oc.ID,
			ChoiceOptionID: oc.ChoiceOptionID,
			QuestionID:     oc.QuestionID,
			Order:          oc.Order,
			Content:        oc.Content,
			Mark:           oc.Mark,
			Color:          oc.Color,
			Correct:        oc.Correct,
			CreatedAt:      oc.CreatedAt,
			UpdatedAt:      oc.UpdatedAt,
			DeletedAt:      oc.DeletedAt,
		},
	}, nil
}

func (s *service) GetChoiceOptionHistoryByQuestionIDAndChoiceOptionID(ctx context.Context, questionID uuid.UUID, optionID uuid.UUID) (*ChoiceOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	oc, err := s.Repository.GetChoiceOptionHistoryByQuestionIDAndChoiceOptionID(c, questionID, optionID)
	if err != nil {
		return nil, err
	}

	return &ChoiceOptionHistoryResponse{
		ChoiceOptionHistory: ChoiceOptionHistory{
			ID:             oc.ID,
			ChoiceOptionID: oc.ChoiceOptionID,
			QuestionID:     oc.QuestionID,
			Order:          oc.Order,
			Content:        oc.Content,
			Mark:           oc.Mark,
			Color:          oc.Color,
			Correct:        oc.Correct,
			CreatedAt:      oc.CreatedAt,
			UpdatedAt:      oc.UpdatedAt,
			DeletedAt:      oc.DeletedAt,
		},
	}, nil
}

// Text related service methods
func (s *service) CreateTextOption(ctx context.Context, tx *gorm.DB, req *TextOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateTextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	ot := &TextOption{
		ID:            uuid.New(),
		QuestionID:    questionID,
		Order:         req.Order,
		Content:       req.Content,
		Mark:          req.Mark,
		CaseSensitive: req.CaseSensitive,
	}

	oth := &TextOptionHistory{
		ID:            uuid.New(),
		OptionTextID:  ot.ID,
		QuestionID:    questionHistoryID,
		Order:         ot.Order,
		Content:       ot.Content,
		Mark:          ot.Mark,
		CaseSensitive: ot.CaseSensitive,
	}

	optionText, err := s.Repository.CreateTextOption(c, tx, ot)
	if err != nil {
		return &CreateTextOptionResponse{}, err
	}

	_, er := s.Repository.CreateTextOptionHistory(c, tx, oth)
	if er != nil {
		return &CreateTextOptionResponse{}, er
	}

	return &CreateTextOptionResponse{
		TextOption: TextOption{
			ID:            optionText.ID,
			QuestionID:    optionText.QuestionID,
			Order:         optionText.Order,
			Content:       optionText.Content,
			Mark:          optionText.Mark,
			CaseSensitive: optionText.CaseSensitive,
			CreatedAt:     optionText.CreatedAt,
			UpdatedAt:     optionText.UpdatedAt,
			DeletedAt:     optionText.DeletedAt,
		},
	}, nil
}

func (s *service) GetTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionTexts, err := s.Repository.GetTextOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []TextOptionResponse
	for _, ot := range optionTexts {
		res = append(res, TextOptionResponse{
			TextOption: TextOption{
				ID:            ot.ID,
				QuestionID:    ot.QuestionID,
				Order:         ot.Order,
				Content:       ot.Content,
				Mark:          ot.Mark,
				CaseSensitive: ot.CaseSensitive,
				CreatedAt:     ot.CreatedAt,
				UpdatedAt:     ot.UpdatedAt,
				DeletedAt:     ot.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetDeleteTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionTexts, err := s.Repository.GetDeleteTextOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []TextOptionResponse
	for _, ot := range optionTexts {
		res = append(res, TextOptionResponse{
			TextOption: TextOption{
				ID:            ot.ID,
				QuestionID:    ot.QuestionID,
				Order:         ot.Order,
				Content:       ot.Content,
				Mark:          ot.Mark,
				CaseSensitive: ot.CaseSensitive,
				CreatedAt:     ot.CreatedAt,
				UpdatedAt:     ot.UpdatedAt,
				DeletedAt:     ot.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) UpdateTextOption(ctx context.Context, tx *gorm.DB, req *TextOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateTextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionText, err := s.Repository.GetTextOptionByID(c, optionID)
	if err != nil {
		return &UpdateTextOptionResponse{}, err
	}

	if req.Order != 0 {
		optionText.Order = req.Order
	}
	if req.Content != "" {
		optionText.Content = req.Content
	}
	if req.Mark != 0 {
		optionText.Mark = req.Mark
	}
	if !req.CaseSensitive {
		optionText.CaseSensitive = req.CaseSensitive
	}

	oth := &TextOptionHistory{
		ID:            uuid.New(),
		OptionTextID:  optionText.ID,
		QuestionID:    questionHistoryID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
	}

	optionText, er := s.Repository.UpdateTextOption(c, tx, optionText)
	if er != nil {
		return &UpdateTextOptionResponse{}, er
	}

	_, e := s.Repository.CreateTextOptionHistory(c, tx, oth)
	if e != nil {
		return &UpdateTextOptionResponse{}, e
	}

	return &UpdateTextOptionResponse{
		TextOption: TextOption{
			ID:            optionText.ID,
			QuestionID:    optionText.QuestionID,
			Order:         optionText.Order,
			Content:       optionText.Content,
			Mark:          optionText.Mark,
			CaseSensitive: optionText.CaseSensitive,
			CreatedAt:     optionText.CreatedAt,
			UpdatedAt:     optionText.UpdatedAt,
			DeletedAt:     optionText.DeletedAt,
		},
	}, nil
}

func (s *service) DeleteTextOption(ctx context.Context, tx *gorm.DB, textOptionID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteTextOption(c, tx, textOptionID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreTextOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreTextOption(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetTextOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionTexts, err := s.Repository.GetTextOptionHistoriesByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []TextOptionHistoryResponse
	for _, ot := range optionTexts {
		res = append(res, TextOptionHistoryResponse{
			TextOptionHistory: TextOptionHistory{
				ID:            ot.ID,
				OptionTextID:  ot.OptionTextID,
				QuestionID:    ot.QuestionID,
				Order:         ot.Order,
				Content:       ot.Content,
				Mark:          ot.Mark,
				CaseSensitive: ot.CaseSensitive,
				CreatedAt:     ot.CreatedAt,
				UpdatedAt:     ot.UpdatedAt,
				DeletedAt:     ot.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetTextOptionHistoryByQuestionIDAndContent(ctx context.Context, questionID uuid.UUID, content string) (*TextOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	ot, err := s.Repository.GetTextOptionHistoryByQuestionIDAndContent(c, questionID, content)
	if err != nil {
		return nil, err
	}

	return &TextOptionHistoryResponse{
		TextOptionHistory: TextOptionHistory{
			ID:            ot.ID,
			OptionTextID:  ot.OptionTextID,
			QuestionID:    ot.QuestionID,
			Order:         ot.Order,
			Content:       ot.Content,
			Mark:          ot.Mark,
			CaseSensitive: ot.CaseSensitive,
			CreatedAt:     ot.CreatedAt,
			UpdatedAt:     ot.UpdatedAt,
			DeletedAt:     ot.DeletedAt,
		},
	}, nil
}

// ------ Matching Option ------

func (s *service) CreateMatchingOption(ctx context.Context, tx *gorm.DB, req *MatchingOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateMatchingOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	om := &MatchingOption{
		ID:         uuid.New(),
		QuestionID: questionID,
		Order:      req.Order,
		Content:    req.Content,
		Type:       req.Type,
		Color:      req.Color,
		Eliminate:  req.Eliminate,
	}

	omh := &MatchingOptionHistory{
		ID:               uuid.New(),
		OptionMatchingID: om.ID,
		QuestionID:       questionHistoryID,
		Order:            om.Order,
		Content:          om.Content,
		Type:             om.Type,
		Color:            om.Color,
		Eliminate:        om.Eliminate,
	}

	optionMatching, err := s.Repository.CreateMatchingOption(c, tx, om)
	if err != nil {
		return &CreateMatchingOptionResponse{}, err
	}

	optionMatchingHistory, er := s.Repository.CreateMatchingOptionHistory(c, tx, omh)
	if er != nil {
		return &CreateMatchingOptionResponse{}, er
	}

	return &CreateMatchingOptionResponse{
		MatchingOption: MatchingOption{
			ID:         optionMatching.ID,
			QuestionID: optionMatching.QuestionID,
			Order:      optionMatching.Order,
			Content:    optionMatching.Content,
			Type:       optionMatching.Type,
			Color:      optionMatching.Color,
			Eliminate:  optionMatching.Eliminate,
			CreatedAt:  optionMatching.CreatedAt,
			UpdatedAt:  optionMatching.UpdatedAt,
			DeletedAt:  optionMatching.DeletedAt,
		},
		MatchingOptionHistoryID: optionMatchingHistory.ID,
	}, nil
}

func (s *service) GetMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionMatchings, err := s.Repository.GetMatchingOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []MatchingOptionResponse
	for _, om := range optionMatchings {
		res = append(res, MatchingOptionResponse{
			MatchingOption: MatchingOption{
				ID:         om.ID,
				QuestionID: om.QuestionID,
				Order:      om.Order,
				Content:    om.Content,
				Type:       om.Type,
				Color:      om.Color,
				Eliminate:  om.Eliminate,
				CreatedAt:  om.CreatedAt,
				UpdatedAt:  om.UpdatedAt,
				DeletedAt:  om.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetDeleteMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionMatchings, err := s.Repository.GetDeleteMatchingOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []MatchingOptionResponse
	for _, om := range optionMatchings {
		res = append(res, MatchingOptionResponse{
			MatchingOption: MatchingOption{
				ID:         om.ID,
				QuestionID: om.QuestionID,
				Order:      om.Order,
				Content:    om.Content,
				Type:       om.Type,
				Color:      om.Color,
				Eliminate:  om.Eliminate,
				CreatedAt:  om.CreatedAt,
				UpdatedAt:  om.UpdatedAt,
				DeletedAt:  om.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetMatchingOptionByQuestionIDAndOrder(ctx context.Context, questionID uuid.UUID, order int) (*MatchingOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionMatching, err := s.Repository.GetMatchingOptionByQuestionIDAndOrder(c, questionID, order)
	if err != nil {
		return nil, err
	}

	return &MatchingOptionResponse{
		MatchingOption: MatchingOption{
			ID:         optionMatching.ID,
			QuestionID: optionMatching.QuestionID,
			Order:      optionMatching.Order,
			Content:    optionMatching.Content,
			Type:       optionMatching.Type,
			Color:      optionMatching.Color,
			Eliminate:  optionMatching.Eliminate,
			CreatedAt:  optionMatching.CreatedAt,
			UpdatedAt:  optionMatching.UpdatedAt,
			DeletedAt:  optionMatching.DeletedAt,
		},
	}, nil
}

func (s *service) UpdateMatchingOption(ctx context.Context, tx *gorm.DB, req *MatchingOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateMatchingOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionMatching, err := s.Repository.GetMatchingOptionByID(c, optionID)
	if err != nil {
		return &UpdateMatchingOptionResponse{}, err
	}

	if req.Order != 0 {
		optionMatching.Order = req.Order
	}
	if req.Content != "" {
		optionMatching.Content = req.Content
	}
	if req.Type != "" {
		optionMatching.Type = req.Type
	}
	if !req.Eliminate {
		optionMatching.Eliminate = req.Eliminate
	}

	omh := &MatchingOptionHistory{
		ID:               uuid.New(),
		OptionMatchingID: optionMatching.ID,
		QuestionID:       questionHistoryID,
		Order:            optionMatching.Order,
		Content:          optionMatching.Content,
		Type:             optionMatching.Type,
		Color:            optionMatching.Color,
		Eliminate:        optionMatching.Eliminate,
	}

	optionMatching, er := s.Repository.UpdateMatchingOption(c, tx, optionMatching)
	if er != nil {
		return &UpdateMatchingOptionResponse{}, er
	}

	_, e := s.Repository.CreateMatchingOptionHistory(c, tx, omh)
	if e != nil {
		return &UpdateMatchingOptionResponse{}, e
	}

	return &UpdateMatchingOptionResponse{
		MatchingOption: MatchingOption{
			ID:         optionMatching.ID,
			QuestionID: optionMatching.QuestionID,
			Order:      optionMatching.Order,
			Content:    optionMatching.Content,
			Type:       optionMatching.Type,
			Color:      optionMatching.Color,
			Eliminate:  optionMatching.Eliminate,
			CreatedAt:  optionMatching.CreatedAt,
			UpdatedAt:  optionMatching.UpdatedAt,
			DeletedAt:  optionMatching.DeletedAt,
		},
	}, nil
}

func (s *service) DeleteMatchingOption(ctx context.Context, tx *gorm.DB, matchingOptionID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteMatchingOption(c, tx, matchingOptionID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreMatchingOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreMatchingOption(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetMatchingOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionMatchings, err := s.Repository.GetMatchingOptionHistoriesByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []MatchingOptionHistoryResponse
	for _, om := range optionMatchings {
		res = append(res, MatchingOptionHistoryResponse{
			MatchingOptionHistory: MatchingOptionHistory{
				ID:               om.ID,
				OptionMatchingID: om.OptionMatchingID,
				QuestionID:       om.QuestionID,
				Order:            om.Order,
				Content:          om.Content,
				Type:             om.Type,
				Color:            om.Color,
				Eliminate:        om.Eliminate,
				CreatedAt:        om.CreatedAt,
				UpdatedAt:        om.UpdatedAt,
				DeletedAt:        om.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetMatchingOptionHistoryByQuestionIDAndContent(ctx context.Context, questionID uuid.UUID, content string) (*MatchingOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	om, err := s.Repository.GetMatchingOptionHistoryByQuestionIDAndContent(c, questionID, content)
	if err != nil {
		return nil, err
	}

	return &MatchingOptionHistoryResponse{
		MatchingOptionHistory: MatchingOptionHistory{
			ID:               om.ID,
			OptionMatchingID: om.OptionMatchingID,
			QuestionID:       om.QuestionID,
			Order:            om.Order,
			Content:          om.Content,
			Type:             om.Type,
			Color:            om.Color,
			Eliminate:        om.Eliminate,
			CreatedAt:        om.CreatedAt,
			UpdatedAt:        om.UpdatedAt,
			DeletedAt:        om.DeletedAt,
		},
	}, nil
}

func (s *service) GetMatchingOptionHistoryByOptionMatchingID(ctx context.Context, optionMatchingID uuid.UUID) (*MatchingOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	om, err := s.Repository.GetMatchingOptionHistoryByOptionMatchingID(c, optionMatchingID)
	if err != nil {
		return nil, err
	}

	return &MatchingOptionHistoryResponse{
		MatchingOptionHistory: MatchingOptionHistory{
			ID:               om.ID,
			OptionMatchingID: om.OptionMatchingID,
			QuestionID:       om.QuestionID,
			Order:            om.Order,
			Content:          om.Content,
			Type:             om.Type,
			Color:            om.Color,
			Eliminate:        om.Eliminate,
			CreatedAt:        om.CreatedAt,
			UpdatedAt:        om.UpdatedAt,
			DeletedAt:        om.DeletedAt,
		},
	}, nil
}

func (s *service) GetMatchingOptionHistoryByID(ctx context.Context, id uuid.UUID) (*MatchingOptionHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	om, err := s.Repository.GetMatchingOptionHistoryByID(c, id)
	if err != nil {
		return nil, err
	}

	return &MatchingOptionHistoryResponse{
		MatchingOptionHistory: MatchingOptionHistory{
			ID:               om.ID,
			OptionMatchingID: om.OptionMatchingID,
			QuestionID:       om.QuestionID,
			Order:            om.Order,
			Content:          om.Content,
			Type:             om.Type,
			Color:            om.Color,
			Eliminate:        om.Eliminate,
			CreatedAt:        om.CreatedAt,
			UpdatedAt:        om.UpdatedAt,
			DeletedAt:        om.DeletedAt,
		},
	}, nil
}

// ------ Matching Answer ------

func (s *service) CreateMatchingAnswer(ctx context.Context, tx *gorm.DB, req *MatchingAnswerRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, matchingOptionPromptHistoryID uuid.UUID, matchingOptionOptionHistoryID uuid.UUID, uid uuid.UUID) (*CreateMatchingAnswerResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	am := &MatchingAnswer{
		ID:         uuid.New(),
		QuestionID: questionID,
		PromptID:   req.PromptID,
		OptionID:   req.OptionID,
		Mark:       req.Mark,
	}

	amh := &MatchingAnswerHistory{
		ID:               uuid.New(),
		AnswerMatchingID: am.ID,
		QuestionID:       questionHistoryID,
		PromptID:         matchingOptionPromptHistoryID,
		OptionID:         matchingOptionOptionHistoryID,
		Mark:             am.Mark,
	}

	answerMatching, err := s.Repository.CreateMatchingAnswer(c, tx, am)
	if err != nil {
		return &CreateMatchingAnswerResponse{}, err
	}

	_, er := s.Repository.CreateMatchingAnswerHistory(c, tx, amh)
	if er != nil {
		return &CreateMatchingAnswerResponse{}, er
	}

	return &CreateMatchingAnswerResponse{
		MatchingAnswer: MatchingAnswer{
			ID:         answerMatching.ID,
			QuestionID: answerMatching.QuestionID,
			PromptID:   answerMatching.PromptID,
			OptionID:   answerMatching.OptionID,
			Mark:       answerMatching.Mark,
			CreatedAt:  answerMatching.CreatedAt,
			UpdatedAt:  answerMatching.UpdatedAt,
			DeletedAt:  answerMatching.DeletedAt,
		},
	}, nil
}

func (s *service) GetMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	answerMatchings, err := s.Repository.GetMatchingAnswersByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []MatchingAnswerResponse
	for _, am := range answerMatchings {
		res = append(res, MatchingAnswerResponse{
			MatchingAnswer: MatchingAnswer{
				ID:         am.ID,
				QuestionID: am.QuestionID,
				PromptID:   am.PromptID,
				OptionID:   am.OptionID,
				Mark:       am.Mark,
				CreatedAt:  am.CreatedAt,
				UpdatedAt:  am.UpdatedAt,
				DeletedAt:  am.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetDeleteMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	answerMatchings, err := s.Repository.GetDeleteMatchingAnswersByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []MatchingAnswerResponse
	for _, am := range answerMatchings {
		res = append(res, MatchingAnswerResponse{
			MatchingAnswer: MatchingAnswer{
				ID:         am.ID,
				QuestionID: am.QuestionID,
				PromptID:   am.PromptID,
				OptionID:   am.OptionID,
				Mark:       am.Mark,
				CreatedAt:  am.CreatedAt,
				UpdatedAt:  am.UpdatedAt,
				DeletedAt:  am.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) UpdateMatchingAnswer(ctx context.Context, tx *gorm.DB, req *MatchingAnswerRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateMatchingAnswerResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	answerMatching, err := s.Repository.GetMatchingAnswerByID(c, optionID)
	if err != nil {
		return &UpdateMatchingAnswerResponse{}, err
	}

	if req.PromptID != uuid.Nil {
		answerMatching.PromptID = req.PromptID
	}
	if req.OptionID != uuid.Nil {
		answerMatching.OptionID = req.OptionID
	}

	answerMatching.Mark = req.Mark

	amh := &MatchingAnswerHistory{
		ID:               uuid.New(),
		AnswerMatchingID: answerMatching.ID,
		QuestionID:       questionHistoryID,
		PromptID:         answerMatching.PromptID,
		OptionID:         answerMatching.OptionID,
		Mark:             answerMatching.Mark,
	}

	answerMatching, er := s.Repository.UpdateMatchingAnswer(c, tx, answerMatching)
	if er != nil {
		return &UpdateMatchingAnswerResponse{}, er
	}

	_, e := s.Repository.CreateMatchingAnswerHistory(c, tx, amh)
	if e != nil {
		return &UpdateMatchingAnswerResponse{}, e
	}

	return &UpdateMatchingAnswerResponse{
		MatchingAnswer: MatchingAnswer{
			ID:         answerMatching.ID,
			QuestionID: answerMatching.QuestionID,
			PromptID:   answerMatching.PromptID,
			OptionID:   answerMatching.OptionID,
			Mark:       answerMatching.Mark,
			CreatedAt:  answerMatching.CreatedAt,
			UpdatedAt:  answerMatching.UpdatedAt,
			DeletedAt:  answerMatching.DeletedAt,
		},
	}, nil
}

func (s *service) DeleteMatchingAnswer(ctx context.Context, tx *gorm.DB, matchingAnswerID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	e := s.Repository.DeleteMatchingAnswer(c, tx, matchingAnswerID)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreMatchingAnswer(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, e := s.Repository.RestoreMatchingAnswer(c, tx, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) GetMatchingAnswerHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	answerMatchings, err := s.Repository.GetMatchingAnswerHistoriesByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []MatchingAnswerHistoryResponse
	for _, am := range answerMatchings {
		res = append(res, MatchingAnswerHistoryResponse{
			MatchingAnswerHistory: MatchingAnswerHistory{
				ID:               am.ID,
				AnswerMatchingID: am.AnswerMatchingID,
				QuestionID:       am.QuestionID,
				PromptID:         am.PromptID,
				OptionID:         am.OptionID,
				Mark:             am.Mark,
				CreatedAt:        am.CreatedAt,
				UpdatedAt:        am.UpdatedAt,
				DeletedAt:        am.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetMatchingAnswerHistoryByPromptIDAndOptionID(ctx context.Context, promptID uuid.UUID, optionID uuid.UUID) (*MatchingAnswerHistoryResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	am, err := s.Repository.GetMatchingAnswerHistoryByPromptIDAndOptionID(c, promptID, optionID)
	if err != nil {
		return nil, err
	}

	return &MatchingAnswerHistoryResponse{
		MatchingAnswerHistory: MatchingAnswerHistory{
			ID:               am.ID,
			AnswerMatchingID: am.AnswerMatchingID,
			QuestionID:       am.QuestionID,
			PromptID:         am.PromptID,
			OptionID:         am.OptionID,
			Mark:             am.Mark,
			CreatedAt:        am.CreatedAt,
			UpdatedAt:        am.UpdatedAt,
			DeletedAt:        am.DeletedAt,
		},
	}, nil
}

// ---------- Live + Quiz related service methods ---------- //
func (s *service) GetLatestQuizVersionByID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	qid, err := s.Repository.GetLatestQuizHistoryByQuizID(c, id)
	if err != nil {
		return nil, err
	}

	return qid, nil
}

func (s *service) getOptionsByQuestionIDForLQS(c context.Context, t string, qid uuid.UUID) (any, error) {
	switch t {
	case util.Choice, util.TrueFalse:
		ocs, err := s.Repository.GetChoiceOptionHistoriesByQuestionID(c, qid)
		if err != nil {
			return nil, err
		}
		options := make([]LQSChoiceOption, 0)
		for _, oc := range ocs {
			options = append(options, LQSChoiceOption{
				ID:      oc.ID,
				Content: oc.Content,
				Color:   oc.Color,
				Order:   oc.Order,
			})
		}
		sort.Sort(ByCOOrder(options))
		return options, nil
	case util.FillBlank, util.Paragraph:
		ots, err := s.Repository.GetTextOptionHistoriesByQuestionID(c, qid)
		if err != nil {
			return nil, err
		}
		options := make([]LQSTextOption, 0)
		for _, ot := range ots {
			options = append(options, LQSTextOption{
				ID:            ot.ID,
				CaseSensitive: ot.CaseSensitive,
				Order:         ot.Order,
			})
		}
		sort.Sort(ByTOOrder(options))
		return options, nil
	case util.Matching:
		oms, err := s.Repository.GetMatchingOptionHistoriesByQuestionID(c, qid)
		if err != nil {
			return nil, err
		}
		omp := make([]LQSMatchingOptionPrompt, 0)
		omo := make([]LQSMatchingOptionOption, 0)
		for _, om := range oms {
			if om.Type == "MATCHING_PROMPT" {
				omp = append(omp, LQSMatchingOptionPrompt{
					ID:      om.ID,
					Content: om.Content,
					Order:   om.Order,
				})
			}
			if om.Type == "MATCHING_OPTION" {
				omo = append(omo, LQSMatchingOptionOption{
					ID:        om.ID,
					Content:   om.Content,
					Color:     om.Color,
					Eliminate: om.Eliminate,
					Order:     om.Order,
				})
			}
		}
		sort.Sort(ByMPOrder(omp))
		sort.Sort(ByMOOrder(omo))
		return LQSMatchingOption{
			Prompts: omp,
			Options: omo,
		}, nil
	}
	return nil, nil
}

func (s *service) getAnswersByQuestionIDForLQS(c context.Context, t string, qid uuid.UUID) (any, error) {
	switch t {
	case util.Choice, util.TrueFalse:
		ocs, err := s.Repository.GetChoiceOptionHistoriesByQuestionID(c, qid)
		if err != nil {
			return nil, err
		}
		answers := make([]LQSChoiceAnswer, 0)
		for _, oc := range ocs {
			answers = append(answers, LQSChoiceAnswer{
				LQSChoiceOption: LQSChoiceOption{
					ID:      oc.ID,
					Content: oc.Content,
					Color:   oc.Color,
					Order:   oc.Order,
				},
				Mark:       oc.Mark,
				Correct:    oc.Correct,
				Type:       t,
				QuestionID: qid,
			})
		}
		sort.Sort(ByCAOrder(answers))
		return answers, nil
	case util.FillBlank, util.Paragraph:
		ots, err := s.Repository.GetTextOptionHistoriesByQuestionID(c, qid)
		if err != nil {
			return nil, err
		}
		answers := make([]LQSTextAnswer, 0)
		for _, ot := range ots {
			answers = append(answers, LQSTextAnswer{
				LQSTextOption: LQSTextOption{
					ID:            ot.ID,
					CaseSensitive: ot.CaseSensitive,
					Order:         ot.Order,
				},
				Content:    ot.Content,
				Mark:       ot.Mark,
				Type:       t,
				QuestionID: qid,
			})
		}
		sort.Sort(ByTAOrder(answers))
		return answers, nil
	case util.Matching:
		oms, err := s.Repository.GetMatchingAnswerHistoriesByQuestionID(c, qid)
		if err != nil {
			return nil, err
		}
		answers := make([]LQSMatchingAnswer, 0)
		for _, om := range oms {
			answers = append(answers, LQSMatchingAnswer{
				PromptID:   om.PromptID,
				OptionID:   om.OptionID,
				Mark:       om.Mark,
				Type:       t,
				QuestionID: qid,
			})
		}
		return answers, nil
	}
	return nil, nil
}

func (s *service) GetQuestionsByQuizIDForLQS(ctx context.Context, id uuid.UUID) ([]any, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	qh, err := s.Repository.GetQuestionHistoriesByQuizID(c, id)
	if err != nil {
		return nil, err
	}

	qph, err := s.Repository.GetQuestionPoolHistoriesByQuizID(c, id)
	if err != nil {
		return nil, err
	}

	var qs []any
	for _, q := range qh {
		if q.PoolOrder == -1 {
			options, err := s.getOptionsByQuestionIDForLQS(c, q.Type, q.ID)
			if err != nil {
				return nil, err
			}
			qs = append(qs, LQSQuestion{
				QuestionHistory: q,
				Options:         options,
			})
		}
	}

	for _, qp := range qph {
		var subQ []LQSQuestion
		for _, q := range qh {
			if q.PoolOrder == qp.Order {
				options, err := s.getOptionsByQuestionIDForLQS(c, q.Type, q.ID)
				if err != nil {
					return nil, err
				}
				subQ = append(subQ, LQSQuestion{
					QuestionHistory: q,
					Options:         options,
				})
			}
		}
		sort.Sort(ByQHOrder(subQ))
		qs = append(qs, LQSQuestionPool{
			QuestionPoolHistory: qp,
			Type:                util.Pool,
			SubQuestions:        subQ,
		})
	}

	sort.Sort(ByQQPOrder(qs))

	return qs, nil
}

func (s *service) GetAnswersByQuizIDForLQS(ctx context.Context, id uuid.UUID) ([]any, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	qh, err := s.Repository.GetQuestionHistoriesByQuizID(c, id)
	if err != nil {
		return nil, err
	}

	qph, err := s.Repository.GetQuestionPoolHistoriesByQuizID(c, id)
	if err != nil {
		return nil, err
	}

	var ans []LQSAnswers
	for _, q := range qh {
		if q.PoolOrder == -1 {
			answers, err := s.getAnswersByQuestionIDForLQS(c, q.Type, q.ID)
			if err != nil {
				return nil, err
			}
			ans = append(ans, LQSAnswers{
				Answers: answers,
				Order:   q.Order,
			})
		}
	}

	for _, qp := range qph {
		var a []LQSAnswers
		for _, q := range qh {
			if q.PoolOrder == qp.Order {
				answers, err := s.getAnswersByQuestionIDForLQS(c, q.Type, q.ID)
				if err != nil {
					return nil, err
				}
				a = append(a, LQSAnswers{
					Answers: answers,
					Order:   q.Order,
				})
			}
		}
		sort.Sort(ByAOrder(a))
		ans = append(ans, LQSAnswers{
			Answers: a,
			Order:   qp.Order,
		})
	}
	sort.Sort(ByAOrder(ans))

	var res []any
	for _, a := range ans {
		var subA []any
		val, ok := a.Answers.([]LQSAnswers)
		if ok {
			for _, aa := range val {
				subA = append(subA, aa.Answers)
			}
			res = append(res, subA)
		} else {
			res = append(res, a.Answers)
		}
	}

	return res, nil
}
