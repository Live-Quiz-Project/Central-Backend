package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
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

// ---------- Quiz related service methods ---------- //
func (s *service) CreateQuiz(ctx context.Context, req *CreateQuizRequest, uid uuid.UUID) (*CreateQuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	q := &Quiz{
		ID:             uuid.New(),
		CreatorID:      uid,
		Title:          req.Title,
		Description:    req.Description,
		CoverImage:     "default.png",
		Visibility:     req.Visibility,
		TimeLimit:      req.TimeLimit,
		HaveTimeFactor: req.HaveTimeFactor,
		TimeFactor:     req.TimeFactor,
		FontSize:       req.FontSize,
		Mark:           req.Mark,
		SelectUpTo:     req.SelectUpTo,
		CaseSensitive:  req.CaseSensitive,
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
		SelectUpTo:     q.SelectUpTo,
		CaseSensitive:  q.CaseSensitive,
	}

	quiz, err := s.Repository.CreateQuiz(c, q)
	if err != nil {
		return &CreateQuizResponse{}, err
	}
	quizH, er := s.Repository.CreateQuizHistory(c, qh)
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
				SelectUpTo:     quiz.SelectUpTo,
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
				SelectUpTo:     q.SelectUpTo,
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
			SelectUpTo:     quiz.SelectUpTo,
			CaseSensitive:  quiz.CaseSensitive,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			DeletedAt:      quiz.DeletedAt,
		},
	}, nil
}

func (s *service) UpdateQuiz(ctx context.Context, req *UpdateQuizRequest, uid uuid.UUID, id uuid.UUID) (*UpdateQuizResponse, error) {
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
	quiz.Title = req.Title
	quiz.Description = req.Description
	quiz.CoverImage = req.CoverImage
	quiz.Visibility = req.Visibility
	quiz.TimeLimit = req.TimeLimit
	quiz.HaveTimeFactor = req.HaveTimeFactor
	quiz.TimeFactor = req.TimeFactor
	quiz.FontSize = req.FontSize
	quiz.Mark = req.Mark
	quiz.SelectUpTo = req.SelectUpTo
	quiz.CaseSensitive = req.CaseSensitive

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
		SelectUpTo:     quiz.SelectUpTo,
		CaseSensitive:  quiz.CaseSensitive,
	}

	quiz, er := s.Repository.UpdateQuiz(c, quiz)
	if er != nil {
		return &UpdateQuizResponse{}, er
	}

	_, e := s.Repository.CreateQuizHistory(c, qh)
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
				SelectUpTo:     quiz.SelectUpTo,
				CaseSensitive:  quiz.CaseSensitive,
				CreatedAt:      quiz.CreatedAt,
				UpdatedAt:      quiz.UpdatedAt,
				DeletedAt:      quiz.DeletedAt,
			},
		},
		QuizHistoryID: qh.ID,
	}, nil
}

