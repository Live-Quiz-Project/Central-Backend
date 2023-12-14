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
		ID:          uuid.New(),
		CreatorID:   uid,
		Title:       req.Title,
		Description: req.Description,
		CoverImage:  "default.png",
	}

	qh := &QuizHistory{
		ID:          uuid.New(),
		QuizID:      q.ID,
		CreatorID:   q.CreatorID,
		Title:       q.Title,
		Description: q.Description,
		CoverImage:  q.CoverImage,
		UpdatedBy:   q.CreatorID,
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
			ID:          quiz.ID,
			CreatorID:   quiz.CreatorID,
			Title:       quiz.Title,
			Description: quiz.Description,
			CoverImage:  quiz.CoverImage,
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
			ID:          q.ID,
			CreatorID:   q.CreatorID,
			Title:       q.Title,
			Description: q.Description,
			CoverImage:  q.CoverImage,
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
		ID:          quiz.ID,
		CreatorID:   quiz.CreatorID,
		Title:       quiz.Title,
		Description: quiz.Description,
		CoverImage:  quiz.CoverImage,
	}, nil
}

func (s *service) UpdateQuiz(ctx context.Context, req *UpdateQuizRequest, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	quiz, err := s.Repository.GetQuizByID(c, id)
	if err != nil {
		return &QuizResponse{}, err
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

	qh := &QuizHistory{
		ID:          uuid.New(),
		QuizID:      quiz.ID,
		CreatorID:   quiz.CreatorID,
		Title:       quiz.Title,
		Description: quiz.Description,
		CoverImage:  quiz.CoverImage,
		UpdatedBy:   uid,
	}

	quiz, er := s.Repository.UpdateQuiz(c, quiz)
	if er != nil {
		return &QuizResponse{}, er
	}

	_, e := s.Repository.CreateQuizHistory(c, qh)
	if e != nil {
		return &QuizResponse{}, e
	}

	return &QuizResponse{
		ID:          quiz.ID,
		CreatorID:   quiz.CreatorID,
		Title:       quiz.Title,
		Description: quiz.Description,
		CoverImage:  quiz.CoverImage,
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
		ID:          uuid.New(),
		QuizID:      quiz.ID,
		CreatorID:   quiz.CreatorID,
		Title:       quiz.Title,
		Description: quiz.Description,
		CoverImage:  quiz.CoverImage,
		UpdatedBy:   uid,
		Deleted:     true,
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
		ID:          uuid.New(),
		QuizID:      quiz.ID,
		CreatorID:   quiz.CreatorID,
		Title:       quiz.Title,
		Description: quiz.Description,
		CoverImage:  quiz.CoverImage,
		UpdatedBy:   uid,
		Deleted:     false,
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
		ID:          quiz.ID,
		CreatorID:   quiz.CreatorID,
		Title:       quiz.Title,
		Description: quiz.Description,
		CoverImage:  quiz.CoverImage,
	}, nil
}

// ---------- Question related service methods ---------- //
func (s *service) CreateQuestion(ctx context.Context, req *CreateQuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID, uid uuid.UUID) (*CreateQuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	q := &Question{
		ID:             uuid.New(),
		QuizID:         quizID,
		ParentID:       nil,
		Type:           req.Type,
		Order:          req.Order,
		Content:        req.Content,
		Note:           req.Note,
		Media:          req.Media,
		TimeLimit:      req.TimeLimit,
		HaveTimeFactor: req.HaveTimeFactor,
		TimeFactor:     req.TimeFactor,
		FontSize:       req.FontSize,
		LayoutIdx:      req.LayoutIdx,
		SelectedUpTo:   req.SelectedUpTo,
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     q.ID,
		QuizHistoryID:  quizHistoryID,
		ParentID:       q.ParentID,
		Type:           q.Type,
		Order:          q.Order,
		Content:        q.Content,
		Note:           q.Note,
		Media:          q.Media,
		TimeLimit:      q.TimeLimit,
		HaveTimeFactor: q.HaveTimeFactor,
		TimeFactor:     q.TimeFactor,
		FontSize:       q.FontSize,
		LayoutIdx:      q.LayoutIdx,
		SelectedUpTo:   q.SelectedUpTo,
		UpdatedBy:      uid,
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
			ID:             question.ID,
			QuizID:         question.QuizID,
			ParentID:       question.ParentID,
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
			SelectedUpTo:   question.SelectedUpTo,
			Options:        options,
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
			ID:             q.ID,
			QuizID:         q.QuizID,
			ParentID:       q.ParentID,
			Type:           q.Type,
			Order:          q.Order,
			Content:        q.Content,
			Note:           q.Note,
			Media:          q.Media,
			TimeLimit:      q.TimeLimit,
			HaveTimeFactor: q.HaveTimeFactor,
			TimeFactor:     q.TimeFactor,
			FontSize:       q.FontSize,
			LayoutIdx:      q.LayoutIdx,
			SelectedUpTo:   q.SelectedUpTo,
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

func (s *service) UpdateQuestion(ctx context.Context, req *UpdateQuestionRequest, id uuid.UUID, uid uuid.UUID) (*QuestionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	question, err := s.Repository.GetQuestionByID(c, id)
	if err != nil {
		return &QuestionResponse{}, err
	}

	if req.Type != "" {
		question.Type = req.Type
	}
	if req.Order != 0 {
		question.Order = req.Order
	}
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
	if req.SelectedUpTo != 0 {
		question.SelectedUpTo = req.SelectedUpTo
	}

	qh := &QuestionHistory{
		ID:             uuid.New(),
		QuestionID:     question.ID,
		QuizHistoryID:  question.QuizID,
		ParentID:       question.ParentID,
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
		SelectedUpTo:   question.SelectedUpTo,
		UpdatedBy:      uid,
	}

	question, er := s.Repository.UpdateQuestion(c, question)
	if er != nil {
		return &QuestionResponse{}, er
	}

	_, e := s.Repository.CreateQuestionHistory(c, qh)
	if e != nil {
		return &QuestionResponse{}, e
	}

	return &QuestionResponse{
		ID:             question.ID,
		QuizID:         question.QuizID,
		ParentID:       question.ParentID,
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
		SelectedUpTo:   question.SelectedUpTo,
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
		QuizHistoryID:  question.QuizID,
		ParentID:       question.ParentID,
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
		SelectedUpTo:   question.SelectedUpTo,
		UpdatedBy:      uid,
		Deleted:        true,
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
		QuizHistoryID:  question.QuizID,
		ParentID:       question.ParentID,
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
		SelectedUpTo:   question.SelectedUpTo,
		UpdatedBy:      uid,
		Deleted:        false,
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
		ID:             question.ID,
		QuizID:         question.QuizID,
		ParentID:       question.ParentID,
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
		SelectedUpTo:   question.SelectedUpTo,
	}, nil
}

// ---------- Options related service methods ---------- //
// Choice related service methods
func (s *service) CreateChoiceOption(ctx context.Context, req *CreateChoiceOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error) {
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
		UpdatedBy:      uid,
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
		ID:         optionChoice.ID,
		QuestionID: optionChoice.QuestionID,
		Order:      optionChoice.Order,
		Content:    optionChoice.Content,
		Mark:       optionChoice.Mark,
		Color:      optionChoice.Color,
		Correct:    optionChoice.Correct,
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
			ID:         oc.ID,
			QuestionID: oc.QuestionID,
			Order:      oc.Order,
			Content:    oc.Content,
			Mark:       oc.Mark,
			Color:      oc.Color,
			Correct:    oc.Correct,
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
			ID:         oc.ID,
			QuestionID: oc.QuestionID,
			Order:      oc.Order,
			Content:    oc.Content,
			Mark:       oc.Mark,
			Color:      oc.Color,
			Correct:    oc.Correct,
		})
	}

	return res, nil
}

