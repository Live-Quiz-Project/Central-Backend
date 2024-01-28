package v1

import (
	"net/http"

	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
	quizService q.Service
	liveService l.Service
}

func NewHandler(s Service, qServ q.Service, lServ l.Service) *Handler {
	return &Handler{
		Service:     s,
		quizService: qServ,
		liveService: lServ,
	}
}

// func (h *Handler) CreateAnswerResponse(c *gin.Context) {

// 	var req LiveAnswerRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	uid, ok := c.Get("uid")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	userID, err := uuid.Parse(uid.(string))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
// 		return
// 	}

// 	res, er := h.Service.CreateAnswerResponse(c.Request.Context(), &req)
// 	if er != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
// 		return
// 	}

// 	for _, answer := range {

// 	}

// 	for _, q := range req.Questions {
// 		var qRes *CreateQuestionResponse
// 		var qpRes *CreateQuestionPoolResponse

// 		if q.Type == util.Pool {
// 			qpRes, err = h.Service.CreateQuestionPool(c.Request.Context(), &q, res.ID, res.QuizHistoryID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			qpResID = &qpRes.ID
// 			qphResID = &qpRes.QuestionPoolHistoryID

// 			res.Questions = append(res.Questions, QuestionResponse{
// 				Question: Question{
// 					ID:             qpRes.ID,
// 					QuizID:         qpRes.QuizID,
// 					Type:           "POOL",
// 					Order:          qpRes.Order,
// 					Content:        qpRes.Content,
// 					Note:           qpRes.Note,
// 					Media:          qpRes.Media,
// 					TimeLimit:      qpRes.TimeLimit,
// 					HaveTimeFactor: qpRes.HaveTimeFactor,
// 					TimeFactor:     qpRes.TimeFactor,
// 					FontSize:       qpRes.FontSize,
// 					CreatedAt:      qpRes.CreatedAt,
// 					UpdatedAt:      qpRes.UpdatedAt,
// 					DeletedAt:      qpRes.DeletedAt,
// 				},
// 			})

// 			continue
// 		} else {
// 			if q.IsInPool == true {
// 				qRes, err = h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, qpResID, qphResID, userID)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 					return
// 				}
// 			} else {
// 				qRes, err = h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, nil, nil, userID)
// 				if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 					return
// 				}
// 			}

// 			for _, qt := range qRes.Options {
// 				if qst, ok := qt.(map[string]any); ok {
// 					if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
// 						_, err := h.Service.CreateChoiceOption(c.Request.Context(), &ChoiceOptionRequest{
// 							ChoiceOption: ChoiceOption{
// 								Order:   int(qst["order"].(float64)),
// 								Content: qst["content"].(string),
// 								Mark:    int(qst["mark"].(float64)),
// 								Color:   qst["color"].(string),
// 								Correct: qst["correct"].(bool),
// 							},
// 						}, qRes.ID, qRes.QuestionHistoryID, userID)

// 						if err != nil {
// 							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 							return
// 						}

// 					} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
// 						_, err := h.Service.CreateTextOption(c.Request.Context(), &TextOptionRequest{
// 							TextOption: TextOption{
// 								Order:         int(qst["order"].(float64)),
// 								Content:       qst["content"].(string),
// 								Mark:          int(qst["mark"].(float64)),
// 								CaseSensitive: qst["case_sensitive"].(bool),
// 							},
// 						}, qRes.ID, qRes.QuestionHistoryID, userID)

// 						if err != nil {
// 							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 							return
// 						}

// 					} else if qRes.Type == util.Matching {
// 						if qst["type"].(string) != "MATCHING_ANSWER" {
// 							_, err := h.Service.CreateMatchingOption(c.Request.Context(), &MatchingOptionRequest{
// 								MatchingOption: MatchingOption{
// 									Order:     int(qst["order"].(float64)),
// 									Content:   qst["content"].(string),
// 									Type:      qst["type"].(string),
// 									Eliminate: qst["eliminate"].(bool),
// 								},
// 							}, qRes.ID, qRes.QuestionHistoryID, userID)

// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}

// 						} else {

// 							prompt, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["prompt"].(float64)))
// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}

// 							option, err := h.Service.GetMatchingOptionByQuestionIDAndOrder(c.Request.Context(), qRes.ID, int(qst["option"].(float64)))
// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}

// 							_, err = h.Service.CreateMatchingAnswer(c.Request.Context(), &MatchingAnswerRequest{
// 								MatchingAnswer: MatchingAnswer{
// 									PromptID: prompt.ID,
// 									OptionID: option.ID,
// 									Mark:     int(qst["mark"].(float64)),
// 								},
// 							}, qRes.ID, qRes.QuestionHistoryID, userID)

// 							if err != nil {
// 								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 								return
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}

