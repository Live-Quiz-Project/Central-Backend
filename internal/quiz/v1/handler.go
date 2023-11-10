package v1

import (
	"net/http"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// ---------- Quiz related handlers ---------- //
func (h *Handler) CreateQuiz(c *gin.Context) {
	var req CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	res, er := h.Service.CreateQuiz(c.Request.Context(), &req, userID)
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
		return
	}

	for _, q := range req.Questions {
		qRes, err := h.Service.CreateQuestion(c.Request.Context(), &q, res.ID, res.QuizHistoryID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, qt := range qRes.Options {
			if qst, ok := qt.(map[string]any); ok {
				if qRes.Type == util.Choice || qRes.Type == util.TrueFalse {
					_, err := h.Service.CreateChoiceOption(c.Request.Context(), &CreateChoiceOptionRequest{
						Order:   int(qst["order"].(float64)),
						Content: qst["content"].(string),
						Mark:    int(qst["mark"].(float64)),
						Color:   qst["color"].(string),
						Correct: qst["correct"].(bool),
					}, qRes.ID, qRes.QuestionHistoryID, userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				} else if qRes.Type == util.ShortText || qRes.Type == util.Paragraph {
					_, err := h.Service.CreateTextOption(c.Request.Context(), &CreateTextOptionRequest{
						Order:         int(qst["order"].(float64)),
						Content:       qst["content"].(string),
						Mark:          int(qst["mark"].(float64)),
						CaseSensitive: qst["case_sensitive"].(bool),
					}, qRes.ID, qRes.QuestionHistoryID, userID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
			}
		}

		res.Questions = append(res.Questions, QuestionResponse{
			ID:             qRes.ID,
			QuizID:         qRes.QuizID,
			ParentID:       qRes.ParentID,
			Type:           qRes.Type,
			Order:          qRes.Order,
			Content:        qRes.Content,
			Note:           qRes.Note,
			Media:          qRes.Media,
			TimeLimit:      qRes.TimeLimit,
			HaveTimeFactor: qRes.HaveTimeFactor,
			TimeFactor:     qRes.TimeFactor,
			FontSize:       qRes.FontSize,
			LayoutIdx:      qRes.LayoutIdx,
			SelectedUpTo:   qRes.SelectedUpTo,
			SubQuestions:   qRes.SubQuestions,
			Options:        qRes.Options,
		})
	}

	c.JSON(http.StatusCreated, &QuizResponse{
		ID:          res.ID,
		CreatorID:   res.CreatorID,
		Title:       res.Title,
		Description: res.Description,
		CoverImage:  res.CoverImage,
		Questions:   res.Questions,
	})
}

func (h *Handler) GetQuizzes(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	res, err := h.Service.GetQuizzes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r := make([]QuizResponse, 0)
	for _, q := range res {
		qRes, err := h.Service.GetQuestionsByQuizID(c.Request.Context(), q.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		for _, qr := range qRes {
			if qr.Type == util.Choice || qr.Type == util.TrueFalse {
				ocRes, err := h.Service.GetChoiceOptionsByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				var oc []any
				for _, ocr := range ocRes {
					oc = append(oc, ChoiceOptioneResponse{
						ID:      ocr.ID,
						Order:   ocr.Order,
						Mark:    ocr.Mark,
						Color:   ocr.Color,
						Correct: ocr.Correct,
						Content: ocr.Content,
					})
				}

				q.Questions = append(q.Questions, QuestionResponse{
					ID:             qr.ID,
					QuizID:         qr.QuizID,
					ParentID:       qr.ParentID,
					Type:           qr.Type,
					Order:          qr.Order,
					Content:        qr.Content,
					Note:           qr.Note,
					Media:          qr.Media,
					TimeLimit:      qr.TimeLimit,
					HaveTimeFactor: qr.HaveTimeFactor,
					TimeFactor:     qr.TimeFactor,
					FontSize:       qr.FontSize,
					LayoutIdx:      qr.LayoutIdx,
					SelectedUpTo:   qr.SelectedUpTo,
					SubQuestions:   qr.SubQuestions,
					Options:        oc,
				})
			} else if qr.Type == util.ShortText || qr.Type == util.Paragraph {
				otRes, err := h.Service.GetTextOptionsByQuestionID(c.Request.Context(), qr.ID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				var ot []any
				for _, otr := range otRes {
					ot = append(ot, TextOptionResponse{
						ID:            otr.ID,
						Order:         otr.Order,
						Mark:          otr.Mark,
						CaseSensitive: otr.CaseSensitive,
						Content:       otr.Content,
					})
				}

				q.Questions = append(q.Questions, QuestionResponse{
					ID:             qr.ID,
					QuizID:         qr.QuizID,
					ParentID:       qr.ParentID,
					Type:           qr.Type,
					Order:          qr.Order,
					Content:        qr.Content,
					Note:           qr.Note,
					Media:          qr.Media,
					TimeLimit:      qr.TimeLimit,
					HaveTimeFactor: qr.HaveTimeFactor,
					TimeFactor:     qr.TimeFactor,
					FontSize:       qr.FontSize,
					LayoutIdx:      qr.LayoutIdx,
					SelectedUpTo:   qr.SelectedUpTo,
					SubQuestions:   qr.SubQuestions,
					Options:        ot,
				})
			}
		}

		r = append(r, q)
	}

	c.JSON(http.StatusOK, r)
}

func (h *Handler) GetQuizByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	res, err := h.Service.GetQuizByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	qRes, err := h.Service.GetQuestionsByQuizID(c.Request.Context(), res.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, qr := range qRes {
		if qr.Type == util.Choice || qr.Type == util.TrueFalse {
			ocRes, err := h.Service.GetChoiceOptionsByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var oc []any
			for _, ocr := range ocRes {
				oc = append(oc, ChoiceOptioneResponse{
					ID:      ocr.ID,
					Order:   ocr.Order,
					Mark:    ocr.Mark,
					Color:   ocr.Color,
					Correct: ocr.Correct,
					Content: ocr.Content,
				})
			}

			res.Questions = append(res.Questions, QuestionResponse{
				ID:             qr.ID,
				QuizID:         qr.QuizID,
				ParentID:       qr.ParentID,
				Type:           qr.Type,
				Order:          qr.Order,
				Content:        qr.Content,
				Note:           qr.Note,
				Media:          qr.Media,
				TimeLimit:      qr.TimeLimit,
				HaveTimeFactor: qr.HaveTimeFactor,
				TimeFactor:     qr.TimeFactor,
				FontSize:       qr.FontSize,
				LayoutIdx:      qr.LayoutIdx,
				SelectedUpTo:   qr.SelectedUpTo,
				SubQuestions:   qr.SubQuestions,
				Options:        oc,
			})
		} else if qr.Type == util.ShortText || qr.Type == util.Paragraph {
			otRes, err := h.Service.GetTextOptionsByQuestionID(c.Request.Context(), qr.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			var ot []any
			for _, otr := range otRes {
				ot = append(ot, TextOptionResponse{
					ID:            otr.ID,
					Order:         otr.Order,
					Mark:          otr.Mark,
					CaseSensitive: otr.CaseSensitive,
					Content:       otr.Content,
				})
			}

			res.Questions = append(res.Questions, QuestionResponse{
				ID:             qr.ID,
				QuizID:         qr.QuizID,
				ParentID:       qr.ParentID,
				Type:           qr.Type,
				Order:          qr.Order,
				Content:        qr.Content,
				Note:           qr.Note,
				Media:          qr.Media,
				TimeLimit:      qr.TimeLimit,
				HaveTimeFactor: qr.HaveTimeFactor,
				TimeFactor:     qr.TimeFactor,
				FontSize:       qr.FontSize,
				LayoutIdx:      qr.LayoutIdx,
				SelectedUpTo:   qr.SelectedUpTo,
				SubQuestions:   qr.SubQuestions,
				Options:        ot,
			})
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateQuiz(c *gin.Context) {
	var req UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	res, err := h.Service.UpdateQuiz(c.Request.Context(), &req, quizID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteQuiz(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	er := h.Service.DeleteQuiz(c.Request.Context(), quizID, userID)
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": er.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted",
	})
}

func (h *Handler) RestoreQuiz(c *gin.Context) {}