func (s *service) UpdateChoiceOption(ctx context.Context, req *UpdateChoiceOptionRequest, id uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionChoice, err := s.Repository.GetChoiceOptionByID(c, id)
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
		QuestionID:     optionChoice.QuestionID,
		Order:          optionChoice.Order,
		Content:        optionChoice.Content,
		Mark:           optionChoice.Mark,
		Color:          optionChoice.Color,
		Correct:        optionChoice.Correct,
		UpdatedBy:      uid,
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
		ID:         optionChoice.ID,
		QuestionID: optionChoice.QuestionID,
		Order:      optionChoice.Order,
		Content:    optionChoice.Content,
		Mark:       optionChoice.Mark,
		Color:      optionChoice.Color,
		Correct:    optionChoice.Correct,
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
		UpdatedBy:      uid,
		Deleted:        true,
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
		UpdatedBy:      uid,
		Deleted:        false,
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
		ID:         optionChoice.ID,
		QuestionID: optionChoice.QuestionID,
		Order:      optionChoice.Order,
		Content:    optionChoice.Content,
		Mark:       optionChoice.Mark,
		Color:      optionChoice.Color,
		Correct:    optionChoice.Correct,
	}, nil
}

// Text related service methods
func (s *service) CreateTextOption(ctx context.Context, req *CreateTextOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error) {
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
		UpdatedBy:     uid,
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
		ID:            optionText.ID,
		QuestionID:    optionText.QuestionID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
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
			ID:            ot.ID,
			QuestionID:    ot.QuestionID,
			Order:         ot.Order,
			Content:       ot.Content,
			Mark:          ot.Mark,
			CaseSensitive: ot.CaseSensitive,
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
			ID:            ot.ID,
			QuestionID:    ot.QuestionID,
			Order:         ot.Order,
			Content:       ot.Content,
			Mark:          ot.Mark,
			CaseSensitive: ot.CaseSensitive,
		})
	}

	return res, nil
}

func (s *service) UpdateTextOption(ctx context.Context, req *UpdateTextOptionRequest, id uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	optionText, err := s.Repository.GetTextOptionByID(c, id)
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
		QuestionID:    optionText.QuestionID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
		UpdatedBy:     uid,
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
		ID:            optionText.ID,
		QuestionID:    optionText.QuestionID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
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
		UpdatedBy:     uid,
		Deleted:       true,
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
		UpdatedBy:     uid,
		Deleted:       false,
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
		ID:            optionText.ID,
		QuestionID:    optionText.QuestionID,
		Order:         optionText.Order,
		Content:       optionText.Content,
		Mark:          optionText.Mark,
		CaseSensitive: optionText.CaseSensitive,
	}, nil
}