// 		res.Questions = append(res.Questions, QuestionResponse{
// 			Question: Question{
// 				ID:             qRes.ID,
// 				QuizID:         qRes.QuizID,
// 				QuestionPoolID: qRes.QuestionPoolID,
// 				Type:           qRes.Type,
// 				Order:          qRes.Order,
// 				Content:        qRes.Content,
// 				Note:           qRes.Note,
// 				Media:          qRes.Media,
// 				UseTemplate:    qRes.UseTemplate,
// 				TimeLimit:      qRes.TimeLimit,
// 				HaveTimeFactor: qRes.HaveTimeFactor,
// 				TimeFactor:     qRes.TimeFactor,
// 				FontSize:       qRes.FontSize,
// 				LayoutIdx:      qRes.LayoutIdx,
// 				SelectUpTo:     qRes.SelectUpTo,
// 				CreatedAt:      qRes.CreatedAt,
// 				UpdatedAt:      qRes.UpdatedAt,
// 				DeletedAt:      qRes.DeletedAt,
// 			},
// 			Options: qRes.Options,
// 		})
// 	}

// 	c.JSON(http.StatusCreated, &QuizResponse{
// 		Quiz: Quiz{
// 			ID:             res.ID,
// 			CreatorID:      res.CreatorID,
// 			Title:          res.Title,
// 			Description:    res.Description,
// 			CoverImage:     res.CoverImage,
// 			Visibility:     res.Visibility,
// 			TimeLimit:      res.TimeLimit,
// 			HaveTimeFactor: res.HaveTimeFactor,
// 			TimeFactor:     res.TimeFactor,
// 			FontSize:       res.FontSize,
// 			Mark:           res.Mark,
// 			SelectUpTo:     res.SelectUpTo,
// 			CaseSensitive:  res.CaseSensitive,
// 			CreatedAt:      res.CreatedAt,
// 			UpdatedAt:      res.UpdatedAt,
// 			DeletedAt:      res.DeletedAt,
// 		},
// 		Questions: res.Questions,
// 	})
// }