func (s *service) DeleteQuiz(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetQuizByID(c, id)
	if err != nil {
		return err
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
		SelectUpTo:     quiz.SelectUpTo,
		CaseSensitive:  quiz.CaseSensitive,
	}

	_, er := s.Repository.CreateQuizHistory(c, qh)
	if er != nil {
		return er
	}

	e := s.Repository.DeleteQuiz(c, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreQuiz(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetQuizByID(c, id)
	if err != nil {
		return nil, err
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
		SelectUpTo:     quiz.SelectUpTo,
		CaseSensitive:  quiz.CaseSensitive,
	}

	_, er := s.Repository.CreateQuizHistory(c, qh)
	if er != nil {
		return nil, er
	}

	quiz, e := s.Repository.RestoreQuiz(c, id)
	if e != nil {
		return nil, e
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
			SelectUpTo:     quiz.SelectUpTo,
			CaseSensitive:  quiz.CaseSensitive,
			CreatedAt:      quiz.CreatedAt,
			UpdatedAt:      quiz.UpdatedAt,
			DeletedAt:      quiz.DeletedAt,
		},
	}, nil
}

// ---------- Question Pool related service methods ---------- //
func (s *service) CreateQuestionPool(ctx context.Context, req *QuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID) (*CreateQuestionPoolResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	qp := &QuestionPool{
		ID:             uuid.New(),
		QuizID:         quizID,
		Order:          req.Order,
		Content:        req.Content,
		Note:           req.Note,
		Media:          req.Media,
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
		Content:        qp.Content,
		Note:           qp.Note,
		Media:          qp.Media,
		TimeLimit:      qp.TimeLimit,
		HaveTimeFactor: qp.HaveTimeFactor,
		TimeFactor:     qp.TimeFactor,
		FontSize:       qp.FontSize,
	}

	questionPool, err := s.Repository.CreateQuestionPool(c, qp)
	if err != nil {
		return &CreateQuestionPoolResponse{}, err
	}
	questionPoolH, er := s.Repository.CreateQuestionPoolHistory(c, qph)
	if er != nil {
		return &CreateQuestionPoolResponse{}, er
	}

	return &CreateQuestionPoolResponse{
		QuestionPoolResponse: QuestionPoolResponse{
			QuestionPool: QuestionPool{
				ID:             questionPool.ID,
				QuizID:         questionPool.QuizID,
				Order:          questionPool.Order,
				Content:        questionPool.Content,
				Note:           questionPool.Note,
				Media:          questionPool.Media,
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
				Content:        qp.Content,
				Note:           qp.Note,
				Media:          qp.Media,
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

func (s *service) UpdateQuestionPool(ctx context.Context, req *QuestionRequest, user_id uuid.UUID, questionPoolID uuid.UUID, quizHistoryID uuid.UUID) (*UpdateQuestionPoolResponse, error) {
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
		Content:        questionPool.Content,
		Note:           questionPool.Note,
		Media:          questionPool.Media,
		TimeLimit:      questionPool.TimeLimit,
		HaveTimeFactor: questionPool.HaveTimeFactor,
		TimeFactor:     questionPool.TimeFactor,
		FontSize:       questionPool.FontSize,
	}

	questionPool, er := s.Repository.UpdateQuestionPool(c, questionPool)
	if er != nil {
		return &UpdateQuestionPoolResponse{}, er
	}

	_, e := s.Repository.CreateQuestionPoolHistory(c, qph)
	if e != nil {
		return &UpdateQuestionPoolResponse{}, e
	}

	return &UpdateQuestionPoolResponse{
		QuestionPoolResponse: QuestionPoolResponse{
			QuestionPool: QuestionPool{
				ID:             questionPool.ID,
				QuizID:         questionPool.QuizID,
				Order:          questionPool.Order,
				Content:        questionPool.Content,
				Note:           questionPool.Note,
				Media:          questionPool.Media,
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

// ---------- Question related service methods ---------- //
func (s *service) CreateQuestion(ctx context.Context, req *QuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID, questionPoolID *uuid.UUID, questionPoolHistoryID *uuid.UUID, uid uuid.UUID) (*CreateQuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	q := &Question{
		ID:             uuid.New(),
		QuizID:         quizID,
		QuestionPoolID: questionPoolID,
		Type:           req.Type,
		Order:          req.Order,
		Content:        req.Content,
		Note:           req.Note,
		Media:          req.Media,
		UseTemplate:    req.UseTemplate,
		TimeLimit:      req.TimeLimit,
		HaveTimeFactor: req.HaveTimeFactor,
		TimeFactor:     req.TimeFactor,
		FontSize:       req.FontSize,
		LayoutIdx:      req.LayoutIdx,
		SelectUpTo:     req.SelectUpTo,
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     q.ID,
		QuizID:         quizHistoryID,
		QuestionPoolID: questionPoolHistoryID,
		Type:           q.Type,
		Order:          q.Order,
		Content:        q.Content,
		Note:           q.Note,
		Media:          q.Media,
		UseTemplate:    q.UseTemplate,
		TimeLimit:      q.TimeLimit,
		HaveTimeFactor: q.HaveTimeFactor,
		TimeFactor:     q.TimeFactor,
		FontSize:       q.FontSize,
		LayoutIdx:      q.LayoutIdx,
		SelectUpTo:     q.SelectUpTo,
	}

	question, err := s.Repository.CreateQuestion(c, q)
	if err != nil {
		return &CreateQuestionResponse{}, err
	}

	_, er := s.Repository.CreateQuestionHistory(c, qh)
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
				Content:        question.Content,
				Note:           question.Note,
				Media:          question.Media,
				UseTemplate:    question.UseTemplate,
				TimeLimit:      question.TimeLimit,
				HaveTimeFactor: question.HaveTimeFactor,
				TimeFactor:     question.TimeFactor,
				FontSize:       question.FontSize,
				LayoutIdx:      question.LayoutIdx,
				SelectUpTo:     question.SelectUpTo,
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
				Content:        q.Content,
				Note:           q.Note,
				Media:          q.Media,
				UseTemplate:    q.UseTemplate,
				TimeLimit:      q.TimeLimit,
				HaveTimeFactor: q.HaveTimeFactor,
				TimeFactor:     q.TimeFactor,
				FontSize:       q.FontSize,
				LayoutIdx:      q.LayoutIdx,
				SelectUpTo:     q.SelectUpTo,
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

func (s *service) UpdateQuestion(ctx context.Context, req *QuestionRequest, user_id uuid.UUID, questionID uuid.UUID, quizHistoryID uuid.UUID, questionPoolHistoryID *uuid.UUID) (*UpdateQuestionResponse, error) {
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
	if req.SelectUpTo != 0 {
		question.SelectUpTo = req.SelectUpTo
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     question.ID,
		QuizID:         quizHistoryID,
		QuestionPoolID: questionPoolHistoryID,
		Type:           question.Type,
		Order:          question.Order,
		Content:        question.Content,
		Note:           question.Note,
		Media:          question.Media,
		UseTemplate:    question.UseTemplate,
		TimeLimit:      question.TimeLimit,
		HaveTimeFactor: question.HaveTimeFactor,
		TimeFactor:     question.TimeFactor,
		FontSize:       question.FontSize,
		LayoutIdx:      question.LayoutIdx,
		SelectUpTo:     question.SelectUpTo,
	}

	question, er := s.Repository.UpdateQuestion(c, question)
	if er != nil {
		return &UpdateQuestionResponse{}, er
	}

	_, e := s.Repository.CreateQuestionHistory(c, qh)
	if e != nil {
		return &UpdateQuestionResponse{}, e
	}

	return &UpdateQuestionResponse{
		QuestionResponse: QuestionResponse{
			Question: Question{
				ID:             question.ID,
				QuizID:         question.QuizID,
				QuestionPoolID: question.QuestionPoolID,
				Type:           question.Type,
				Order:          question.Order,
				Content:        question.Content,
				Note:           question.Note,
				Media:          question.Media,
				TimeLimit:      question.TimeLimit,
				HaveTimeFactor: question.HaveTimeFactor,
				TimeFactor:     question.TimeFactor,
				FontSize:       question.FontSize,
				LayoutIdx:      question.LayoutIdx,
				SelectUpTo:     question.SelectUpTo,
				CreatedAt:      question.CreatedAt,
				UpdatedAt:      question.UpdatedAt,
				DeletedAt:      question.DeletedAt,
			},
		},
		QuestionHistoryID: qh.ID,
	}, nil
}

func (s *service) DeleteQuestion(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	question, err := s.Repository.GetQuestionByID(c, id)
	if err != nil {
		return err
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     question.ID,
		QuizID:         question.QuizID,
		QuestionPoolID: question.QuestionPoolID,
		Type:           question.Type,
		Order:          question.Order,
		Content:        question.Content,
		Note:           question.Note,
		Media:          question.Media,
		UseTemplate:    question.UseTemplate,
		TimeLimit:      question.TimeLimit,
		HaveTimeFactor: question.HaveTimeFactor,
		TimeFactor:     question.TimeFactor,
		FontSize:       question.FontSize,
		LayoutIdx:      question.LayoutIdx,
		SelectUpTo:     question.SelectUpTo,
	}

	_, er := s.Repository.CreateQuestionHistory(c, qh)
	if er != nil {
		return er
	}

	e := s.Repository.DeleteQuestion(c, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreQuestion(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	question, err := s.Repository.GetQuestionByID(c, id)
	if err != nil {
		return nil, err
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     question.ID,
		QuizID:         question.QuizID,
		QuestionPoolID: question.QuestionPoolID,
		Type:           question.Type,
		Order:          question.Order,
		Content:        question.Content,
		Note:           question.Note,
		Media:          question.Media,
		UseTemplate:    question.UseTemplate,
		TimeLimit:      question.TimeLimit,
		HaveTimeFactor: question.HaveTimeFactor,
		TimeFactor:     question.TimeFactor,
		FontSize:       question.FontSize,
		LayoutIdx:      question.LayoutIdx,
		SelectUpTo:     question.SelectUpTo,
	}

	_, er := s.Repository.CreateQuestionHistory(c, qh)
	if er != nil {
		return nil, er
	}

	question, e := s.Repository.RestoreQuestion(c, id)
	if e != nil {
		return nil, e
	}

	return &QuestionResponse{
		Question: Question{
			ID:             question.ID,
			QuizID:         question.QuizID,
			QuestionPoolID: question.QuestionPoolID,
			Type:           question.Type,
			Order:          question.Order,
			Content:        question.Content,
			Note:           question.Note,
			Media:          question.Media,
			TimeLimit:      question.TimeLimit,
			HaveTimeFactor: question.HaveTimeFactor,
			TimeFactor:     question.TimeFactor,
			FontSize:       question.FontSize,
			LayoutIdx:      question.LayoutIdx,
			SelectUpTo:     question.SelectUpTo,
			CreatedAt:      question.CreatedAt,
			UpdatedAt:      question.UpdatedAt,
			DeletedAt:      question.DeletedAt,
		},
	}, nil
}

// ---------- Options related service methods ---------- //
// Choice related service methods
func (s *service) CreateChoiceOption(ctx context.Context, req *ChoiceOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error) {
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

	optionChoice, err := s.Repository.CreateChoiceOption(c, oc)
	if err != nil {
		return &ChoiceOptioneResponse{}, err
	}

	_, er := s.Repository.CreateChoiceOptionHistory(c, och)
	if er != nil {
		return &ChoiceOptioneResponse{}, er
	}

	return &ChoiceOptioneResponse{
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

func (s *service) GetChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptioneResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoices, err := s.Repository.GetChoiceOptionsByQuestionID(c, questionID)
	if err != nil {
		return nil, err
	}

	var res []ChoiceOptioneResponse
	for _, oc := range optionChoices {
		res = append(res, ChoiceOptioneResponse{
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

func (s *service) GetChoiceAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptioneResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoices, err := s.Repository.GetChoiceAnswersByQuestionID(c, id)
	if err != nil {
		return nil, err
	}

	var res []ChoiceOptioneResponse
	for _, oc := range optionChoices {
		res = append(res, ChoiceOptioneResponse{
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

func (s *service) UpdateChoiceOption(ctx context.Context, req *ChoiceOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*ChoiceOptioneResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoice, err := s.Repository.GetChoiceOptionByID(c, optionID)
	if err != nil {
		return &ChoiceOptioneResponse{}, err
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

	optionChoice, er := s.Repository.UpdateChoiceOption(c, optionChoice)
	if er != nil {
		return &ChoiceOptioneResponse{}, er
	}

	_, e := s.Repository.CreateChoiceOptionHistory(c, och)
	if e != nil {
		return &ChoiceOptioneResponse{}, e
	}

	return &ChoiceOptioneResponse{
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

func (s *service) DeleteChoiceOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoice, err := s.Repository.GetChoiceOptionByID(c, id)
	if err != nil {
		return err
	}

	och := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: optionChoice.ID,
		QuestionID:     optionChoice.QuestionID,
		Order:          optionChoice.Order,
		Content:        optionChoice.Content,
		Mark:           optionChoice.Mark,
		Color:          optionChoice.Color,
		Correct:        optionChoice.Correct,
	}

	_, er := s.Repository.CreateChoiceOptionHistory(c, och)
	if er != nil {
		return er
	}

	e := s.Repository.DeleteChoiceOption(c, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreChoiceOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoice, err := s.Repository.GetChoiceOptionByID(c, id)
	if err != nil {
		return nil, err
	}

	och := &ChoiceOptionHistory{
		ID:             uuid.New(),
		ChoiceOptionID: optionChoice.ID,
		QuestionID:     optionChoice.QuestionID,
		Order:          optionChoice.Order,
		Content:        optionChoice.Content,
		Mark:           optionChoice.Mark,
		Color:          optionChoice.Color,
		Correct:        optionChoice.Correct,
	}

	_, er := s.Repository.CreateChoiceOptionHistory(c, och)
	if er != nil {
		return nil, er
	}

	optionChoice, e := s.Repository.RestoreChoiceOption(c, id)
	if e != nil {
		return nil, e
	}

	return &ChoiceOptioneResponse{
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

// Text related service methods
func (s *service) CreateTextOption(ctx context.Context, req *TextOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error) {
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

	optionText, err := s.Repository.CreateTextOption(c, ot)
	if err != nil {
		return &TextOptionResponse{}, err
	}

	_, er := s.Repository.CreateTextOptionHistory(c, oth)
	if er != nil {
		return &TextOptionResponse{}, er
	}

	return &TextOptionResponse{
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

func (s *service) GetTextAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]TextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionTexts, err := s.Repository.GetTextAnswersByQuestionID(c, id)
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

func (s *service) UpdateTextOption(ctx context.Context, req *TextOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*TextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionText, err := s.Repository.GetTextOptionByID(c, optionID)
	if err != nil {
		return &TextOptionResponse{}, err
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

	optionText, er := s.Repository.UpdateTextOption(c, optionText)
	if er != nil {
		return &TextOptionResponse{}, er
	}

	_, e := s.Repository.CreateTextOptionHistory(c, oth)
	if e != nil {
		return &TextOptionResponse{}, e
	}

	return &TextOptionResponse{
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

func (s *service) DeleteTextOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionText, err := s.Repository.GetTextOptionByID(c, id)
	if err != nil {
		return err
	}

	oth := &TextOptionHistory{
		ID:            uuid.New(),
		OptionTextID:  optionText.ID,
		QuestionID:    optionText.QuestionID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
	}

	_, er := s.Repository.CreateTextOptionHistory(c, oth)
	if er != nil {
		return er
	}

	e := s.Repository.DeleteTextOption(c, id)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) RestoreTextOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionText, err := s.Repository.GetTextOptionByID(c, id)
	if err != nil {
		return nil, err
	}

	oth := &TextOptionHistory{
		ID:            uuid.New(),
		OptionTextID:  optionText.ID,
		QuestionID:    optionText.QuestionID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
	}

	_, er := s.Repository.CreateTextOptionHistory(c, oth)
	if er != nil {
		return nil, er
	}

	optionText, e := s.Repository.RestoreTextOption(c, id)
	if e != nil {
		return nil, e
	}

	return &TextOptionResponse{
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