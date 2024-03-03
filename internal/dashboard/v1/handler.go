package v1

import (
	"net/http"
	"strings"

	l "github.com/Live-Quiz-Project/Backend/internal/live/v1"
	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"
	"github.com/Live-Quiz-Project/Backend/internal/util"

	// "github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
	quizService q.Service
	liveService l.Service
	userService u.Service
}

func NewHandler(s Service, qServ q.Service, lServ l.Service, uServ u.Service) *Handler {
	return &Handler{
		Service:     s,
		quizService: qServ,
		liveService: lServ,
		userService: uServ,
	}
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

	lqs, err := h.liveService.GetLiveQuizSessionBySessionID(c.Request.Context(), id)
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
			ID:        qph.ID,
			Type:      "POOL",
			PoolOrder: qph.PoolOrder,
			Order:     qph.Order,
			Content:   qph.Content,
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

			answerResponse, err := h.Service.GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(c.Request.Context(), lqs.ID, qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var oc []any
			for _, ocr := range ocRes {

				var answerParticipants []ParticipantResponse
				answerParticipants = nil

				for _, answerData := range answerResponse {
					if ocr.Content == answerData.Answer {
						participant, err := h.Service.GetParticipantByID(c.Request.Context(), answerData.ParticipantID)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}

						answerParticipants = append(answerParticipants, ParticipantResponse{
							ID:     participant.ID,
							UserID: participant.UserID,
							//LiveQuizSessionID: participant.LiveQuizSessionID,
							//Status:            participant.Status,
							Name: participant.Name,
							//Marks:             participant.Marks,
						})
					}
				}

				oc = append(oc, QuestionViewOptionChoice{
					ID:           ocr.ID,
					Order:        ocr.Order,
					Content:      ocr.Content,
					Mark:         ocr.Mark,
					Participants: answerParticipants,
				})
			}

			res.Questions = append(res.Questions, QuestionViewQuestionResponse{
				ID:             qr.ID,
				Type:           qr.Type,
				PoolOrder:      qr.PoolOrder,
				Order:          qr.Order,
				Content:        qr.Content,
				Note:           qr.Note,
				Media:          qr.Media,
				UseTemplate:    qr.UseTemplate,
				TimeLimit:      qr.TimeLimit,
				HaveTimeFactor: qr.HaveTimeFactor,
				TimeFactor:     qr.TimeFactor,
				FontSize:       qr.FontSize,
				SelectMin:      qr.SelectMin,
				SelectMax:      qr.SelectMax,
				Options:        oc,
			})
		}
		if qr.Type == util.FillBlank || qr.Type == util.Paragraph {
			otRes, err := h.quizService.GetTextOptionHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			answerResponse, err := h.Service.GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(c.Request.Context(), lqs.ID, qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var ot []any
			for _, otr := range otRes {

				var answerParticipants []ParticipantResponse
				answerParticipants = nil

				for _, answerData := range answerResponse {
					if otr.Content == answerData.Answer {
						participant, err := h.Service.GetParticipantByID(c.Request.Context(), answerData.ParticipantID)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}

						answerParticipants = append(answerParticipants, ParticipantResponse{
							ID:     participant.ID,
							UserID: participant.UserID,
							//LiveQuizSessionID: participant.LiveQuizSessionID,
							//Status:            participant.Status,
							Name: participant.Name,
							//Marks:             participant.Marks,
						})
					}
				}

				ot = append(ot, QuestionViewOptionText{
					ID:            otr.ID,
					Order:         otr.Order,
					Content:       otr.Content,
					Mark:          otr.Mark,
					CaseSensitive: otr.CaseSensitive,
					Participants:  answerParticipants,
				})
			}

			res.Questions = append(res.Questions, QuestionViewQuestionResponse{
				ID:             qr.ID,
				Type:           qr.Type,
				PoolOrder:      qr.PoolOrder,
				Order:          qr.Order,
				Content:        qr.Content,
				Note:           qr.Note,
				Media:          qr.Media,
				UseTemplate:    qr.UseTemplate,
				TimeLimit:      qr.TimeLimit,
				HaveTimeFactor: qr.HaveTimeFactor,
				TimeFactor:     qr.TimeFactor,
				FontSize:       qr.FontSize,
				SelectMin:      qr.SelectMin,
				SelectMax:      qr.SelectMax,
				Options:        ot,
			})
		}

		if qr.Type == util.Matching {
			amRes, err := h.quizService.GetMatchingAnswerHistoriesByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			answerResponse, err := h.Service.GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(c.Request.Context(), lqs.ID, qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var optionContent string
			var promptContent string

			var om []any
			for _, omr := range amRes {

				var answerParticipants []ParticipantResponse
				answerParticipants = nil

				for _, answerData := range answerResponse {
					splitAnswer := strings.Split(answerData.Answer, util.ANSWER_SPLIT)

					option, err := h.quizService.GetMatchingOptionHistoryByID(c.Request.Context(), omr.OptionID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					prompt, err := h.quizService.GetMatchingOptionHistoryByID(c.Request.Context(), omr.PromptID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}

					optionContent = option.Content
					promptContent = prompt.Content

					for _, pair := range splitAnswer {
						ans := strings.Split(pair, ":")

						if ans[0] == prompt.Content && ans[1] == option.Content {
							participant, err := h.Service.GetParticipantByID(c.Request.Context(), answerData.ParticipantID)
							if err != nil {
								c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
								return
							}

							answerParticipants = append(answerParticipants, ParticipantResponse{
								ID:     participant.ID,
								UserID: participant.UserID,
								// LiveQuizSessionID: participant.LiveQuizSessionID,
								// Status:            participant.Status,
								Name: participant.Name,
								// Marks:             participant.Marks,
							})
						}
					}
				}

				om = append(om, QuestionViewMatching{
					ID:            omr.ID,
					OptionID:      omr.OptionID,
					OptionContent: optionContent,
					PromptID:      omr.PromptID,
					PromptContent: promptContent,
					Mark:          omr.Mark,
					Participants:  answerParticipants,
				})
			}

			res.Questions = append(res.Questions, QuestionViewQuestionResponse{
				ID:             qr.ID,
				Type:           qr.Type,
				PoolOrder:      qr.PoolOrder,
				Order:          qr.Order,
				Content:        qr.Content,
				Note:           qr.Note,
				Media:          qr.Media,
				UseTemplate:    qr.UseTemplate,
				TimeLimit:      qr.TimeLimit,
				HaveTimeFactor: qr.HaveTimeFactor,
				TimeFactor:     qr.TimeFactor,
				FontSize:       qr.FontSize,
				SelectMin:      qr.SelectMin,
				SelectMax:      qr.SelectMax,
				Options:        om,
			})
		}
	}

	// Bubble Sort the Questions and Question Pool by Order
	n := len(res.Questions)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if res.Questions[j].Order > res.Questions[j+1].Order {
				res.Questions[j], res.Questions[j+1] = res.Questions[j+1], res.Questions[j]
			}
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetDashboardAnswerViewByID(c *gin.Context) {
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

	lqs, err := h.liveService.GetLiveQuizSessionBySessionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	quizH, err := h.quizService.GetQuizHistoryByID(c.Request.Context(), lqs.QuizID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := AnswerViewQuizResponse{
		ID:          quizH.ID,
		CreatorID:   quizH.CreatorID,
		Title:       quizH.Title,
		Description: quizH.Description,
		CoverImage:  quizH.CoverImage,
		CreatedAt:   quizH.CreatedAt,
	}

	participants, err := h.Service.GetOrderParticipantsByLiveQuizSessionID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, p := range participants {
		answers, err := h.Service.GetAnswerResponsesByLiveQuizSessionIDAndParticipantID(c.Request.Context(), lqs.ID, p.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var totalMarks int
		var totalTimeUsed int
		var correctAns int
		var incorrectAns int
		var unanswered int
		var questions []AnswerViewQuestionResponse

		for _, a := range answers {
			ansList := strings.Split(a.Answer, util.ANSWER_SPLIT)
			answerString := strings.Join(ansList, ", ")
			questionMark := 0

			q, err := h.quizService.GetQuestionHistoryByID(c.Request.Context(), a.QuestionID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if a.Type == util.Choice || a.Type == util.TrueFalse {
				for _, ans := range ansList {
					optionInfo, err := h.quizService.GetChoiceOptionHistoryByQuestionIDAndContent(c.Request.Context(), a.QuestionID, ans)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					questionMark += optionInfo.Mark
				}

				totalMarks += questionMark
				totalTimeUsed += a.UseTime

				questions = append(questions, AnswerViewQuestionResponse{
					ID:      a.ID,
					Type:    q.Type,
					Order:   q.Order,
					Content: q.Content,
					Answer:  answerString,
					Mark:    questionMark,
					UseTime: a.UseTime,
				})
			}
			if a.Type == util.FillBlank || a.Type == util.Paragraph {
				for _, ans := range ansList {
					optionInfo, err := h.quizService.GetTextOptionHistoryByQuestionIDAndContent(c.Request.Context(), a.QuestionID, ans)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					questionMark += optionInfo.Mark
				}

				totalMarks += questionMark
				totalTimeUsed += a.UseTime

				questions = append(questions, AnswerViewQuestionResponse{
					ID:      a.ID,
					Type:    q.Type,
					Order:   q.Order,
					Content: q.Content,
					Answer:  answerString,
					Mark:    questionMark,
					UseTime: a.UseTime,
				})
			}
			if a.Type == util.Matching {

				for _, ans := range ansList {
					pair := strings.Split(ans, ":")
					promptInfo, err := h.quizService.GetMatchingOptionHistoryByQuestionIDAndContent(c.Request.Context(), a.QuestionID, pair[0])
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}

					optionInfo, err := h.quizService.GetMatchingOptionHistoryByQuestionIDAndContent(c.Request.Context(), a.QuestionID, pair[1])
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}

					checkMatchingAnswer, err := h.quizService.GetMatchingAnswerHistoryByPromptIDAndOptionID(c.Request.Context(), promptInfo.ID, optionInfo.ID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}

					questionMark += checkMatchingAnswer.Mark
				}

				totalMarks += questionMark
				totalTimeUsed += a.UseTime

				questions = append(questions, AnswerViewQuestionResponse{
					ID:      a.ID,
					Type:    q.Type,
					Order:   q.Order,
					Content: q.Content,
					Answer:  answerString,
					Mark:    questionMark,
					UseTime: a.UseTime,
				})
			}
			
			// Check Correct Answer
			if a.Answer == "" {
				unanswered += 1
			} else if questionMark != 0 && a.Answer != "" {
				correctAns += 1
			} else if questionMark == 0 && a.Answer != "" {
				incorrectAns += 1
			}

		}

		n := len(questions)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if questions[j].Order > questions[j+1].Order {
					questions[j], questions[j+1] = questions[j+1], questions[j]
				}
			}
		}

		res.Participants = append(res.Participants, AnswerViewParticipantResponse{
			ID:             p.ID,
			UserID:         p.UserID,
			Name:           p.Name,
			Marks:          p.Marks,
			Corrects:       correctAns,
			Incorrects:     incorrectAns,
			Unanswered:     unanswered,
			TotalQuestions: len(questions),
			TotalMarks:     totalMarks,
			TotalTimeUsed:  totalTimeUsed,
			Questions:      questions,
		},
		)
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetDashboardHistoryByUserID(c *gin.Context) {
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

	sessions, err := h.liveService.GetLiveQuizSessionsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var sessionHistory []SessionHistory

	for _, eachSession := range sessions {
		userInfo, err := h.userService.GetUserByID(c.Request.Context(), eachSession.HostID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalParticipant, err := h.Service.CountTotalParticipants(c.Request.Context(), eachSession.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		quizInfo, err := h.quizService.GetQuizHistoryByID(c.Request.Context(), eachSession.QuizID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sessionHistory = append(sessionHistory, SessionHistory{
			ID:                eachSession.ID,
			CreatorName:       userInfo.Name,
			Title:             quizInfo.Title,
			Description:       quizInfo.Description,
			CoverImage:        quizInfo.CoverImage,
			Visibility:        quizInfo.Visibility,
			TotalParticipants: totalParticipant,
			CreatedAt:         eachSession.CreatedAt,
			UpdatedAt:         quizInfo.UpdatedAt,
			DeletedAt:         quizInfo.DeletedAt,
		})
	}

	c.JSON(200, sessionHistory)
}