func (h *Handler) GetAnswerResponseByLiveQuizSessionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.Service.GetAnswerResponseByLiveQuizSessionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) GetAnswerResponseByQuestionID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.Service.GetAnswerResponseByQuestionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *Handler) GetAnswerResponseByParticipantID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.Service.GetAnswerResponseByParticipantID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetDashboardQuestionViewByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id")) // id = live_quiz_session_id
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lqs, err := h.liveService.GetLiveQuizSessionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quizH, err := h.quizService.GetQuizHistoryByID(c.Request.Context(), lqs.QuizID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := QuestionViewQuizResponse{
		ID:          quizH.ID,
		CreatorID:   quizH.CreatorID,
		Title:       quizH.Title,
		Description: quizH.Description,
		CoverImage:  quizH.CoverImage,
		CreatedAt:   quizH.CreatedAt,
	}

	questionPoolH, err := h.quizService.GetQuestionPoolHistoriesByQuizID(c.Request.Context(), lqs.QuizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, qph := range questionPoolH {
		res.Questions = append(res.Questions, QuestionViewQuestionResponse{
			ID:      qph.ID,
			Order:   qph.Order,
			Content: qph.Content,
			Type:    "POOL",
		})
	}

	questionH, err := h.quizService.GetQuestionHistoriesByQuizID(c.Request.Context(), lqs.QuizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, qr := range questionH {
		if qr.Type == util.Choice || qr.Type == util.TrueFalse {
			ocRes, err := h.quizService.GetChoiceOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			answerResponse, err := h.Service.GetAnswerResponsesByLiveQuizSessionIDAndQuestionID(c.Request.Context(), lqs.ID, qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var allParticipants []ParticipantResponse

			for _, answerData := range answerResponse {
				participant, err := h.Service.GetParticipantByID(c.Request.Context(), answerData.ParticipantID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				allParticipants = append(allParticipants, ParticipantResponse{
					Participant: Participant{
						ID:                participant.ID,
						UserID:            participant.UserID,
						LiveQuizSessionID: participant.LiveQuizSessionID,
						Status:            participant.Status,
						Name:              participant.Name,
						Marks:             participant.Marks,
					},
				})
			}

			var oc []any
			for _, ocr := range ocRes {
				oc = append(oc, QuestionViewOptionChoice{
					ID:      ocr.ID,
					Order:   ocr.Order,
					Content: ocr.Content,
					Mark:    ocr.Mark,
				})
			}

			res.Questions = append(res.Questions, QuestionViewQuestionResponse{
				ID:             qr.ID,
				Type:           qr.Type,
				Order:          qr.Order,
				Content:        qr.Content,
				Note:           qr.Note,
				Media:          qr.Media,
				UseTemplate:    qr.UseTemplate,
				TimeLimit:      qr.TimeLimit,
				HaveTimeFactor: qr.HaveTimeFactor,
				TimeFactor:     qr.TimeFactor,
				FontSize:       qr.FontSize,
				SelectUpTo:     qr.SelectUpTo,
				Options:        oc,
			})

			// } else if qr.Type == util.ShortText || qr.Type == util.Paragraph {
			// 	otRes, err := h.Service.GetTextOptionsByQuestionID(c.Request.Context(), qr.ID)
			// 	if err != nil {
			// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 		return
			// 	}

			// 	var ot []any
			// 	for _, otr := range otRes {
			// 		ot = append(ot, TextOptionResponse{
			// 			TextOption: TextOption{
			// 				ID:            otr.ID,
			// 				QuestionID:    otr.QuestionID,
			// 				Order:         otr.Order,
			// 				Content:       otr.Content,
			// 				Mark:          otr.Mark,
			// 				CaseSensitive: otr.CaseSensitive,
			// 				CreatedAt:     otr.CreatedAt,
			// 				UpdatedAt:     otr.UpdatedAt,
			// 				DeletedAt:     otr.DeletedAt,
			// 			},
			// 		})
			// 	}

			// 	res.Questions = append(res.Questions, QuestionResponse{
			// 		Question: Question{
			// 			ID:             qr.ID,
			// 			QuizID:         qr.QuizID,
			// 			QuestionPoolID: qr.QuestionPoolID,
			// 			Type:           qr.Type,
			// 			Order:          qr.Order,
			// 			Content:        qr.Content,
			// 			Note:           qr.Note,
			// 			Media:          qr.Media,
			// 			UseTemplate:    qr.UseTemplate,
			// 			TimeLimit:      qr.TimeLimit,
			// 			HaveTimeFactor: qr.HaveTimeFactor,
			// 			TimeFactor:     qr.TimeFactor,
			// 			FontSize:       qr.FontSize,
			// 			LayoutIdx:      qr.LayoutIdx,
			// 			SelectUpTo:     qr.SelectUpTo,
			// 			CreatedAt:      qr.CreatedAt,
			// 			UpdatedAt:      qr.UpdatedAt,
			// 		},
			// 		Options: ot,
			// 	})
			// } else if qr.Type == util.Matching {
			// 	omRes, err := h.Service.GetMatchingOptionsByQuestionID(c.Request.Context(), qr.ID)
			// 	if err != nil {
			// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 		return
			// 	}

			// 	amRes, err := h.Service.GetMatchingAnswersByQuestionID(c.Request.Context(), qr.ID)
			// 	if err != nil {
			// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// 		return
			// 	}

			// 	var o []any
			// 	for _, omr := range omRes {
			// 		o = append(o, MatchingOptionAndAnswerResponse{
			// 			ID:         omr.ID,
			// 			QuestionID: omr.QuestionID,
			// 			Type:       omr.Type,
			// 			Order:      omr.Order,
			// 			Content:    omr.Content,
			// 			Eliminate:  omr.Eliminate,
			// 			PromptID:   uuid.Nil,
			// 			OptionID:   uuid.Nil,
			// 			Mark:       0,
			// 			CreatedAt:  omr.CreatedAt,
			// 			UpdatedAt:  omr.UpdatedAt,
			// 			DeletedAt:  omr.DeletedAt,
			// 		})
			// 	}

			// 	for _, amr := range amRes {
			// 		o = append(o, MatchingOptionAndAnswerResponse{
			// 			ID:         amr.ID,
			// 			QuestionID: amr.QuestionID,
			// 			Type:       "MATCHING_ANSWER",
			// 			Order:      0,
			// 			Content:    "",
			// 			Eliminate:  false,
			// 			PromptID:   amr.PromptID,
			// 			OptionID:   amr.OptionID,
			// 			Mark:       amr.Mark,
			// 			CreatedAt:  amr.CreatedAt,
			// 			UpdatedAt:  amr.UpdatedAt,
			// 			DeletedAt:  amr.DeletedAt,
			// 		})
			// 	}

			// 	res.Questions = append(res.Questions, QuestionResponse{
			// 		Question: Question{
			// 			ID:             qr.ID,
			// 			QuizID:         qr.QuizID,
			// 			QuestionPoolID: qr.QuestionPoolID,
			// 			Type:           qr.Type,
			// 			Order:          qr.Order,
			// 			Content:        qr.Content,
			// 			Note:           qr.Note,
			// 			Media:          qr.Media,
			// 			UseTemplate:    qr.UseTemplate,
			// 			TimeLimit:      qr.TimeLimit,
			// 			HaveTimeFactor: qr.HaveTimeFactor,
			// 			TimeFactor:     qr.TimeFactor,
			// 			FontSize:       qr.FontSize,
			// 			LayoutIdx:      qr.LayoutIdx,
			// 			SelectUpTo:     qr.SelectUpTo,
			// 			CreatedAt:      qr.CreatedAt,
			// 			UpdatedAt:      qr.UpdatedAt,
			// 		},
			// 		Options: o,
			// 	})
		}
	}

}
